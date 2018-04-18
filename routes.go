package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func Index(c echo.Context) error {
	type extendedItem struct {
		item
		Basics        []basic
		Commands      []command
		TotalManaCost int
	}

	var extItems []extendedItem

	for _, i := range items {
		basics := RecurBasics(i.Recipe)
		basics = SplitBasics(basics)

		commands := RecurCommands(i.Recipe)
		commands = SplitCommands(commands)

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
		Basics        []basic
		Commands      []command
		TotalManaCost int
	}

	var extItems []extendedItem

	for _, r := range resources {
		// skip basic
		if r.Composite == false {
			continue
		}

		basics := RecurBasics(r.Recipe)
		basics = SplitBasics(basics)

		commands := RecurCommands(r.Recipe)
		commands = SplitCommands(commands)

		// add craft itself
		commands = append(commands, command{
			r.ID,
			r.Name,
			1,
			r.ManaCost,
		})

		extItem := extendedItem{
			r,
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
			basics := RecurBasics(i.Recipe)
			basics = SplitBasics(basics)

			return c.JSON(http.StatusOK, basics)
		}
	}

	return echo.ErrNotFound
}

func getCommands(c echo.Context) error {
	id := c.Param("id")

	for _, i := range items {
		if i.ID == id {
			commands := RecurCommands(i.Recipe)
			commands = SplitCommands(commands)

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