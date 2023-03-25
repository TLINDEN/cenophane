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

package api

import (
	//"github.com/alecthomas/repr"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tlinden/cenophane/cfg"
	"github.com/tlinden/cenophane/common"

	"os"
	"path/filepath"
	"strings"
	"time"
)

type SetContext struct {
	Apicontext string `json:"apicontext" form:"apicontext"`
}

func UploadPost(c *fiber.Ctx, cfg *cfg.Config, db *Db) error {
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

	id := uuid.NewString()

	var returnUrl string
	var formdata Meta

	os.MkdirAll(filepath.Join(cfg.StorageDir, id), os.ModePerm)

	// fetch auxiliary form data
	form, err := c.MultipartForm()
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"mime/multipart error "+err.Error())
	}

	// init upload obj
	entry := &common.Upload{Id: id, Uploaded: common.Timestamp{Time: time.Now()}}

	// retrieve the API Context name from the session
	apicontext, err := GetApicontext(c)
	if err != nil {
		return JsonStatus(c, fiber.StatusInternalServerError,
			"Unable to initialize session store from context: "+err.Error())
	}
	entry.Context = apicontext

	// retrieve files, if any
	files := form.File["upload[]"]
	members, err := SaveFormFiles(c, cfg, files, id)
	if err != nil {
		return JsonStatus(c, fiber.StatusInternalServerError,
			"Could not store uploaded file[s]: "+err.Error())
	}
	entry.Members = members

	// extract auxilliary form data (expire field et al)
	if err := c.BodyParser(&formdata); err != nil {
		return JsonStatus(c, fiber.StatusInternalServerError,
			"bodyparser error : "+err.Error())
	}

	// post process expire
	if len(formdata.Expire) == 0 {
		entry.Expire = "asap"
	} else {
		ex, err := common.Untaint(formdata.Expire, cfg.RegDuration) // duration or asap allowed
		if err != nil {
			return JsonStatus(c, fiber.StatusForbidden,
				"Invalid data: "+err.Error())
		}
		entry.Expire = ex
	}

	// get url [and zip if there are multiple files]
	returnUrl, Newfilename, err := ProcessFormFiles(cfg, entry.Members, id)
	if err != nil {
		return JsonStatus(c, fiber.StatusInternalServerError,
			"Could not process uploaded file[s]: "+err.Error())
	}
	entry.File = Newfilename
	entry.Url = returnUrl

	Log("Now serving %s from %s/%s", returnUrl, cfg.StorageDir, id)
	Log("Expire set to: %s", entry.Expire)
	Log("Uploaded with API-Context %s", entry.Context)

	// we do this in the background to not thwart the server
	go db.Insert(id, entry)

	// everything went well so far
	res := &common.Response{Uploads: []*common.Upload{entry}}
	res.Success = true
	res.Code = fiber.StatusOK
	return c.Status(fiber.StatusOK).JSON(res)
}

func UploadFetch(c *fiber.Ctx, cfg *cfg.Config, db *Db, shallExpire ...bool) error {
	// deliver  a file and delete  it if expire is set to asap

	// we ignore c.Params("file"), cause  it may be malign. Also we've
	// got it in the db anyway
	id, err := common.Untaint(c.Params("id"), cfg.RegKey)
	if err != nil {
		return fiber.NewError(403, "Invalid id provided!")
	}

	// retrieve the API Context name from the session
	apicontext, err := GetApicontext(c)
	if err != nil {
		return JsonStatus(c, fiber.StatusInternalServerError,
			"Unable to initialize session store from context: "+err.Error())
	}

	response, err := db.Lookup(apicontext, id, common.TypeUpload)
	if err != nil {
		// non existent db entry with that id, or other db error, see logs
		return fiber.NewError(404, "No download with that id could be found!")
	}

	var upload *common.Upload
	if len(response.Uploads) > 0 {
		upload = response.Uploads[0]
	}

	file := upload.File
	filename := filepath.Join(cfg.StorageDir, id, file)

	if _, err := os.Stat(filename); err != nil {
		// db entry is there, but file isn't (anymore?)
		go db.Delete(apicontext, id)
		return fiber.NewError(404, "No download with that id could be found!")
	}

	// finally put the file to the client
	err = c.Download(filename, file)

	if len(shallExpire) > 0 {
		if shallExpire[0] == true {
			go func() {
				// check if we need to delete the file now and do it in the background
				if upload.Expire == "asap" {
					cleanup(filepath.Join(cfg.StorageDir, id))
					db.Delete(apicontext, id)
				}
			}()
		}
	}

	return err
}

// delete file, id dir and db entry
func UploadDelete(c *fiber.Ctx, cfg *cfg.Config, db *Db) error {

	id, err := common.Untaint(c.Params("id"), cfg.RegKey)
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"Invalid id provided!")
	}

	if len(id) == 0 {
		return JsonStatus(c, fiber.StatusForbidden,
			"No id specified!")
	}

	// retrieve the API Context name from the session
	apicontext, err := GetApicontext(c)
	if err != nil {
		return JsonStatus(c, fiber.StatusInternalServerError,
			"Unable to initialize session store from context: "+err.Error())
	}

	err = db.Delete(apicontext, id)
	if err != nil {
		// non existent db entry with that id, or other db error, see logs
		return JsonStatus(c, fiber.StatusForbidden,
			"No upload with that id could be found!")
	}

	cleanup(filepath.Join(cfg.StorageDir, id))

	return nil
}

// returns the whole list + error code, no post processing by server
func UploadsList(c *fiber.Ctx, cfg *cfg.Config, db *Db) error {
	// fetch filter from body(json expected)
	setcontext := new(SetContext)
	if err := c.BodyParser(setcontext); err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"Unable to parse body: "+err.Error())
	}

	filter, err := common.Untaint(setcontext.Apicontext, cfg.RegKey)
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"Invalid api context filter provided!")
	}

	// retrieve the API Context name from the session
	apicontext, err := GetApicontext(c)
	if err != nil {
		return JsonStatus(c, fiber.StatusInternalServerError,
			"Unable to initialize session store from context: "+err.Error())
	}

	// get list
	uploads, err := db.UploadsList(apicontext, filter, common.TypeUpload)
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"Unable to list uploads: "+err.Error())
	}

	// if we reached this point we can signal success
	uploads.Success = true
	uploads.Code = fiber.StatusOK

	return c.Status(fiber.StatusOK).JSON(uploads)
}

// returns just one upload obj + error code, no post processing by server
func UploadDescribe(c *fiber.Ctx, cfg *cfg.Config, db *Db) error {
	id, err := common.Untaint(c.Params("id"), cfg.RegKey)
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"Invalid id provided!")
	}

	// retrieve the API Context name from the session
	apicontext, err := GetApicontext(c)
	if err != nil {
		return JsonStatus(c, fiber.StatusInternalServerError,
			"Unable to initialize session store from context: "+err.Error())
	}

	response, err := db.Get(apicontext, id, common.TypeUpload)
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"No upload with that id could be found!")
	}

	for _, upload := range response.Uploads {
		upload.Url = strings.Join([]string{cfg.Url, "download", id, upload.File}, "/")
	}

	// if we reached this point we can signal success
	response.Success = true
	response.Code = fiber.StatusOK

	return c.Status(fiber.StatusOK).JSON(response)
}
