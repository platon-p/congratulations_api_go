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
const background_path = "./backgrounds"

type Preset struct {
	name         string
	background   string
	greetingText string
	greetingY    int
	text         string
	textX        int
	textY        int
}

var gramota = Preset{
	name:         "Благодарственное письмо",
	background:   background_path + "/gramota.jpg",
	greetingText: "Здравствуйте, {}!",
	greetingY:    50,
	text:         "Laboris ad esse ullamco consectetur eiusmod do nisi nisi sunt dolor in incididunt et veniam do. Laboris ad esse ullamco consectetur eiusmod do nisi nisi sunt dolor in incididunt et veniam do.",
	textX:        40,
	textY:        70,
}

func renderPdf(name, greeting string) string {
	doc, err := gopdfwrapper.NewDoc(18, 1)
	if err != nil {
		fmt.Print(err)
	}

	liberation, _ := fonts.NewLiberationSansFamily()
	doc.AddPage()

	doc.SetFontFamily(*liberation)

	doc.Image("gramota.jpg", 0, 0, gopdf.PageSizeA4)

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
	//greeting := c.Query("greeting")
	fname := proPdf(name, gramota)
	return c.SendFile(fname)
}

func proPdf(name string, cType Preset) string {
	doc, err := gopdfwrapper.NewDoc(18, 1)
	if err != nil {
		fmt.Print(err)
	}

	liberation, _ := fonts.NewLiberationSansFamily()
	doc.SetFontFamily(*liberation)
	doc.AddPage()

	doc.Image(cType.background, 0, 0, gopdf.PageSizeA4)

	doc.SetY(float64(cType.greetingY))
	doc.SetFontStyle("bold")
	greeting := strings.Replace(cType.greetingText, "{}", "%v", 1)
	doc.CellWithOption(gopdf.PageSizeA4, fmt.Sprintf(greeting, name), gopdf.CellOption{Align: gopdf.Center})

	doc.SetFontStyle("")

	//w := 420 - cType.textX
	//symbW, _ := doc.MeasureTextWidth("Ж")
	//onOneLineSymbs := int(float64(w) / symbW)
	lineN := 0.
	lastSt := ""

	for _, st := range strings.Split(cType.text, " ") {
		l, _ := doc.MeasureTextWidth(lastSt + st + " ")
		if l <= float64(210-cType.textX*2+1) {
			lastSt += st + " "
		} else {
			doc.SetX(float64(cType.textX))
			doc.SetY(float64(cType.textY) + lineN*(doc.LineHeight(18)+1))
			doc.Text(lastSt)

			lineN++
			lastSt = st + " "
		}
		//
		//if len(lastSt+st+" ") <= onOneLineSymbs+1 {
		//	lastSt += st + " "
		//} else {
		//	doc.AddText(float64(cType.textX), float64(cType.textY)+lineN*doc.LineHeight(18), lastSt)
		//	lineN++
		//	lastSt = st + " "
		//}
	}
	doc.SetX(float64(cType.textX))
	doc.SetY(float64(cType.textY) + lineN*(doc.LineHeight(18)+1))
	doc.Text(lastSt)
	//doc.AddWrapText(float64(cType.textX), float64(cType.textY), 210, cType.text)

	pathName := "results/" + generateName() + ".pdf"
	doc.WritePdf(pathName)
	return pathName
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("./templates/index.html", fiber.Map{})
	})

	api := app.Group("/api")
	api.Get("/pdf", sendPdf)

	app.Listen(":8000")
}
