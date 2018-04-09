//go:generate abigen --sol ./contracts/MocrowCoinCrowdsale.sol --pkg gocontracts  --out ./gocontracts/contracts.go

package main

import (
	"fmt"

	"MCW-btc-module/model"
	"MCW-btc-module/server"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sevlyar/go-daemon"
	config "github.com/spf13/viper"
)

func main() {
	fmt.Println("GET CONFIG")

	config.AddConfigPath(".")
	config.SetConfigName("config")

	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error getting config from file: %s \n", err))
	}
	if config.GetBool("daemon.enabled") {

		cntxt := &daemon.Context{
			PidFileName: config.GetString("daemon.pidfile"),
			PidFilePerm: 0644,
			LogFileName: config.GetString("daemon.logfile"),
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
		}

		d, err := cntxt.Reborn()
		if err != nil {
			panic("Unable to run: " + err.Error())
		}
		if d != nil {
			return
		}
		defer cntxt.Release()

		fmt.Println("- - - - - - - - - - - - - - -")
		fmt.Println("daemon started")
	}
	fmt.Println("CONNECT TO DATABASE")
	database, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s user=%s dbname=%s password=%s",
			config.GetString("postgres.host"),
			config.GetString("postgres.user"),
			config.GetString("postgres.dbname"),
			config.GetString("postgres.password"),
		),
	)

	if err != nil {
		panic("COULDN'T CONNECT TO DATABASE " + err.Error())
	}

	fmt.Println("BEGIN MIGRATIONS")
	database.AutoMigrate(&model.BTCTransaction{})
	fmt.Println("END MIGRATIONS")

	server, err := server.New(database)

	if err != nil {
		panic(err)
	}

	server.Logger.Fatal(server.Start(":4000"))
}
