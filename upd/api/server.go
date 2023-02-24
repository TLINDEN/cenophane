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
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/tlinden/up/upd/cfg"
	bolt "go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"time"
)

func Runserver(cfg *cfg.Config, args []string) error {
	dst := cfg.StorageDir
	router := fiber.New()
	router.Use(requestid.New())
	router.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))
	api := router.Group(cfg.ApiPrefix + ApiVersion)

	db, err := bolt.Open(cfg.DbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	{
		api.Post("/file/put", func(c *fiber.Ctx) error {
			uri, err := Putfile(c, cfg, db)

			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(Result{
					Code:    fiber.StatusBadRequest,
					Message: err.Error(),
					Success: false,
				})
			} else {
				return c.Status(fiber.StatusOK).JSON(Result{
					Code:    fiber.StatusOK,
					Message: uri,
					Success: true,
				})
			}
		})

		api.Get("/file/get/:id/:file", func(c *fiber.Ctx) error {
			// deliver  a file and delete  it after a delay  (FIXME: check
			// when gin  has done delivering it?). Redirect  to the static
			// handler for actual delivery.
			id := c.Params("id")
			file := c.Params("file")

			filename := filepath.Join(dst, id, file)

			if _, err := os.Stat(filename); err == nil {
				go func() {
					time.Sleep(500 * time.Millisecond)
					cleanup(filepath.Join(dst, id))
				}()
			}

			return c.Download(filename, file)
		})
	}

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("welcome to upload api, use /api enpoint!"))
	})

	router.Listen(cfg.Listen)

	return nil
}
