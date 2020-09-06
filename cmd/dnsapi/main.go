//go:generate pkger
package main

import (
	"github.com/dgoujard/rfc2136_web_api/util/rfc2136"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgoujard/rfc2136_web_api/app/app"
	"github.com/dgoujard/rfc2136_web_api/app/router"
	"github.com/dgoujard/rfc2136_web_api/config"
	lr "github.com/dgoujard/rfc2136_web_api/util/logger"
	vr "github.com/dgoujard/rfc2136_web_api/util/validator"
)

func main() {
	appConf := config.AppConfig()

	logger := lr.New(appConf.Debug)

	rfc, err := rfc2136.New(
		appConf.Dns.Host,
		appConf.Dns.Port,
		false,
		appConf.Dns.TsigSecret,
		appConf.Dns.TsigSecretAlg,
		appConf.Dns.TsigKeyname,
		)
	if err != nil {
		logger.Error().Msg(err.Error())
	}

	validator := vr.New()

	application := app.New(logger, validator, rfc, strings.Split(appConf.Dns.Zone,","), appConf.Server.Login, appConf.Server.Password)

	appRouter := router.New(application)

	address := fmt.Sprintf(":%d", appConf.Server.Port)

	logger.Info().Msgf("Starting server %v", address)

	s := &http.Server{
		Addr:         address,
		Handler:      appRouter,
		ReadTimeout:  appConf.Server.TimeoutRead,
		WriteTimeout: appConf.Server.TimeoutWrite,
		IdleTimeout:  appConf.Server.TimeoutIdle,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server startup failed")
	}
}
