package app

import (
	"github.com/dgoujard/rfc2136_web_api/util/rfc2136"
	"github.com/go-playground/validator/v10"

	"github.com/dgoujard/rfc2136_web_api/util/logger"
)

const (
	appErrDataCreationFailure    = "data creation failure"
	appErrDataAccessFailure      = "data access failure"
	appErrDataUpdateFailure      = "data update failure"
	appErrJsonCreationFailure    = "json creation failure"
	appErrFormDecodingFailure    = "form decoding failure"
	appErrFormErrResponseFailure = "form error response failure"
)

type App struct {
	logger    *logger.Logger
	validator *validator.Validate
	rfc2136 *rfc2136.Rfc2136
	zones []string
	username string
	password string
}

func New(
	logger *logger.Logger,
	validator *validator.Validate,
	rfc2136 *rfc2136.Rfc2136,
	zones []string,
	username string,
	password string,
) *App {
	return &App{
		logger:    logger,
		validator: validator,
		rfc2136: rfc2136,
		zones: zones,
		username: username,
		password: password,
	}
}

func (app *App) Logger() *logger.Logger {
	return app.logger
}
