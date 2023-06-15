package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/torvald2/hack-fs-2023-promise-card/router"
	"go.uber.org/zap"
)

func main() {
	InitLogger("DEBUG")

	ctx, cancel := context.WithCancel(context.Background())

	conf := GetConfig()
	r := router.NewRouter(conf.PolybaseKey, conf.PolybaseUrl, conf.PolybaseCollection, conf.TimelockHost, conf.TimelockHash)
	srv := &http.Server{
		Addr:    conf.TCPPort,
		Handler: r,
	}

	go func() {
		Logger.Info(fmt.Sprintf("Listen started on port %s", conf.TCPPort))
		if err := srv.ListenAndServe(); err != nil {
			Logger.Panic("Handle server error", zap.Error(err))
		}
	}()

	// Listen for os sygnals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	Logger.Info("App Interrputtes. Waiting for graseful shutdown")
	srv.Shutdown(ctx)
	Logger.Info("Http server stopped")
	cancel()
	os.Exit(0)

}
