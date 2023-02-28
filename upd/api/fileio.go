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
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// cleanup an upload directory, either because  we got an error in the
// middle of an upload or something else  went wrong.
func cleanup(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		Log("Failed to remove dir %s: %s", dir, err.Error())
	}
}

func ZipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	// source must be an absolute path, target a zip file
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	// don't chdir the server itself
	go func() {
		defer wg.Done()

		os.Chdir(source)
		newDir, err := os.Getwd()
		if err != nil {
		}
		if newDir != source {
			err = errors.New("Failed to changedir!")
			return
		}

		err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			// 2. Go through all the files of the source
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
			header.Name = path
			//Log("a: <%s>, b: <%s>, rel: <%s>", filepath.Dir(source), path, header.Name)
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
	}()

	wg.Wait()

	return err
}
