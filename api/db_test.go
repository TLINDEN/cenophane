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
	"github.com/tlinden/cenophane/cfg"
	"os"
	"testing"
)

func finalize(db *Db) {
	if db.bolt != nil {
		db.Close()
	}
	if _, err := os.Stat(db.cfg.DbFile); err == nil {
		os.Remove(db.cfg.DbFile)
	}
}

func TestNew(t *testing.T) {
	var tests = []struct {
		name     string
		file     string
		wantfail bool
	}{
		{"opennew", "test.db", false},
		{"openfail", "/hopefully/not/existing/directory/test.db", true},
	}

	for _, tt := range tests {
		c := &cfg.Config{DbFile: tt.file}
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDb(c)
			defer finalize(db)
			if err != nil && !tt.wantfail {
				t.Errorf("expected: &Db{}, got err: " + err.Error())
			}

			if err == nil && tt.wantfail {
				t.Errorf("expected: fail, got &Db{}")
			}
		})
	}
}
