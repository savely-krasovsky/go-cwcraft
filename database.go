package main

import (
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/spf13/viper"
)

func InitConnectionPool() error {
	var err error

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{viper.GetString("database.address")},
	})
	if err != nil {
		return err
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
		Authentication: driver.BasicAuthentication(
			viper.GetString("database.username"),
			viper.GetString("database.password"),
		),
	})
	if err != nil {
		return err
	}

	db, err = client.Database(nil, viper.GetString("database.name"))
	if err != nil {
		return err
	}

	usersCol, err = db.Collection(nil, "users")
	if err != nil {
		return err
	}

	return nil
}
