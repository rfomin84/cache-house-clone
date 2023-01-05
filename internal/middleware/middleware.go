package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

//Auth middleware
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("api_token")

		if token == os.Getenv("API_TOKEN") {
			return next(c)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	}
}
