package main

import (
	"bonus-client-back/config"
	"bonus-client-back/internal/domain"
	"bonus-client-back/internal/handler"
	"bonus-client-back/internal/httpserver"
	"bonus-client-back/internal/repository"
	"bonus-client-back/internal/service"
	"bonus-client-back/pkg/database"
	"bonus-client-back/pkg/logger"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

const (
	fileName = "./config/config.yaml"
)

func main() {
	zapLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	ctx, cancelContext := context.WithCancel(context.Background())

	conf, err := config.NewConfig(fileName)
	if err != nil {
		zapLogger.Error("error init config", zap.Error(err))
		return
	}

	dbConn, err := database.ConnectToDatabase(&conf.DatabaseConfig)
	if err != nil {
		zapLogger.Error("error connect to database", zap.Error(err))
		return
	}
	defer dbConn.Close()

	if err := database.Migrate(dbConn, zapLogger); !errors.Is(err, domain.ErrExistsTable) {
		zapLogger.Error("error migrate to database", zap.Error(err))
	}

	repo := repository.NewRepository(dbConn)
	service := service.NewServices(ctx, conf, zapLogger, repo)
	handler := handler.NewHandler(service, zapLogger, conf)
	server := httpserver.NewServer(handler.InitHandler())

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(http.ErrServerClosed, err) {
			zapLogger.Error("error server start", zap.Error(err))
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	cancelContext()

	/*
		for i := 3; i > 0; i-- {
			time.Sleep(time.Second)
			fmt.Println(i)
		}
	*/

	zapLogger.Info("application closed")
}
