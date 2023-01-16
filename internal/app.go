package internal

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/clickadilla/cache-house/internal/controllers"
	"github.com/clickadilla/cache-house/internal/controllers/api/public"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/clickadilla/cache-house/internal/middleware"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type App struct {
	router            *echo.Echo
	homeController    *controllers.HomeController
	feedsController   *public.FeedsController
	sspController     *public.SupplySidePlatformsController
	networkController *public.NetworkController
	discrepController *public.DiscrepancyController
	logger            *logrus.Logger
}

func (a *App) Run() {

	a.boot()
	a.bootRouting()

	log.Fatal(a.router.Start(":" + os.Getenv("APP_PORT")))
}

func (a *App) boot() {
	ctx := context.Background()

	a.logger = logrus.New()
	a.logger.SetOutput(os.Stdout)

	err := godotenv.Load()
	if err != nil {
		a.logger.Fatal("Error loading .env file")
	}

	clickadillaClient := managers.NewClickadillaClient(os.Getenv("CLICKADILLA_API_ENDPOINT"), os.Getenv("CLICKADILLA_API_TOKEN"))
	feedState := managers.NewFeedState(clickadillaClient, a.logger)
	go feedState.RunUpdate()
	go feedState.RunUpdateAllFeeds()
	go feedState.RunUpdateFeedsNetworks()
	go feedState.RunUpdateFeedsAccountManagers()

	sspState := managers.NewSupplySidePlatformState(clickadillaClient, a.logger)
	go sspState.RunUpdate()

	networkState := managers.NewNetworkState(clickadillaClient, a.logger)
	networkStateUpdater := managers.NewUpdater(networkState, time.Minute*2)
	networkStateUpdater.StartPeriodicUpdate(ctx)

	discrepancyState := managers.NewDiscrepancyState(clickadillaClient, a.logger, feedState)
	go discrepancyState.RunUpdate()

	a.router = echo.New()
	a.feedsController = &public.FeedsController{
		FeedState: feedState,
	}
	a.sspController = &public.SupplySidePlatformsController{
		SupplySidePlatformState: sspState,
	}
	a.networkController = &public.NetworkController{
		NetworkState: networkState,
	}
	a.discrepController = &public.DiscrepancyController{
		DiscrepancyState: discrepancyState,
	}
	a.homeController = &controllers.HomeController{
		FeedState:    feedState,
		SspState:     sspState,
		DiscrepState: discrepancyState,
	}
}

func (a *App) bootRouting() {
	a.router.GET("/", a.homeController.Index)
	a.router.GET("/api/feeds", a.feedsController.Index, middleware.Auth)
	a.router.GET("/api/feeds/tsv", a.feedsController.FeedListTsv, middleware.Auth)
	a.router.GET("/api/feeds/list-account/tsv", a.feedsController.ListAccountTsv, middleware.Auth)
	a.router.GET("/api/feeds/list-network/tsv", a.feedsController.ListNetworkTsv, middleware.Auth)
	a.router.GET("/api/feeds/account-managers/tsv", a.feedsController.FeedsAccountManagers, middleware.Auth)

	a.router.GET("/api/supply-side-platforms", a.sspController.Index, middleware.Auth)

	a.router.GET("/api/networks", a.networkController.Index, middleware.Auth)
	a.router.GET("/api/networks/tsv", a.networkController.Tsv, middleware.Auth)

	a.router.GET("/api/discrepancies", a.discrepController.Index, middleware.Auth)
}
