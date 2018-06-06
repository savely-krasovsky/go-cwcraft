package main

import (
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
	"sync"
)

var (
	equipment []equipmentItem
	alchemy   []alchemyItem
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
	// Init Configurator, Logger, Database, CW API, resources
	if err := Init(); err != nil {
		panic(err)
	}

	// API Responses
	go func() {
		for update := range client.Updates {
			if err := HandleUpdate(update); err != nil {
				log.Warn(err)
			}
		}
	}()

	// Yellow pages
	go func() {
		for pages := range client.YellowPages {
			if err := HandlePages(pages); err != nil {
				log.Warn(err)
			}
		}
	}()

	// Update all user stocks
	go func() {
		err := UpdateStocks()
		log.Warn(err)
	}()

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
	e.GET("/shops", Shops)

	e.GET("/login", LoginGet)
	e.POST("/login", LoginPost)

	e.GET("/stock", Stock)

	e.GET("/api/equipment", getEquipment)
	e.GET("/api/equipment/:id", getEquipmentItem)

	e.GET("/api/alchemy", getAlchemy)
	e.GET("/api/alchemy/:id", getAlchemyItem)

	e.GET("/api/resources", getResources)
	e.GET("/api/resources/:id", getResource)

	e.GET("/api/basics/:type/:id", getBasics)
	e.GET("/api/commands/:type/:id", getCommands)

	e.GET("/api/shops", getShops)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", viper.GetString("port"))))
}
