package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

// Embed a directory
//go:embed public/*
var embedDirStatic embed.FS

func main() {
	debugValue, hasDev := os.LookupEnv("DEBUG")
	isProd := !hasDev || debugValue == ""
	app := fiber.New(fiber.Config{
		Prefork:               isProd,
		ReduceMemoryUsage:     isProd,
		DisableStartupMessage: isProd,
		ServerHeader:          "Stream::IO",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Retreive the custom statuscode if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code := e.Code
				log.Print(fmt.Sprintln("Code %d", code))
				if code != 404 {

					// Send custom error page
					err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
					if err != nil {
						// In case the SendFile fails
						return ctx.Status(500).SendString("Internal Server Error")
					}
				}
			}
			log.Println("Redirect")
			// Return from handler
			ctx.Redirect("/")
			return nil
		},
	})

	subFS, _ := fs.Sub(embedDirStatic, "public")
	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(subFS),
	}))

	log.Fatal(app.Listen(":3000"))
}
