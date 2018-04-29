package main

import (
	"encoding/json"
	"fmt"
	"github.com/L11R/go-chatwars-api"
	"github.com/arangodb/go-driver"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"html/template"
	"io"
	"io/ioutil"
	"sync"
)

var (
	items     []item
	resources []resource

	client *cwapi.Client

	db       driver.Database
	usersCol driver.Collection

	waiters sync.Map
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Initialize CW API client
	client = cwapi.NewClient(viper.GetString("cwapi.username"), viper.GetString("cwapi.password"))

	// Log
	go func() {
		for update := range client.Updates {
			if err := HandleUpdate(update); err != nil {
				log.Error(err)
			}
		}
	}()

	// Database pool init
	if err := InitConnectionPool(); err != nil {
		log.Fatal(err)
	}

	// Update all user stocks
	go func() {
		err := UpdateStocks()
		log.Error(err)
	}()

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
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(viper.GetString("sessions_secret")))))

	// Set renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// Routes
	e.GET("/", Index)
	e.GET("/resources", Resources)
	e.GET("/alchemist", Alchemist)

	e.GET("/login", LoginGet)
	e.POST("/login", LoginPost)

	e.GET("/stock", Stock)

	e.GET("/api/items", getItems)
	e.GET("/api/items/:id", getItem)

	e.GET("/api/resources", getResources)
	e.GET("/api/resources/:id", getResource)

	e.GET("/api/basics/:id", getBasics)
	e.GET("/api/commands/:id", getCommands)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
