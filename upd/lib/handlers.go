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
	//"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tlinden/up/upd/cfg"
	bolt "go.etcd.io/bbolt"
	//"io"
	// "net/http"
	"os"
	"path/filepath"
	//"regexp"
	"strings"
	"time"
)

func Putfile(c *gin.Context, cfg *cfg.Config, db *bolt.DB) (string, error) {
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

	var returnUrl string
	var formdata Meta

	os.MkdirAll(filepath.Join(cfg.StorageDir, id), os.ModePerm)

	// fetch auxiliary form data
	form, _ := c.MultipartForm()

	entry := &Upload{Id: id, Uploaded: time.Now()}

	// init upload obj

	// retrieve files, if any
	files := form.File["upload[]"]
	for _, file := range files {
		filename := NormalizeFilename(filepath.Base(file.Filename))
		path := filepath.Join(cfg.StorageDir, id, filename)
		entry.Members = append(entry.Members, filename)
		Log("Received: %s => %s/%s", file.Filename, id, filename)

		if err := c.SaveUploadedFile(file, path); err != nil {
			cleanup(filepath.Join(cfg.StorageDir, id))
			return "", err
		}
	}

	if err := c.ShouldBind(&formdata); err != nil {
		return "", err
	}
	if len(formdata.Expire) == 0 {
		entry.Expire = "asap"
	} else {
		entry.Expire = formdata.Expire // FIXME: validate
	}

	if len(entry.Members) == 1 {
		returnUrl = cfg.Url + cfg.ApiPrefix + "/getfile/" + id + "/" + entry.Members[0]
		entry.File = entry.Members[0]
	} else {
		zipfile := Ts() + "data.zip"
		tmpzip := filepath.Join(cfg.StorageDir, zipfile)
		finalzip := filepath.Join(cfg.StorageDir, id, zipfile)
		iddir := filepath.Join(cfg.StorageDir, id)

		if err := zipSource(iddir, tmpzip); err != nil {
			cleanup(iddir)
			return "", err
		}

		if err := os.Rename(tmpzip, finalzip); err != nil {
			cleanup(iddir)
			return "", err
		}

		returnUrl = strings.Join([]string{cfg.Url + cfg.ApiPrefix, "getfile", id, zipfile}, "/")
		entry.File = zipfile

		// clean up after us
		go func() {
			for _, file := range entry.Members {
				if err := os.Remove(filepath.Join(cfg.StorageDir, id, file)); err != nil {
					Log("ERROR: unable to delete %s: %s", file, err)
				}
			}
		}()

	}

	Log("Now serving %s from %s/%s", returnUrl, cfg.StorageDir, id)
	Log("Expire set to: %s", entry.Expire)

	go func() {
		// => db.go !
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

	return returnUrl, nil
}
