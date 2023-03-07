package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"user-check/api/handlers"
	"user-check/api/middleware"
	"user-check/configuration"
	"user-check/utils/go-stats/concurrency"
	"user-check/utils/logger"
	"time"
)

const httpServerShutdownGracePeriodSeconds = 20

func StartGin(ctx context.Context) {
	defer concurrency.GlobalWaitGroup.Done()

	conf := configuration.AppConfig()
	log := logger.SugaredLogger()


	concurrency.GlobalWaitGroup.Add(1)
	defer func() {
		defer concurrency.GlobalWaitGroup.Done()
		_, localCancel := context.WithTimeout(context.Background(), httpServerShutdownGracePeriodSeconds*time.Second)
		defer localCancel()

	}()
	

	// Set up gin
	log.Debugf("Setting up Gin")
	if !conf.GinLogger {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Set up the middleware
	if conf.GinLogger {
		log.Warnf("Gin's logger is active! Logs will be unstructured!")
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.Use(middleware.CorrelationId())

	// Set up the groups
	userAPI := router.Group("/api/v1")
	{

		// check user exists in ldap
		userAPI.GET("/usercheck/:isid", handlers.UserCheck)
		// count users in ldap
		userAPI.GET("/usercount", handlers.UserGroupCount)
		// health check endpoint
		userAPI.GET("status", handlers.Status)

	}

	// Activate swagger if configured
	if conf.UseSwagger {
		log.Infof("Swagger is active, enabling endpoints")
		url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}




	// Set up the listener
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.HttpPort),
		Handler: router,
	}

	// Start the HTTPS Server
	if conf.Tls {
		go func() {
			log.Infof("API TLS is active, enabling secure communication on port %d", conf.HttpPort)
			log.Debugf("crt file: %s and key file:%s", conf.ApiCertCrtFile, conf.ApiCertKeyFile)
			if err := httpSrv.ListenAndServeTLS(conf.ApiCertCrtFile, conf.ApiCertKeyFile); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalf("Unrecoverable HTTPS Server failure: %s", err.Error())
				}
			}
		}()
	// Start the HTTP Server
	} else {

		go func() {
			log.Infof("Listening on port %d", conf.HttpPort)
			if err := httpSrv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalf("Unrecoverable HTTP Server failure: %s", err.Error())
				}
			}
		}()
	}

	// Block until SIGTERM/SIGINT
	<-ctx.Done()

	// Clean up and shutdown the HTTP server
	cleanCtx, cancel := context.WithTimeout(context.Background(), httpServerShutdownGracePeriodSeconds*time.Second)
	defer cancel()
	log.Infof("Attempting to shutdown the HTTP server with a timeout of %d seconds", httpServerShutdownGracePeriodSeconds)
	if err := httpSrv.Shutdown(cleanCtx); err != nil {
		log.Errorf("HTTP server failed to shutdown gracefully: %s", err.Error())
	} else {
		log.Infof("HTTP Server was shutdown successfully")
	}
}
