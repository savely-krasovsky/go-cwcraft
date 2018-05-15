package main

import (
	"fmt"
	"github.com/L11R/go-chatwars-api"
	"github.com/arangodb/go-driver"
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
		equipmentItem
		Recipe        []basic
		Basics        []basic
		Purchases     []basic
		Commands      []command
		TotalManaCost int
		ShowUserData  bool
	}

	var extItems []extendedItem

	for _, e := range equipment {
		basics := RecurBasics(e.Recipe)
		basics = SplitBasics(basics)

		commands := RecurCommands(e.Recipe)
		commands = SplitCommands(commands)

		// make recipe to add user amount field
		var recipe []basic
		for name, amount := range e.Recipe {
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

		var purchases []basic

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

			purchases = RecurPurchases(e.Recipe, user.Stock)
			purchases = SplitPurchases(purchases)

			commands = RecurUserCommands(e.Recipe, user.Stock)
			commands = SplitUserCommands(commands)
		}

		// add craft itself
		commands = append(commands, command{
			e.ID,
			e.Name,
			1,
			e.ManaCost,
		})

		extItem := extendedItem{
			e,
			recipe,
			basics,
			purchases,
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
		sort.Slice(extItem.Purchases, func(i, j int) bool { return extItem.Purchases[i].Name < extItem.Purchases[j].Name })

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
		Purchases     []basic
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

		var purchases []basic

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

			purchases = RecurPurchases(r.Recipe, user.Stock)
			purchases = SplitPurchases(purchases)

			commands = RecurUserCommands(r.Recipe, user.Stock)
			commands = SplitUserCommands(commands)
		}

		// add craft itself
		commands = append(commands, command{
			r.ID,
			r.Name,
			1,
			r.ManaCost,
		})

		extItem := extendedItem{
			r,
			recipe,
			basics,
			purchases,
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
		sort.Slice(extItem.Purchases, func(i, j int) bool { return extItem.Purchases[i].Name < extItem.Purchases[j].Name })

		extItems = append(extItems, extItem)
	}

	return c.Render(http.StatusOK, "resources", extItems)
}

