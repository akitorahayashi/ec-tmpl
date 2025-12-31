package api

import (
	"net/http"

	"example.com/ec-tmpl/internal/protocols"
	"github.com/labstack/echo/v4"
)

type Dependencies struct {
	GreetingService protocols.GreetingService
}

type healthResponse struct {
	Status string `json:"status"`
}

type helloResponse struct {
	Message string `json:"message"`
}

func RegisterRoutes(e *echo.Echo, deps Dependencies) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, healthResponse{Status: "ok"})
	})

	e.GET("/hello/:name", func(c echo.Context) error {
		name := c.Param("name")

		message, err := deps.GreetingService.Greet(c.Request().Context(), name)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "greeting failed")
		}

		return c.JSON(http.StatusOK, helloResponse{Message: message})
	})
}
