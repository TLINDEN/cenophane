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
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tlinden/up/upd/cfg"

	"os"
	"path/filepath"
	"time"
)

func FilePut(c *fiber.Ctx, cfg *cfg.Config, db *Db) (string, error) {
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
	form, err := c.MultipartForm()
	if err != nil {
		Log("multipart error %s", err.Error())
		return "", err
	}

	// init upload obj
	entry := &Upload{Id: id, Uploaded: Timestamp{Time: time.Now()}}

	// retrieve files, if any
	files := form.File["upload[]"]
	members, err := SaveFormFiles(c, cfg, files, id)
	if err != nil {
		return "", fiber.NewError(500, "Could not store uploaded file[s]!")
	}
	entry.Members = members

	// extract auxilliary form data (expire field et al)
	if err := c.BodyParser(&formdata); err != nil {
		Log("bodyparser error %s", err.Error())
		return "", err
	}

	// post process expire
	if len(formdata.Expire) == 0 {
		entry.Expire = "asap"
	} else {
		ex, err := Untaint(formdata.Expire, `[^dhms0-9]`) // duration or asap allowed
		if err != nil {
			return "", err
		}
		entry.Expire = ex
	}

	// get url [and zip if there are multiple files]
	returnUrl, Newfilename, err := ProcessFormFiles(cfg, entry.Members, id)
	if err != nil {
		return "", fiber.NewError(500, "Could not process uploaded file[s]!")
	}
	entry.File = Newfilename

	Log("Now serving %s from %s/%s", returnUrl, cfg.StorageDir, id)
	Log("Expire set to: %s", entry.Expire)

	// we do this in the background to not thwart the server
	go db.Insert(id, entry)

	return returnUrl, nil
}

func FileGet(c *fiber.Ctx, cfg *cfg.Config, db *Db) error {
	// deliver  a file and delete  it if expire is set to asap

	// we ignore c.Params("file"), cause  it may be malign. Also we've
	// got it in the db anyway
	id, err := Untaint(c.Params("id"), `[^a-zA-Z0-9\-]`)
	if err != nil {
		return fiber.NewError(403, "Invalid id provided!")
	}

	upload, err := db.Lookup(id)
	if err != nil {
		// non existent db entry with that id, or other db error, see logs
		return fiber.NewError(404, "No download with that id could be found!")
	}

	file := upload.File
	filename := filepath.Join(cfg.StorageDir, id, file)

	if _, err := os.Stat(filename); err != nil {
		// db entry is there, but file isn't (anymore?)
		go db.Delete(id)
		return fiber.NewError(404, "No download with that id could be found!")
	}

	// finally put the file to the client
	err = c.Download(filename, file)

	go func() {
		// check if we need to delete the file now and do it in the background
		if upload.Expire == "asap" {
			cleanup(filepath.Join(cfg.StorageDir, id))
			db.Delete(id)
		}
	}()

	return err
}

type Id struct {
	Id string `json:"name" xml:"name" form:"name"`
}

func FileDelete(c *fiber.Ctx, cfg *cfg.Config, db *Db) error {
	// delete file, id dir and db entry

	id, err := Untaint(c.Params("id"), `[^a-zA-Z0-9\-]`)
	if err != nil {
		return fiber.NewError(403, "Invalid id provided!")
	}

	if len(id) == 0 {
		return fiber.NewError(403, "No id given!")
	}

	cleanup(filepath.Join(cfg.StorageDir, id))

	err = db.Delete(id)
	if err != nil {
		// non existent db entry with that id, or other db error, see logs
		return fiber.NewError(404, "No upload with that id could be found!")
	}

	return nil
}
