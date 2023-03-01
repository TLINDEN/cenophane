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
)

func Runserver(cfg *cfg.Config, args []string) error {
	router := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		Immutable:     true,
		Prefork:       cfg.Prefork,
		ServerHeader:  "upd",
		AppName:       cfg.AppName,
		BodyLimit:     cfg.BodyLimit,
		Network:       cfg.Network,
	})

	router.Use(requestid.New())
	router.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	db, err := NewDb(cfg.DbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	api := router.Group(cfg.ApiPrefix + ApiVersion)
	{
		api.Post("/file/", func(c *fiber.Ctx) error {
			msg, err := FilePut(c, cfg, db)
			return SendResponse(c, msg, err)
		})

		api.Get("/file/:id/:file", func(c *fiber.Ctx) error {
			return FileGet(c, cfg, db)
		})

		api.Get("/file/:id/", func(c *fiber.Ctx) error {
			return FileGet(c, cfg, db)
		})

		api.Delete("/file/:id/", func(c *fiber.Ctx) error {
			return FileDelete(c, cfg, db)
		})

		api.Delete("/file/", func(c *fiber.Ctx) error {
			return FileDelete(c, cfg, db)
		})
	}

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("welcome to upload api, use /api enpoint!"))
	})

	return router.Listen(cfg.Listen)

}

func SendResponse(c *fiber.Ctx, msg string, err error) error {
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Result{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(Result{
		Code:    fiber.StatusOK,
		Message: msg,
		Success: true,
	})
}
