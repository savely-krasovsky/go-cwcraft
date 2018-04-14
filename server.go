package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"net/http"
	"io/ioutil"
	"github.com/labstack/gommon/log"
	"encoding/json"
)



type (
	quickMap = map[string]int

	item struct {
		ID       string         `json:"id"`
		Name     string         `json:"name"`
		Stats    stats          `json:"stats"`
		Type     string         `json:"type"`
		ManaCost int            `json:"mana_cost,omitempty"`
		Recipe   quickMap `json:"recipe,omitempty"`
	}

	stats struct {
		Attack  int `json:"attack,omitempty"`
		Defense int `json:"defense,omitempty"`
		Mana    int `json:"mana,omitempty"`
	}

	resource struct {
		ID       string         `json:"id"`
		Name     string         `json:"name"`
		ManaCost int            `json:"mana_cost,omitempty"`
		Recipe   quickMap `json:"recipe,omitempty"`
	}

	command struct {
		ID string `json:"id"`
		Name string `json:"name"`
		Amount int `json:"amount"`
	}
)

var (
	items []item
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
		*item
		Basics  quickMap
		Commands []command
	}

	var extItem []extendedItem



	for _, item := range items {
		basics, commands := getBasicsRecursively(quickMap{}, &[]command{}, item.Recipe)

		// Don't forget to reverse array
		for i, j := 0, len(commands)-1; i < j; i, j = i+1, j-1 {
			commands[i], commands[j] = commands[j], commands[i]
		}

		extItem = append(extItem, extendedItem{
			&item,
			basics,
			commands,
		})
	}

	return c.Render(http.StatusOK, "index", extItem)
}

func getItems(c echo.Context) error {
	v, err := c.FormParams()
	if err != nil {
		return echo.ErrMethodNotAllowed
	}

	name := v.Get("name")
	itemType := v.Get("type")

	if name != "" || itemType != "" {
		for _, item := range items {
			if item.Name == name {
				return c.JSON(http.StatusOK, item)
			}

			if item.Type == itemType {
				return c.JSON(http.StatusOK, item)
			}
		}
	}

	return c.JSON(http.StatusOK, items)
}

func getItem(c echo.Context) error {
	id := c.Param("id")

	for _, item := range items {
		if item.ID == id {
			return c.JSON(http.StatusOK, item)
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
		for _, resource := range resources {
			if resource.Name == name {
				return c.JSON(http.StatusOK, resource)
			}
		}
	}

	return c.JSON(http.StatusOK, resources)
}

func getResource(c echo.Context) error {
	id := c.Param("id")

	for _, resource := range resources {
		if resource.ID == id {
			return c.JSON(http.StatusOK, resource)
		}
	}

	return echo.ErrNotFound
}

func getBasics(c echo.Context) error {
	id := c.Param("id")

	for _, item := range items {
		if item.ID == id {
			basics, _ := getBasicsRecursively(quickMap{}, &[]command{}, item.Recipe)

			return c.JSON(http.StatusOK, basics)
		}
	}

	return echo.ErrNotFound
}

func getCommands(c echo.Context) error {
	id := c.Param("id")

	for _, item := range items {
		if item.ID == id {
			_, commands := getBasicsRecursively(quickMap{}, &[]command{}, item.Recipe)

			// Don't forget to reverse array
			for i, j := 0, len(commands)-1; i < j; i, j = i+1, j-1 {
				commands[i], commands[j] = commands[j], commands[i]
			}

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
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// Routes
	e.GET("/", Index)

	e.GET("/items", getItems)
	e.GET("/items/:id", getItem)

	e.GET("/resources", getResources)
	e.GET("/resources/:id", getResource)

	e.GET("/basics/:id", getBasics)
	e.GET("/commands/:id", getCommands)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
