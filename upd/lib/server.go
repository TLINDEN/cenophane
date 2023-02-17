/*
Copyright © 2023 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package lib

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tlinden/up/upd/cfg"
	bolt "go.etcd.io/bbolt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Result struct {
	success bool
	url     string
	error   string
}

func Log(format string, values ...any) {
	fmt.Fprintf(gin.DefaultWriter, "[GIN] "+format+"\n", values...)
}

func Ts() string {
	t := time.Now()
	return t.Format("2006-01-02-15-04-")
}

func NormalizeFilename(file string) string {
	r := regexp.MustCompile(`[^\w\d\-_\\.]`)

	return Ts() + r.ReplaceAllString(file, "")
}

// Binding from JSON, data coming from user, not tainted
type Meta struct {
	Expire string `form:"expire"`
}

type DbEntry struct {
	Id     string `json:"id"`
	Expire string `json:"expire"`
	File   string `json:"file"`
}

func Runserver(cfg *cfg.Config, args []string) error {
	dst := cfg.StorageDir
	router := gin.Default()
	router.SetTrustedProxies(nil)
	api := router.Group(cfg.ApiPrefix)

	db, err := bolt.Open(cfg.DbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	// FIXME: put these beast into their own funcs!!!!!!!
	{
		api.POST("/putfile", func(c *gin.Context) {
			// supports upload of multiple files with:
			//
			// curl -X POST localhost:8080/putfile \
			//   -F "upload[]=@/home/scip/2023-02-06_10-51.png" \
			//   -F "upload[]=@/home/scip/pgstat.png" \
			//   -H "Content-Type: multipart/form-data"
			//
			// If multiple files are  uploaded, they are zipped, otherwise
			// the  file is being stored  as is.
			//
			// Returns the  name of the uploaded file.
			//
			// FIXME: normalize or rename filename of single file to avoid dumb file names
			id := uuid.NewString()

			os.MkdirAll(filepath.Join(dst, id), os.ModePerm)

			var formdata Meta
			if err := c.ShouldBind(&formdata); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			form, _ := c.MultipartForm()
			files := form.File["upload[]"]
			uploaded := []string{}

			for _, file := range files {
				filename := NormalizeFilename(filepath.Base(file.Filename))
				path := filepath.Join(dst, id, filename)
				uploaded = append(uploaded, filename)
				Log("Received: %s => %s/%s", file.Filename, id, filename)

				if err := c.SaveUploadedFile(file, path); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusOK,
						"message": "upload file err: " + err.Error(),
						"success": false,
					})

					cleanup(filepath.Join(dst, id))

					return
				}
			}

			var returnUrl string
			entry := &DbEntry{Id: id, Expire: formdata.Expire}

			if len(uploaded) == 1 {
				returnUrl = cfg.Url + cfg.ApiPrefix + "/getfile/" + id + "/" + uploaded[0]
				entry.File = uploaded[0]
			} else {
				zipfile := Ts() + "data.zip"

				if err := zipSource(filepath.Join(dst, id), filepath.Join(dst, id, zipfile)); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": "upload file err: " + err.Error(),
						"success": false,
					})

					cleanup(filepath.Join(dst, id))
				}

				returnUrl = strings.Join([]string{cfg.Url + cfg.ApiPrefix, "getfile", id, zipfile}, "/")
				entry.File = zipfile

				// clean up after us
				go func() {
					for _, file := range uploaded {
						if err := os.Remove(filepath.Join(dst, id, file)); err != nil {
							Log("ERROR: unable to delete %s: %s", file, err)
						}
					}
				}()

			}

			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": returnUrl,
				"success": true,
			})

			Log("Now serving %s from %s/%s", returnUrl, dst, id)
			Log("Expire set to: %s", formdata.Expire)

			go func() {
				err := db.Update(func(tx *bolt.Tx) error {
					b, err := tx.CreateBucketIfNotExists([]byte("uploads"))
					if err != nil {
						return fmt.Errorf("create bucket: %s", err)
					}

					jsonentry, err := json.Marshal(entry)
					if err != nil {
						return fmt.Errorf("json marshalling failure: %s", err)
					}

					err = b.Put([]byte(id), []byte(jsonentry))
					if err != nil {
						return fmt.Errorf("insert data: %s", err)
					}

					// results in:
					// bbolt get /tmp/uploads.db uploads fb242922-86cb-43a8-92bc-b027c35f0586
					// {"id":"fb242922-86cb-43a8-92bc-b027c35f0586","expire":"1d","file":"2023-02-17-13-09-data.zip"}
					return nil
				})
				if err != nil {
					Log("DB error: %s", err.Error())
				}
			}()
		})

		api.GET("/getfile/:id/:file", func(c *gin.Context) {
			// deliver  a file and delete  it after a delay  (FIXME: check
			// when gin  has done delivering it?). Redirect  to the static
			// handler for actual delivery.
			id := c.Param("id")
			file := c.Param("file")
			c.Request.URL.Path = cfg.ApiPrefix + "/static/" + id + "/" + file
			filename := filepath.Join(dst, id, file)

			if _, err := os.Stat(filename); err == nil {
				go func() {
					time.Sleep(500 * time.Millisecond)
					cleanup(filepath.Join(dst, id))
				}()
			}

			router.HandleContext(c)
		})

		// actual  delivery of static  files, uri's  must be known  to the
		// user, mostly being redirected here internally from /f
		api.Static("/static", dst)
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome to upload api, use /api enpoint!")
	})

	router.Run(cfg.Listen)

	return nil
}

// cleanup an upload directory, either because  we got an error in the
// middle of an upload or something else  went wrong. we fork off a go
// routine because this may block.
func cleanup(dir string) {
	go func() {
		err := os.RemoveAll(dir)
		if err != nil {
			Log("Failed to remove dir %s: %s", dir, err.Error())
		}
	}()
}

func zipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
