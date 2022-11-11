package internal

import (
	"github.com/clickadilla/cache-house/internal/controllers"
	"github.com/clickadilla/cache-house/internal/controllers/api/public"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/clickadilla/cache-house/internal/middleware"
	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

type App struct {
	router            *router.Router
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

	log.Fatal(fasthttp.ListenAndServe(":"+os.Getenv("APP_PORT"), a.router.Handler))
}

func (a *App) boot() {
	a.logger = logrus.New()
	a.logger.SetOutput(os.Stdout)

	err := godotenv.Load()
	if err != nil {
		a.logger.Fatal("Error loading .env file")
	}

	clickadillaClient := managers.NewClickadillaClient(os.Getenv("CLICKADILLA_API_ENDPOINT"))
	feedState := managers.NewFeedState(clickadillaClient, a.logger)
	go feedState.RunUpdate()

	sspState := managers.NewSupplySidePlatformState(clickadillaClient, a.logger)
	go sspState.RunUpdate()

	networkState := managers.NewNetworkState(clickadillaClient, a.logger)
	go networkState.RunUpdate()

	discrepancyState := managers.NewDiscrepancyState(clickadillaClient, a.logger, feedState)
	go discrepancyState.RunUpdate()

	a.router = router.New()
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

	a.router.GET("/api/feeds", middleware.AuthMiddleware(a.feedsController.Index))
	a.router.GET("/api/feeds/tsv", middleware.AuthMiddleware(a.feedsController.FeedListTsv))
	a.router.GET("/api/feeds/list-account/tsv", middleware.AuthMiddleware(a.feedsController.ListAccountTsv))
	a.router.GET("/api/feeds/list-network/tsv", middleware.AuthMiddleware(a.feedsController.ListNetworkTsv))

	a.router.GET("/api/supply-side-platforms", middleware.AuthMiddleware(a.sspController.Index))

	a.router.GET("/api/networks", middleware.AuthMiddleware(a.networkController.Index))
	a.router.GET("/api/networks/tsv", middleware.AuthMiddleware(a.networkController.Tsv))

	a.router.GET("/api/discrepancies", middleware.AuthMiddleware(a.discrepController.Index))
}
