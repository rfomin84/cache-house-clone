package internal

import (
	"github.com/clickadilla/cache-house/internal/controllers"
	"github.com/clickadilla/cache-house/internal/controllers/api/public"
	"github.com/clickadilla/cache-house/internal/managers"
	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
)

type App struct {
	router          *router.Router
	homeController  *controllers.HomeController
	feedsController *public.FeedsController
	logger          *logrus.Logger
}

func (a *App) Run() {
	a.boot()
	a.bootRouting()

	_ = fasthttp.ListenAndServe(":"+os.Getenv("APP_PORT"), a.router.Handler)
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

	a.router = router.New()
	a.feedsController = &public.FeedsController{
		FeedState: feedState,
	}
	a.homeController = &controllers.HomeController{
		FeedState: feedState,
	}
}

func (a *App) bootRouting() {
	a.router.GET("/", a.homeController.Index)
	a.router.GET("/api/feeds", a.feedsController.Index)
}
