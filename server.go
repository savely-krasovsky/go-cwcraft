package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

type (
	quickMap = map[string]int

	item struct {
		ID       string   `json:"id"`
		Name     string   `json:"name"`
		Stats    stats    `json:"stats"`
		Type     string   `json:"type"`
		ManaCost int      `json:"mana_cost,omitempty"`
		Recipe   quickMap `json:"recipe,omitempty"`
	}

	stats struct {
		Attack  int `json:"attack,omitempty"`
		Defense int `json:"defense,omitempty"`
		Mana    int `json:"mana,omitempty"`
	}

	resource struct {
		ID       string   `json:"id"`
		Name     string   `json:"name"`
		ManaCost int      `json:"mana_cost,omitempty"`
		Recipe   quickMap `json:"recipe,omitempty"`
	}

	command struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Amount int    `json:"amount"`
		CommandManaCost int `json:"command_mana_cost"`
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

func Index(c echo.Context) error {
	type extendedItem struct {
		item
		Basics   quickMap
		Commands []command
		TotalManaCost int
	}

	var extItems []extendedItem

	for _, i := range items {
		basics, commands := getBasicsRecursively(quickMap{}, &[]command{}, i.Recipe)

		// don't forget to reverse array
		for i, j := 0, len(commands)-1; i < j; i, j = i+1, j-1 {
			commands[i], commands[j] = commands[j], commands[i]
		}

		// add craft itself
		commands = append(commands, command{
			i.ID,
			i.Name,
			1,
			i.ManaCost,
		})

		extItem := extendedItem{
			i,
			basics,
			commands,
			0,
		}

		// count total mana cost
		for _, com := range commands {
			extItem.TotalManaCost += com.CommandManaCost
		}

		extItems = append(extItems, extItem)
	}

	return c.Render(http.StatusOK, "index", extItems)
}

func Resources(c echo.Context) error {
	type extendedItem struct {
		resource
		Basics   quickMap
		ShowBasics bool
		Commands []command
		TotalManaCost int
	}

	var extItems []extendedItem

	for _, r := range resources {
		// skip basic
		if len(r.Recipe) == 0 {
			continue
		}

		basics, commands := getBasicsRecursively(quickMap{}, &[]command{}, r.Recipe)

		// don't forget to reverse array
		for i, j := 0, len(commands)-1; i < j; i, j = i+1, j-1 {
			commands[i], commands[j] = commands[j], commands[i]
		}

		// add craft itself
		commands = append(commands, command{
			r.ID,
			r.Name,
			1,
			r.ManaCost,
		})

		// show Basics by default
		sb := true

		// if they equils hide
		if reflect.DeepEqual(r.Recipe, basics) {
			sb = false
		}

		extItem := extendedItem{
			r,
			basics,
			sb,
			commands,
			0,
		}

		// count total mana cost
		for _, com := range commands {
			extItem.TotalManaCost += com.CommandManaCost
		}

		extItems = append(extItems, extItem)
	}

	return c.Render(http.StatusOK, "resources", extItems)
}

func Alchemist(c echo.Context) error {
	return c.String(http.StatusOK, "Coming soon!")
}

func getItems(c echo.Context) error {
	v, err := c.FormParams()
	if err != nil {
		return echo.ErrMethodNotAllowed
	}

	name := v.Get("name")
	itemType := v.Get("type")

	if name != "" || itemType != "" {
		for _, i := range items {
			if i.Name == name {
				return c.JSON(http.StatusOK, i)
			}

			if i.Type == itemType {
				return c.JSON(http.StatusOK, i)
			}
		}
	}

	return c.JSON(http.StatusOK, items)
}

func getItem(c echo.Context) error {
	id := c.Param("id")

	for _, i := range items {
		if i.ID == id {
			return c.JSON(http.StatusOK, i)
		}
	}

	return echo.ErrNotFound
}

func getResources(c echo.Context) error {
	v, err := c.FormParams()
	if err != nil {
		return echo.ErrMethodNotAllowed
	}

	name := v.Get("name")

	if name != "" {
		for _, r := range resources {
			if r.Name == name {
				return c.JSON(http.StatusOK, r)
			}
		}
	}

	return c.JSON(http.StatusOK, resources)
}

func getResource(c echo.Context) error {
	id := c.Param("id")

	for _, r := range resources {
		if r.ID == id {
			return c.JSON(http.StatusOK, r)
		}
	}

	return echo.ErrNotFound
}

func getBasics(c echo.Context) error {
	id := c.Param("id")

	for _, i := range items {
		if i.ID == id {
			basics, _ := getBasicsRecursively(quickMap{}, &[]command{}, i.Recipe)
			return c.JSON(http.StatusOK, basics)
		}
	}

	return echo.ErrNotFound
}

func getCommands(c echo.Context) error {
	id := c.Param("id")

	for _, i := range items {
		if i.ID == id {
			_, commands := getBasicsRecursively(quickMap{}, &[]command{}, i.Recipe)

			// don't forget to reverse array
			for i, j := 0, len(commands)-1; i < j; i, j = i+1, j-1 {
				commands[i], commands[j] = commands[j], commands[i]
			}

			// add craft itself
			commands = append(commands, command{
				i.ID,
				i.Name,
				1,
				i.ManaCost,
			})

			return c.JSON(http.StatusOK, commands)
		}
	}

	return echo.ErrNotFound
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
