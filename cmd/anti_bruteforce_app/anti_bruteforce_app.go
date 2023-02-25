package main

import (
	"Anti-bruteforce-service/internal/app"
	"Anti-bruteforce-service/internal/config"
	"go.uber.org/zap"
	"log"
)

func main() {
	log.Println("Init Config from env")
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("can't read env config: %v", err)
		return
	}

	log.Println("Init logger")
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugaredLogger := logger.Sugar()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}(logger)

	application := app.NewAntiBruteforceApp(sugaredLogger, cfg)
	application.StartAppApi()

}
