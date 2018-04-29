package main

import (
	"fmt"
	"github.com/L11R/go-chatwars-api"
	"github.com/arangodb/go-driver"
	"github.com/labstack/gommon/log"
	"time"
)

func HandleUpdate(update cwapi.Response) error {
	log.Info(update)

	switch update.Action {
	case "createAuthCode":
		if update.Result == "Ok" {
			_, err := usersCol.CreateDocument(nil, user{
				ID: fmt.Sprint(update.ParsedPayload.(cwapi.ResCreateAuthCode).UserID),
			})
			if err != nil && err.(driver.ArangoError).ErrorNum == 1210 {
				// pass it
			} else if err != nil {
				return err
			}
		}

		if waiter, found := waiters.Load(update.ParsedPayload.(cwapi.ResCreateAuthCode).UserID); found {
			// found? send it to waiter channel
			if update.Result == "Ok" {
				waiter.(chan map[string]string) <- map[string]string{"createAuthCode": "waitingCode"}
			} else {
				waiter.(chan map[string]string) <- map[string]string{"createAuthCode": "internalError"}
			}

			// trying to prevent memory leak
			close(waiter.(chan map[string]string))
		}
	case "grantToken":
		if update.Result == "Ok" {
			_, err := usersCol.UpdateDocument(
				nil,
				fmt.Sprint(update.ParsedPayload.(cwapi.ResGrantToken).UserID),
				user{
					Token: update.ParsedPayload.(cwapi.ResGrantToken).Token,
				},
			)
			if err != nil {
				return err
			}
		}

		if waiter, found := waiters.Load(update.ParsedPayload.(cwapi.ResGrantToken).UserID); found {
			// found? send it to waiter channel
			if update.Result == "Ok" {
				waiter.(chan map[string]string) <- map[string]string{"grantToken": "success"}
			} else {
				waiter.(chan map[string]string) <- map[string]string{"grantToken": "internalError"}
			}

			// trying to prevent memory leak
			close(waiter.(chan map[string]string))
		}
	case "requestStock":
		if update.Result == "Ok" {
			_, err := usersCol.UpdateDocument(
				nil,
				fmt.Sprint(update.ParsedPayload.(cwapi.ResRequestStock).UserID),
				user{
					Stock: update.ParsedPayload.(cwapi.ResRequestStock).Stock,
				},
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateStocks() error {
	for {
		// Get tokens cursor
		cursor, err := db.Query(
			nil,
			`FOR u IN users
						RETURN u`,
			nil,
		)
		if err != nil {
			return err
		}

		// Don't forget to close
		defer cursor.Close()

		// Get all tokens
		for {
			var user user
			_, err := cursor.ReadDocument(nil, &user)
			if driver.IsNoMoreDocuments(err) {
				break
			} else if err != nil {
				return err
			}

			// Request new stock
			if err := client.RequestStock(user.Token); err != nil {
				return err
			}
		}

		// Wait before new update
		time.Sleep(30 * time.Second)
	}
}
