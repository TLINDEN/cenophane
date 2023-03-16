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
	"github.com/tlinden/up/upd/cfg"
	//"github.com/alecthomas/repr"
	bolt "go.etcd.io/bbolt"
)

const Bucket string = "uploads"

// wrapper for bolt db
type Db struct {
	bolt *bolt.DB
	cfg  *cfg.Config
}

func NewDb(c *cfg.Config) (*Db, error) {
	b, err := bolt.Open(c.DbFile, 0600, nil)
	db := Db{bolt: b, cfg: c}
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

func (db *Db) Delete(apicontext string, id string) error {
	err := db.bolt.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))

		if bucket == nil {
			return fmt.Errorf("id %s not found", id)
		}

		j := bucket.Get([]byte(id))

		if len(j) == 0 {
			return fmt.Errorf("id %s not found", id)
		}

		upload := &Upload{}
		if err := json.Unmarshal(j, &upload); err != nil {
			return fmt.Errorf("unable to unmarshal json: %s", err)
		}

		if (apicontext != "" && (db.cfg.Super == apicontext || upload.Context == apicontext)) || apicontext == "" {
			return bucket.Delete([]byte(id))
		}

		return nil
	})

	if err != nil {
		Log("DB error: %s", err.Error())
	}

	return err
}

func (db *Db) List(apicontext string, filter string) (*Uploads, error) {
	uploads := &Uploads{}

	err := db.bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))
		if bucket == nil {
			return nil
		}

		err := bucket.ForEach(func(id, j []byte) error {
			upload := &Upload{}
			if err := json.Unmarshal(j, &upload); err != nil {
				return fmt.Errorf("unable to unmarshal json: %s", err)
			}

			fmt.Printf("apicontext: %s, filter: %s\n", apicontext, filter)
			if apicontext != "" && db.cfg.Super != apicontext {
				// only return the uploads for this context
				if apicontext == upload.Context {
					// unless a filter needed OR no filter specified
					if (filter != "" && upload.Context == filter) || filter == "" {
						uploads.Entries = append(uploads.Entries, upload)
					}
				}
			} else {
				// return all, because we operate a public service or current==super
				if (filter != "" && upload.Context == filter) || filter == "" {
					uploads.Entries = append(uploads.Entries, upload)
				}
			}

			return nil
		})

		return err // might be nil as well
	})

	return uploads, err
}

// we only return one obj here, but could return more later
func (db *Db) Get(apicontext string, id string) (*Uploads, error) {
	uploads := &Uploads{}

	err := db.bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))
		if bucket == nil {
			return nil
		}

		j := bucket.Get([]byte(id))
		if j == nil {
			return fmt.Errorf("No upload object found with id %s", id)
		}

		upload := &Upload{}
		if err := json.Unmarshal(j, &upload); err != nil {
			return fmt.Errorf("unable to unmarshal json: %s", err)
		}

		if (apicontext != "" && (db.cfg.Super == apicontext || upload.Context == apicontext)) || apicontext == "" {
			// allowed if no context (public or download)
			// or if context matches or if context==super
			uploads.Entries = append(uploads.Entries, upload)
		}

		return nil
	})

	return uploads, err
}

// a wrapper around Lookup() which extracts the 1st upload, if any
func (db *Db) Lookup(apicontext string, id string) (*Upload, error) {
	uploads, err := db.Get(apicontext, id)

	if err != nil {
		// non existent db entry with that id, or other db error, see logs
		return &Upload{}, fmt.Errorf("No upload object found with id %s", id)
	}

	if len(uploads.Entries) == 0 {
		return &Upload{}, fmt.Errorf("No upload object found with id %s", id)
	}

	return uploads.Entries[0], nil
}
