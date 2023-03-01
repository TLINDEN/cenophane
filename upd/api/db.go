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
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
)

const Bucket string = "uploads"

// wrapper for bolt db
type Db struct {
	bolt *bolt.DB
}

// stores 1 upload object, gets into db
type Upload struct {
	Id       string    `json:"id"`
	Expire   string    `json:"expire"`
	File     string    `json:"file"`    // final filename (visible to the downloader)
	Members  []string  `json:"members"` // contains multiple files, so File is an archive
	Uploaded Timestamp `json:"uploaded"`
}

func NewDb(file string) (*Db, error) {
	b, err := bolt.Open(file, 0600, nil)
	db := Db{bolt: b}
	return &db, err
}

func (db *Db) Close() {
	db.bolt.Close()
}

func (db *Db) Insert(id string, entry *Upload) error {
	err := db.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(Bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		jsonentry, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("json marshalling failure: %s", err)
		}

		err = bucket.Put([]byte(id), []byte(jsonentry))
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

	return err
}

func (db *Db) Lookup(id string) (Upload, error) {
	var upload Upload

	err := db.bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))
		j := bucket.Get([]byte(id))

		if len(j) == 0 {
			return fmt.Errorf("id %s not found", id)
		}

		if err := json.Unmarshal(j, &upload); err != nil {
			return fmt.Errorf("unable to unmarshal json: %s", err)
		}

		return nil
	})

	if err != nil {
		Log("DB error: %s", err.Error())
		return upload, err
	}

	return upload, nil
}

func (db *Db) Delete(id string) error {
	err := db.bolt.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))

		j := bucket.Get([]byte(id))

		if len(j) == 0 {
			return fmt.Errorf("id %s not found", id)
		}

		err := bucket.Delete([]byte(id))
		return err
	})

	if err != nil {
		Log("DB error: %s", err.Error())
	}

	return err
}
