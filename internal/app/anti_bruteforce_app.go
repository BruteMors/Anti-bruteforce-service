package app

import (
	"Anti-bruteforce-service/internal/config"
	"Anti-bruteforce-service/internal/controller/clinterface"
	"Anti-bruteforce-service/internal/controller/grpcapi"
	"Anti-bruteforce-service/internal/controller/httpapi"
	"Anti-bruteforce-service/internal/controller/httpapi/handlers"
	"Anti-bruteforce-service/internal/domain/service"
	"Anti-bruteforce-service/internal/store/postgressql/adapters"
	"Anti-bruteforce-service/internal/store/postgressql/client"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type AntiBruteforceApp struct {
	router                  *httpapi.ApiRouter
	grpcBlackListServer     *grpcapi.BlackListServer
	grpcWhiteListServer     *grpcapi.WhiteListServer
	grpcBucketServer        *grpcapi.BucketServer
	grpcAuthorizationServer *grpcapi.AuthorizationServer
	cli                     *clinterface.CommandLineInterface
	clientDb                *client.PostgresSql
	logger                  *zap.SugaredLogger
	cfg                     *config.Config
}

func NewAntiBruteforceApp(logger *zap.SugaredLogger, cfg *config.Config) *AntiBruteforceApp {
	logger.Infoln("Init http router")
	clientDb := client.NewPostgresSql(logger, cfg)
	err := clientDb.Open()
	if err != nil {
		logger.Fatalf("Troubels with connect to database: %v", err)
	}

	blackListStore := adapters.NewBlackListRepository(clientDb)
	blackListService := service.NewBlackList(blackListStore, logger)
	blacklist := handlers.NewBlackList(blackListService, logger)
	grpcBlackListServer := grpcapi.NewBlackListServer(blackListService, logger)

	whitelistStore := adapters.NewWhiteListRepository(clientDb)
	whitelistService := service.NewWhiteList(whitelistStore, logger)
	whitelist := handlers.NewWhiteList(whitelistService, logger)
	grpcWhiteListServer := grpcapi.NewWhiteListServer(whitelistService, logger)

	authorizationService := service.NewAuthorization(blackListService, whitelistService, cfg, logger)
	auth := handlers.NewAuthorization(authorizationService, logger)
	bucket := handlers.NewBucket(authorizationService, logger)
	grpcBucketServer := grpcapi.NewBucketServer(authorizationService, logger)
	grpcAuthorizationServer := grpcapi.NewAuthorizationServer(authorizationService, logger)

	router := httpapi.NewRouter(auth, blacklist, whitelist, bucket, logger)

	cli := clinterface.New(authorizationService, whitelistService, blackListService)
	//router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	//router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)
	return &AntiBruteforceApp{grpcBlackListServer: grpcBlackListServer, grpcWhiteListServer: grpcWhiteListServer, grpcBucketServer: grpcBucketServer, grpcAuthorizationServer: grpcAuthorizationServer, cli: cli, router: router, clientDb: clientDb, logger: logger, cfg: cfg}
}

func (a *AntiBruteforceApp) StartAppApi() {
	c := make(chan os.Signal, 1)
	go a.cli.Run(c)
	switch a.cfg.Server.ServerType {
	case "grpc":
		a.logger.Infoln("Init grpc server")
		grpcServer := grpcapi.NewServer(a.grpcBlackListServer, a.grpcWhiteListServer, a.grpcBucketServer, a.grpcAuthorizationServer, a.cfg, a.logger)
		go grpcServer.Shutdown(c)
		err := grpcServer.Start()
		if err != nil {
			a.logger.Fatal(err)
		}
		err = a.clientDb.Close()
		if err != nil {
			fmt.Println(err)
		}

	case "http":
		a.logger.Infoln("Register routes")
		a.router.RegisterRoutes()
		a.logger.Infoln("Init http server")

		server := httpapi.NewHttpApiServer(a.router.GetRouter(), a.cfg, a.logger)
		go server.ShutdownService(c)
		err := server.Start()
		if err != nil {
			if err == http.ErrServerClosed {
				a.logger.Infoln(err)
				err = a.clientDb.Close()
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			a.logger.Fatal(err)
		}
	}
}
