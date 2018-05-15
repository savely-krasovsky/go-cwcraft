package main

import (
	"encoding/json"
	"github.com/L11R/go-chatwars-api"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
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

// Init configuration manager, logger, bot, database
func Init() error {
	// Init and read config file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Configuration defaults
	// Log level: INFO (-1 for DEBUG)
	viper.SetDefault("log.level", 0)
	// Log type: "production" or "development"
	viper.SetDefault("log.type", "production")

	// Init logger
	var loggerConfig zap.Config
	if viper.GetString("log.type") == "production" {
		loggerConfig = zap.NewProductionConfig()
	}
	if viper.GetString("log.type") == "development" {
		loggerConfig = zap.NewDevelopmentConfig()
	}
	loggerConfig.Level.SetLevel(zapcore.Level(viper.GetInt("log.level")))

	logger, _ := loggerConfig.Build()
	defer logger.Sync()

	sugar = logger.Sugar()

	// Init database
	err = InitConnectionPool()
	if err != nil {
		return err
	}

	// Init Chat Wars API client
	client, err = cwapi.NewClient(viper.GetString("cwapi.user"), viper.GetString("cwapi.password"))
	if err != nil {
		return err
	}

	err = client.InitYellowPages()
	if err != nil {
		return err
	}

	// Read equipment
	b, err := ioutil.ReadFile("res/equipment.json")
	if err != nil {
		return err
	}

	// Unmarshal equipment
	err = json.Unmarshal(b, &equipment)
	if err != nil {
		return err
	}

	// Read equipment
	b, err = ioutil.ReadFile("res/alchemy.json")
	if err != nil {
		return err
	}

	// Unmarshal equipment
	err = json.Unmarshal(b, &alchemy)
	if err != nil {
		return err
	}

	// Read resources
	b, err = ioutil.ReadFile("res/resources.json")
	if err != nil {
		return err
	}

	// Unmarshal resources
	err = json.Unmarshal(b, &resources)
	if err != nil {
		return err
	}

	return nil
}
