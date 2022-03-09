package main

import (
	"fmt"
	wrapper "github.com/72nd/gopdf-wrapper"
	"github.com/72nd/gopdf-wrapper/fonts"
	"github.com/gofiber/fiber/v2"
	"github.com/signintech/gopdf"
)

func renderPdf() string {
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
	}, "HELLO", gopdf.CellOption{Align: gopdf.Center})
	doc.AddWrapText(10, 80, 200, "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem.")
	doc.WritePdf("hello.pdf")
	//
	//a, _ := fonts.NewLatoFamily()
	//doc.SetFontFamily(*a)
	//doc.AddPage()
	//
	//doc.AddWrapText(10, 10, gopdf.PageSizeA5.W-10, "Hello world asdasdasdasd Hello world asdasdasdasdHello world asdasdasdasd Hello world asdasdasdasd Hello world asdasdasdasd")
	//
	//doc.Close()
	//
	//err := doc.WritePdf("hello.pdf")
	//if err != nil {
	//	fmt.Print(err)
	//}
	return ""
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	app.Get("otkr/:sample/:fio/:gender/:image", func(c *fiber.Ctx) error {
		renderPdf()
		return c.SendFile("hello.pdf")
		//return c.SendString(fmt.Sprintf("%v, %v", c.Params("sample"), c.Params("fio")))
	})

	app.Listen(":8000")
}
