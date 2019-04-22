package cmd

import (
	"fmt"
	"net/http"

	"github.com/fhofherr/golf-zap/golfzap"
	"github.com/fhofherr/golf/log"
	"github.com/fhofherr/hylc/pkg/web"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start hylc's http server and serve requests.",
	Run: func(cmd *cobra.Command, args []string) {
		zaplogger, err := zap.NewProduction()
		if err != nil {
			fmt.Printf("%v", errors.Wrap(err, "zap logger"))
		}
		logger := golfzap.New(zaplogger)
		publicRouter := web.NewPublicRouter(web.PublicRouterConfig{
			Logger: logger,
		})
		publicServer := http.Server{
			Addr:    ":8080",
			Handler: publicRouter,
		}
		if err = publicServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Log(logger,
				"level", "error",
				"message", fmt.Sprintf("%+v", errors.Wrap(err, "stop public server")))
		}
	},
}
