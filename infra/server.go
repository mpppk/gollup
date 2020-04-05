package infra

import (
	"fmt"

	"github.com/mpppk/cli-template/handler"

	"github.com/comail/colog"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func registerHandlers(e *echo.Echo, handlers *handler.Handlers) {
	e.GET("/api/sum", handlers.Sum)
	e.GET("/api/sum-history", handlers.SumHistory)
}

func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
}

// NewServer create new echo server with handlers
func NewServer(handlers *handler.Handlers) *echo.Echo {
	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: fmt.Sprintln("method=${method}, uri=${uri}, status=${status}"),
	}))
	e.Use(middleware.BodyDump(bodyDumpHandler))
	registerHandlers(e, handlers)
	return e
}

// InitializeLog initialize log settings
func InitializeLog(verbose bool) {
	colog.Register()
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LInfo)

	if verbose {
		colog.SetMinLevel(colog.LDebug)
	}
}
