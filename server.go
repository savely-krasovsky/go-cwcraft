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
	item struct {
		ID       string         `json:"id"`
		Name     string         `json:"name"`
		Stats    stats          `json:"stats"`
		Type     string         `json:"type"`
		ManaCost int            `json:"mana_cost,omitempty"`
		Recipe   map[string]int `json:"recipe,omitempty"`
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
		Recipe   map[string]int `json:"recipe,omitempty"`
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
	return c.Render(http.StatusOK, "index", items)
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

func getBasic(c echo.Context) error {
	id := c.Param("id")

	for _, item := range items {
		if item.ID == id {
			res := getBasicRecursively(map[string]int{}, item.Recipe)
			return c.JSON(http.StatusOK, res)
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

	e.GET("/basic/:id", getBasic)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
