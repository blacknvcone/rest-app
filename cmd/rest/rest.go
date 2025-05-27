package rest

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"rest-app/pkg/constants"
	"rest-app/pkg/validations"

	"rest-app/cmd/rest/middleware"
	"rest-app/config"

	"rest-app/internal/setup"

	ocrServer "rest-app/internal/app/ocr/server"
)

func StartServer(setupData *setup.SetupData) *http.Server {
	conf := config.GetConfig()
	if conf.App.Env == constants.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	// GIN Init
	router := gin.Default()
	router.UseRawPath = true
	validations.InitStructValidation()

	// router.GET("/health", setupData.InternalApp.Handler.HealthCheckHandler.HealthCheck)

	router.Use(middleware.CORSMiddleware())

	initPublicRoute(router, setupData.InternalApp)

	router.Use(middleware.JWTAuthMiddleware())

	initRoute(router, setupData.InternalApp)

	port := config.GetConfig().Http.Port
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}

	go func() {
		// service connections
		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("listen:", err)
		}
	}()
	log.Println("Webserver started")
	return httpServer
}

func initRoute(router *gin.Engine, internalAppStruct setup.InternalAppStruct) {
	//apiRouter := router.Group("/v1/api")
	//ocrServer.Routes.New(apiRouter.Group("/ocr"), internalAppStruct.Handler.OCRHandler)

}

func initPublicRoute(router *gin.Engine, internalAppStruct setup.InternalAppStruct) {
	apiRouter := router.Group("/v1/public-api")

	ocrServer.Routes.New(apiRouter.Group("/ocr"), internalAppStruct.Handler.OCRHandler)

}
