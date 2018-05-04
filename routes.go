package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/gommon/log"
	"net/http"
	"sort"
	"time"
)

func Index(c echo.Context) error {
	sess, _ := session.Get("user", c)

	var user user
	if id, found := sess.Values["id"]; found {
		_, err := usersCol.ReadDocument(nil, fmt.Sprint(id), &user)
		if err != nil {
			return echo.ErrNotFound
		}
	}

	type extendedItem struct {
		item
		Recipe        []basic
		Basics        []basic
		Commands      []command
		TotalManaCost int
		ShowUserData  bool
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

		// make recipe to add user amount field
		var recipe []basic
		for name, amount := range i.Recipe {
			if userAmount, found := user.Stock[name]; found {
				recipe = append(recipe, basic{
					Name:       name,
					Amount:     amount,
					UserAmount: userAmount,
				})
			} else {
				recipe = append(recipe, basic{
					Name:   name,
					Amount: amount,
				})
			}
		}

		showUserData := false
		if user.ID != "" {
			showUserData = true

			// add user amount to basics
			for name, amount := range user.Stock {
				for i := range basics {
					if basics[i].Name == name {
						basics[i].UserAmount = amount
					}
				}
			}
		}

		extItem := extendedItem{
			i,
			recipe,
			basics,
			commands,
			0,
			showUserData,
		}

		// count total mana cost
		for _, c := range commands {
			extItem.TotalManaCost += c.CommandManaCost
		}

		// Sort recipe and basics to fix it
		sort.Slice(extItem.Recipe, func(i, j int) bool { return extItem.Recipe[i].Name < extItem.Recipe[j].Name })
		sort.Slice(extItem.Basics, func(i, j int) bool { return extItem.Basics[i].Name < extItem.Basics[j].Name })

		extItems = append(extItems, extItem)
	}

	return c.Render(http.StatusOK, "index", extItems)
}

func Resources(c echo.Context) error {
	sess, _ := session.Get("user", c)

	var user user
	if id, found := sess.Values["id"]; found {
		_, err := usersCol.ReadDocument(nil, fmt.Sprint(id), &user)
		if err != nil {
			return echo.ErrNotFound
		}
	}

	type extendedItem struct {
		resource
		Recipe        []basic
		Basics        []basic
		Commands      []command
		TotalManaCost int
		ShowUserData  bool
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

		// make recipe to add user amount field
		var recipe []basic
		for name, amount := range r.Recipe {
			if userAmount, found := user.Stock[name]; found {
				recipe = append(recipe, basic{
					Name:       name,
					Amount:     amount,
					UserAmount: userAmount,
				})
			} else {
				recipe = append(recipe, basic{
					Name:   name,
					Amount: amount,
				})
			}
		}

		showUserData := false
		if user.ID != "" {
			showUserData = true

			// add user amount to basics
			for name, amount := range user.Stock {
				for i := range basics {
					if basics[i].Name == name {
						basics[i].UserAmount = amount
					}
				}
			}
		}

		extItem := extendedItem{
			r,
			recipe,
			basics,
			commands,
			0,
			showUserData,
		}

		// count total mana cost
		for _, c := range commands {
			extItem.TotalManaCost += c.CommandManaCost
		}

		// Sort recipe and basics to fix it
		sort.Slice(extItem.Recipe, func(i, j int) bool { return extItem.Recipe[i].Name < extItem.Recipe[j].Name })
		sort.Slice(extItem.Basics, func(i, j int) bool { return extItem.Basics[i].Name < extItem.Basics[j].Name })

		extItems = append(extItems, extItem)
	}

	return c.Render(http.StatusOK, "resources", extItems)
}

func Alchemist(c echo.Context) error {
	return c.String(http.StatusOK, "Coming soon!")
}

func LoginGet(c echo.Context) error {
	l := login{
		Status: "unknown",
	}

	sess, _ := session.Get("user", c)
	if _, found := sess.Values["id"]; found {
		l.Status = "alreadyLogged"
	}

	return c.Render(http.StatusOK, "login", l)
}

func LoginPost(c echo.Context) error {
	l := new(login)
	if err := c.Bind(l); err != nil {
		log.Error(err)
		l.Status = "internalError"
	}

	if l.ID == 0 {
		l.Status = "idNotSpecified"
	} else if l.Code == "" {
		if err := client.CreateAuthCode(l.ID); err != nil {
			log.Error(err)
			l.Status = "internalError"
		} else {
			// create waiter chan and save it in waiters
			waiter := make(chan map[string]string, 1)
			waiters.Store(l.ID, waiter)

			select {
			// wait auth from updates handler
			case status := <-waiter:
				if s, found := status["createAuthCode"]; found {
					l.Status = s
				}

				waiters.Delete(l.ID)
				// or timeout
			case <-time.After(5 * time.Second):
				l.Status = "timeout"
				waiters.Delete(l.ID)
			}
		}
	} else {
		if err := client.GrantToken(l.ID, l.Code); err != nil {
			log.Error(err)
			l.Status = "internalError"
		} else {
			waiter := make(chan map[string]string, 1)
			waiters.Store(l.ID, waiter)

			select {
			// wait auth from updates handler
			case status := <-waiter:
				if s, found := status["grantToken"]; found {
					if s == "success" {
						// save cookie with ID
						sess, _ := session.Get("user", c)
						sess.Values["id"] = fmt.Sprint(l.ID)
						if err := sess.Save(c.Request(), c.Response()); err != nil {
							log.Error(err)
							l.Status = "internalError"
						}

						// request stock after login
						// Request new stock
						if err := client.RequestStock(fmt.Sprint(l.ID)); err != nil {
							log.Error(err)
							l.Status = "internalError"
						}

						l.Status = s
					} else {
						l.Status = "internalError"
					}
				}

				waiters.Delete(l.ID)
				// or timeout
			case <-time.After(5 * time.Second):
				l.Status = "timeout"
				waiters.Delete(l.ID)
			}
		}
	}

	if l.Status != "success" {
		sess, _ := session.Get("user", c)
		if _, found := sess.Values["id"]; found {
			l.Status = "alreadyLogged"
		}
	}

	return c.Render(http.StatusOK, "login", l)
}

func Stock(c echo.Context) error {
	sess, _ := session.Get("user", c)

	var user user
	_, err := usersCol.ReadDocument(nil, fmt.Sprint(sess.Values["id"]), &user)
	if err != nil {
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "stock", user)
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
	type response struct {
		Item   item    `json:"item"`
		Basics []basic `json:"basics"`
	}

	id := c.Param("id")

	for _, i := range items {
		if i.ID == id {
			basics := RecurBasics(i.Recipe)
			basics = SplitBasics(basics)

			return c.JSON(http.StatusOK, response{
				i,
				basics,
			})
		}
	}

	return echo.ErrNotFound
}

func getCommands(c echo.Context) error {
	type response struct {
		Item          item      `json:"item"`
		Commands      []command `json:"commands"`
		TotalManaCost int       `json:"total_mana_cost"`
	}

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

			res := response{
				i,
				commands,
				0,
			}

			// count total mana cost
			for _, c := range commands {
				res.TotalManaCost += c.CommandManaCost
			}

			return c.JSON(http.StatusOK, res)
		}
	}

	return echo.ErrNotFound
}
