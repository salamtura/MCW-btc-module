package server

import (
	"MCW-btc-module/controllers"
	"MCW-btc-module/routing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "github.com/spf13/viper"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	*echo.Echo
}

func New(database *gorm.DB) (*Server, error) {
	server := Server{echo.New()}

	// Middleware
	server.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	//server.Pre(middleware.HTTPSRedirect())
	server.Pre(middleware.RemoveTrailingSlash())
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.CORS())

	mainGroup := server.Group("/exchange")

	blockcypherController := controllers.MakeBlockCypherController(
		config.GetString("blockcypher.accessToken"),
		config.GetBool("blockcypher.isTestnet"),
	)

	blocktrailController := controllers.MakeBlocktrailController(
		config.GetString("blocktrail.apiKey"),
		config.GetBool("blocktrail.isTestnet"),
	)

	monitoringController := controllers.MakeMonitoringController(
		blockcypherController,
		blocktrailController,
	)

	infuraController := controllers.MakeInfuraController(
		config.GetString("infura.accessToken"),
		config.GetBool("infura.isTestnet"),
	)

	tokenManagementController, err := controllers.MakeTokenManagementController(
		infuraController,
		config.GetString("crowdsale.address"),
		config.GetString("crowdsale.ownerAddress"),
		config.GetString("crowdsale.ownerPrivateKey"),
	)

	if err != nil {
		return nil, err
	}

	exchangeController, err := controllers.MakeExchangeController(monitoringController, *tokenManagementController, database, config.GetString("bitcoin.xPub"))

	if err != nil {
		return nil, err
	}

	whitelistController, err := controllers.MakeWhitelistController(infuraController, config.GetString("crowdsale.address"))

	exchangeRouter := routing.MakeExchangeRouter(exchangeController, whitelistController)

	exchangeRouter.Register(mainGroup)

	return &server, nil
}
