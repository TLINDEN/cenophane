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
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/keyauth/v2"
	"github.com/tlinden/up/upd/cfg"
)

// sessions are context specific and can be global savely
var Sessionstore *session.Store

const shallExpire = true

func Runserver(conf *cfg.Config, args []string) error {
	// required for authenticated routes, used to store the api context
	Sessionstore = session.New()

	// bbolt db setup
	db, err := NewDb(conf)
	if err != nil {
		return err
	}
	defer db.Close()

	// setup authenticated endpoints
	auth := SetupAuthStore(conf)

	// setup api server
	router := SetupServer(conf)

	// authenticated routes
	api := router.Group(conf.ApiPrefix + ApiVersion)
	{
		// authenticated routes
		api.Post("/file/", auth, func(c *fiber.Ctx) error {
			return FilePut(c, conf, db)
		})

		api.Get("/file/:id/:file", auth, func(c *fiber.Ctx) error {
			return FileGet(c, conf, db)
		})

		api.Get("/file/:id/", auth, func(c *fiber.Ctx) error {
			return FileGet(c, conf, db)
		})

		api.Delete("/file/:id/", auth, func(c *fiber.Ctx) error {
			err := DeleteUpload(c, conf, db)
			return SendResponse(c, "", err)
		})

		api.Get("/list/", auth, func(c *fiber.Ctx) error {
			return List(c, conf, db)
		})

		api.Get("/upload/:id/", auth, func(c *fiber.Ctx) error {
			return Describe(c, conf, db)
		})
	}

	// public routes
	{
		router.Get("/", func(c *fiber.Ctx) error {
			return c.Send([]byte("welcome to upload api, use /api enpoint!"))
		})

		router.Get("/download/:id/:file", func(c *fiber.Ctx) error {
			return FileGet(c, conf, db, shallExpire)
		})

		router.Get("/download/:id/", func(c *fiber.Ctx) error {
			return FileGet(c, conf, db, shallExpire)
		})
	}

	// setup cleaner
	quitcleaner := BackgroundCleaner(conf, db)

	router.Hooks().OnShutdown(func() error {
		Log("Shutting down cleaner")
		close(quitcleaner)
		return nil
	})

	return router.Listen(conf.Listen)
}

func SetupAuthStore(conf *cfg.Config) func(*fiber.Ctx) error {
	AuthSetEndpoints(conf.ApiPrefix, ApiVersion, []string{"/file"})
	AuthSetApikeys(conf.Apicontext)

	return keyauth.New(keyauth.Config{
		Validator:    AuthValidateAPIKey,
		ErrorHandler: AuthErrHandler,
	})
}

func SetupServer(conf *cfg.Config) *fiber.App {
	router := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		Immutable:     true,
		Prefork:       conf.Prefork,
		ServerHeader:  "upd",
		AppName:       conf.AppName,
		BodyLimit:     conf.BodyLimit,
		Network:       conf.Network,
	})

	router.Use(requestid.New())

	router.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	router.Use(cors.New(cors.Config{
		AllowMethods:  "GET,PUT,POST,DELETE",
		ExposeHeaders: "Content-Type,Authorization,Accept",
	}))

	router.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	return router
}

/*
   Wrapper to respond with proper json status, message and code,
   shall be prepared and called by the handlers directly.
*/
func JsonStatus(c *fiber.Ctx, code int, msg string) error {
	success := true

	if code != fiber.StatusOK {
		success = false
	}

	return c.Status(code).JSON(Result{
		Code:    code,
		Message: msg,
		Success: success,
	})
}

/*
   Used for non json-aware handlers, called by server
*/
func SendResponse(c *fiber.Ctx, msg string, err error) error {
	if err != nil {
		code := fiber.StatusInternalServerError

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		return c.Status(code).JSON(Result{
			Code:    code,
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
