package internal

import (
	"github.com/clickadilla/cache-house/internal/controllers"
	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

type App struct {
	router          *router.Router
	feedsController *controllers.FeedsController
}

func (a *App) Run() {
	a.boot()
	a.bootRouting()
	_ = fasthttp.ListenAndServe(":"+os.Getenv("APP_PORT"), a.router.Handler)
}

func (a *App) boot() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	a.router = router.New()

	a.feedsController = &controllers.FeedsController{}
}

func (a *App) bootRouting() {
	a.router.GET("/api/feeds", a.feedsController.Index)
}
