package main

import (
	"fmt"
	wrapper "github.com/72nd/gopdf-wrapper"
	"github.com/72nd/gopdf-wrapper/fonts"
	"github.com/gofiber/fiber/v2"
	"github.com/signintech/gopdf"
)

func renderPdf(name string) string {
	doc, err := wrapper.NewDoc(28, 1)

	if err != nil {
		fmt.Print(err)
	}
	liberation, err := fonts.NewLiberationSansFamily()
	if err != nil {
		fmt.Print(err)
	}
	doc.SetFontFamily(*liberation)
	doc.AddPage()

	// Text
	doc.CellWithOption(&gopdf.Rect{
		W: 220,
		H: 30,
	}, "Привет, "+name, gopdf.CellOption{Align: gopdf.Center})
	doc.AddWrapText(10, 80, 200, "Lorem ipsum dolores sit amet")
	doc.WritePdf("hello.pdf")
	return "hello.pdf"
}

func main() {
	app := fiber.New()

	api := app.Group("/api", nil)
	pdfApi := api.Group("/gdf", nil)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	pdfApi.Get("otkr/:sample/:fio/:gender/:image", func(c *fiber.Ctx) error {
		return c.SendFile(renderPdf(c.Params("fio")))
	})

	app.Listen(":8000")
}