func Alchemist(c echo.Context) error {
	sess, _ := session.Get("user", c)

	var user user
	if id, found := sess.Values["id"]; found {
		_, err := usersCol.ReadDocument(nil, fmt.Sprint(id), &user)
		if err != nil {
			return echo.ErrNotFound
		}
	}

	type extendedItem struct {
		alchemyItem
		Recipe        []basic
		Basics        []basic
		Purchases     []basic
		Commands      []command
		TotalManaCost int
		ShowUserData  bool
	}

	var extItems []extendedItem

	for _, a := range alchemy {
		basics := RecurBasics(a.Recipe)
		basics = SplitBasics(basics)

		commands := RecurCommands(a.Recipe)
		commands = SplitCommands(commands)

		// make recipe to add user amount field
		var recipe []basic
		for name, amount := range a.Recipe {
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

		var purchases []basic

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

			purchases = RecurPurchases(a.Recipe, user.Stock)
			purchases = SplitPurchases(purchases)

			commands = RecurUserCommands(a.Recipe, user.Stock)
			commands = SplitUserCommands(commands)
		}

		// add craft itself
		commands = append(commands, command{
			a.ID,
			a.Name,
			1,
			a.ManaCost,
		})

		extItem := extendedItem{
			a,
			recipe,
			basics,
			purchases,
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
		sort.Slice(extItem.Purchases, func(i, j int) bool { return extItem.Purchases[i].Name < extItem.Purchases[j].Name })

		extItems = append(extItems, extItem)
	}

	return c.Render(http.StatusOK, "alchemist", extItems)
}

func Shops(c echo.Context) error {
	// Get tokens cursor
	cursor, err := db.Query(
		nil,
		`FOR s IN shops
			RETURN s`,
		nil,
	)
	if err != nil {
		return err
	}

	// Don't forget to close
	defer cursor.Close()

	var shops []cwapi.YellowPage

	// Get all tokens
	for {
		var shop cwapi.YellowPage
		_, err := cursor.ReadDocument(nil, &shop)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return err
		}

		shops = append(shops, shop)
	}

	return c.Render(http.StatusOK, "shops", shops)
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

func getEquipment(c echo.Context) error {
	v, err := c.FormParams()
	if err != nil {
		return echo.ErrMethodNotAllowed
	}

	name := v.Get("name")
	itemType := v.Get("type")

	if name != "" || itemType != "" {
		for _, e := range equipment {
			if e.Name == name {
				return c.JSON(http.StatusOK, e)
			}

			if e.Type == itemType {
				return c.JSON(http.StatusOK, e)
			}
		}
	}

	return c.JSON(http.StatusOK, equipment)
}

func getEquipmentItem(c echo.Context) error {
	id := c.Param("id")

	for _, e := range equipment {
		if e.ID == id {
			return c.JSON(http.StatusOK, e)
		}
	}

	return echo.ErrNotFound
}

func getAlchemy(c echo.Context) error {
	v, err := c.FormParams()
	if err != nil {
		return echo.ErrMethodNotAllowed
	}

	name := v.Get("name")
	itemType := v.Get("type")

	if name != "" || itemType != "" {
		for _, a := range alchemy {
			if a.Name == name {
				return c.JSON(http.StatusOK, a)
			}

			if a.Type == itemType {
				return c.JSON(http.StatusOK, a)
			}
		}
	}

	return c.JSON(http.StatusOK, equipment)
}

func getAlchemyItem(c echo.Context) error {
	id := c.Param("id")

	for _, a := range alchemy {
		if a.ID == id {
			return c.JSON(http.StatusOK, a)
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
		Item   interface{} `json:"item"`
		Basics []basic     `json:"basics"`
	}

	id := c.Param("id")
	itemType := c.Param("type")

	switch itemType {
	case "equipment":
		for _, e := range equipment {
			if e.ID == id {
				basics := RecurBasics(e.Recipe)
				basics = SplitBasics(basics)

				return c.JSON(http.StatusOK, response{
					e,
					basics,
				})
			}
		}
	case "alchemy":
		for _, a := range alchemy {
			if a.ID == id {
				basics := RecurBasics(a.Recipe)
				basics = SplitBasics(basics)

				return c.JSON(http.StatusOK, response{
					a,
					basics,
				})
			}
		}
	}

	return echo.ErrNotFound
}

func getCommands(c echo.Context) error {
	type response struct {
		Item          interface{} `json:"item"`
		Commands      []command   `json:"commands"`
		TotalManaCost int         `json:"total_mana_cost"`
	}

	id := c.Param("id")
	itemType := c.Param("type")

	switch itemType {
	case "equipment":
		for _, e := range equipment {
			if e.ID == id {
				commands := RecurCommands(e.Recipe)
				commands = SplitCommands(commands)

				// add craft itself
				commands = append(commands, command{
					e.ID,
					e.Name,
					1,
					e.ManaCost,
				})

				res := response{
					e,
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
	case "alchemy":
		for _, a := range alchemy {
			if a.ID == id {
				commands := RecurCommands(a.Recipe)
				commands = SplitCommands(commands)

				// add craft itself
				commands = append(commands, command{
					a.ID,
					a.Name,
					1,
					a.ManaCost,
				})

				res := response{
					a,
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
	}

	return echo.ErrNotFound
}

func getShops(c echo.Context) error {
	// Get tokens cursor
	cursor, err := db.Query(
		nil,
		`FOR s IN shops
			RETURN s`,
		nil,
	)
	if err != nil {
		return err
	}

	// Don't forget to close
	defer cursor.Close()

	var shops []cwapi.YellowPage

	// Get all tokens
	for {
		var shop cwapi.YellowPage
		_, err := cursor.ReadDocument(nil, &shop)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return err
		}

		shops = append(shops, shop)
	}

	return c.JSON(http.StatusOK, shops)
}
