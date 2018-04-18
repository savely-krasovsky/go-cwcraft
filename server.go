package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
	"io/ioutil"
)

type (
	item struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		Stats     stats    `json:"stats"`
		Type      string   `json:"type"`
		ManaCost  int      `json:"mana_cost,omitempty"`
		Composite bool     `json:"composite"`
		Recipe    map[string]int `json:"recipe,omitempty"`
	}

	stats struct {
		Attack  int `json:"attack,omitempty"`
		Defense int `json:"defense,omitempty"`
		Mana    int `json:"mana,omitempty"`
	}

	resource struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		ManaCost  int      `json:"mana_cost,omitempty"`
		Composite bool     `json:"composite"`
		Recipe    map[string]int `json:"recipe,omitempty"`
	}

	command struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Amount          int    `json:"amount"`
		CommandManaCost int    `json:"command_mana_cost"`
	}

	basic struct {
		Name string `json:"name"`
		Amount int `json:"amount"`
	}
)

var (
	items     []item
	resources []resource
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Read items
	b, err := ioutil.ReadFile("res/items.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal items
	err = json.Unmarshal(b, &items)
	if err != nil {
		log.Fatal(err)
	}

	// Read resources
	b, err = ioutil.ReadFile("res/resources.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal resources
	err = json.Unmarshal(b, &resources)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Static("/", "static")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// Routes
	e.GET("/", Index)
	e.GET("/resources", Resources)
	e.GET("/alchemist", Alchemist)

	e.GET("/api/items", getItems)
	e.GET("/api/items/:id", getItem)

	e.GET("/api/resources", getResources)
	e.GET("/api/resources/:id", getResource)

	e.GET("/api/basics/:id", getBasics)
	e.GET("/api/commands/:id", getCommands)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
