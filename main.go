package main

import (
	"fmt"
	"github.com/72nd/gopdf-wrapper"
	"github.com/72nd/gopdf-wrapper/fonts"
	"github.com/gofiber/fiber/v2"
	"github.com/signintech/gopdf"
	"math/rand"
	"strings"
)

const letters = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

var gramota = Preset{
	name:       "Грамота",
	background: "background.jpg",
	greetingY:  0,
	textX:      0,
	textY:      0,
}

type Preset struct {
	name       string
	background string
	greetingY  int
	textX      int
	textY      int
}

func renderPdf(name, greeting string) string {
	doc, err := gopdfwrapper.NewDoc(18, 1)
	if err != nil {
		fmt.Print(err)
	}

	liberation, _ := fonts.NewLiberationSansFamily()
	doc.AddPage()

	doc.SetFontFamily(*liberation)

	doc.Image("background.jpg", 0, 0, gopdf.PageSizeA4)

	doc.SetY(50)

	greeting = strings.Replace(greeting, "{}", "%v", 1)
	doc.SetX(0)
	doc.CellWithOption(gopdf.PageSizeA4, fmt.Sprintf(greeting, name), gopdf.CellOption{Align: gopdf.Center})

	pathName := "results/" + generateName() + ".pdf"
	doc.WritePdf(pathName)
	return pathName
}

func generateName() string {
	s := ""
	for i := 0; i < 9; i++ {
		s += string(letters[rand.Int()%len(letters)])
	}
	return s
}

func sendPdf(c *fiber.Ctx) error {
	name := c.Query("name")
	greeting := c.Query("greeting")
	fname := renderPdf(name, greeting)
	return c.SendFile(fname)
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	api := app.Group("/api")
	api.Get("/pdf", sendPdf)

	app.Listen(":8000")
}
