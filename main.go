package main

import (
	"encoding/json"
	"fmt"
	"github.com/72nd/gopdf-wrapper"
	"github.com/72nd/gopdf-wrapper/fonts"
	"github.com/gofiber/fiber/v2"
	"github.com/signintech/gopdf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

const (
	letters        = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	backgroundPath = "./backgrounds"
)

type Preset struct {
	gorm.Model
	Name      string  `json:"Name"`
	PaperSize string  `json:"PaperSize"`
	Text      string  `json:"Text"`
	Greeting  string  `json:"Greeting"`
	TextX     float64 `json:"TextX"`
	TextY     float64 `json:"TextY"`
	GreetingY float64 `json:"GreetingY"`
	Image     string  `json:"Image"`
}

var gramota = Preset{
	Name:      "Благодарственное письмо",
	PaperSize: "A4",
	Text:      "Ea tempor in pariatur ea enim nulla eiusmod. Nulla fugiat consequat occaecat est id consectetur Lorem voluptate ut amet sunt tempor. Nulla fugiat consequat occaecat est id consectetur Lorem voluptate ut amet sunt tempor.",
	Greeting:  "[Уважаемый][Уважаемая] {}!",
	TextX:     30,
	TextY:     60,
	GreetingY: 30,
	Image:     backgroundPath + "/gramota.jpg",
}
var db *gorm.DB

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
	_ = doc.SetFontFamily(*liberation)
	doc.AddPage()

	_ = doc.Image(cType.Image, 0, 0, gopdf.PageSizeA4)

	doc.SetY(cType.GreetingY)
	_ = doc.SetFontStyle("bold")
	greeting := strings.Replace(cType.Greeting, "{}", "%v", 1)
	_ = doc.CellWithOption(gopdf.PageSizeA4, fmt.Sprintf(greeting, name), gopdf.CellOption{Align: gopdf.Center})

	_ = doc.SetFontStyle("")

	lineN := 0.
	lastSt := ""

	for _, st := range strings.Split(cType.Text, " ") {
		l, _ := doc.MeasureTextWidth(lastSt + st + " ")
		if l <= 210-cType.TextX*2+1 {
			lastSt += st + " "
		} else {
			doc.SetX(cType.TextX)
			doc.SetY(cType.TextY + lineN*(doc.LineHeight(18)+1))
			_ = doc.Text(lastSt)

			lineN++
			lastSt = st + " "
		}
	}
	doc.SetX(cType.TextX)
	doc.SetY(cType.TextY + lineN*(doc.LineHeight(18)+1))
	_ = doc.Text(lastSt)

	pathName := "results/" + generateName() + ".pdf"
	_ = doc.WritePdf(pathName)
	return pathName
}

func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Connection error")
	}
	db.AutoMigrate()
}

func newPreset(c *fiber.Ctx) error {
	var preset Preset
	err := json.Unmarshal(c.Body(), &preset)
	if err != nil {
		return fiber.NewError(400, "Can not parse json")
	}
	db.Create(&preset)
	return c.Redirect("/api/presets")
}

func getPreset(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return fiber.NewError(400, "Id field is required")
	}
	var preset Preset
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(400, "Could not parse id")
	}
	db.First(&preset, idInt)
	presetJson, _ := json.Marshal(preset)
	if err != nil {
		return fiber.NewError(500, "Could not create JSON")
	}
	return c.Send(presetJson)
}

func updatePreset(c *fiber.Ctx) error {
	id := c.Query("id")
	param := strings.Split(strings.Split(c.OriginalURL(), "&")[1], "=")

	db.Model(&Preset{}).Where("ID = "+id).Update(param[0], c.Query(param[0]))
	return c.Redirect("preset?id=" + id)
}

func deletePreset(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return fiber.NewError(400, "Id field is required")
	}
	var preset Preset
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(400, "Could not parse id")
	}
	db.First(&preset, idInt)
	db.Delete(&preset)
	return fiber.NewError(200, "Success")
}

func getPresets(c *fiber.Ctx) error {
	var presets []Preset
	db.Find(&presets)
	presetsBody, _ := json.Marshal(presets)
	return c.Send(presetsBody)
}

func main() {
	initDatabase()

	app := fiber.New()
	app.Static("/backgrounds", backgroundPath)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("./templates/index.html", fiber.Map{})
	})

	api := app.Group("/api")
	api.Get("/pdf", sendPdf)
	api.Post("/new_preset", newPreset)
	api.Get("/preset", getPreset)
	api.Get("/update_preset", updatePreset)
	api.Get("/delete_preset", deletePreset)
	api.Get("/presets", getPresets)

	err := app.Listen(":8000")
	if err != nil {
		log.Fatalln("Oops...")
	}
}
