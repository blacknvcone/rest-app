package setup

import (
	"log/slog"
	"rest-app/config"
	"rest-app/pkg/httpclient"
	"time"

	"github.com/otiai10/gosseract/v2"

	ocrHandler "rest-app/internal/app/ocr/handler"
	ocrPort "rest-app/internal/app/ocr/port"
	ocrRepo "rest-app/internal/app/ocr/repository"
	ocrService "rest-app/internal/app/ocr/service"
)

type InternalAppStruct struct {
	Repositories initRepositoriesApp
	Services     initServicesApp
	Handler      InitHandlerApp
	Config       config.Config
	Logger       *slog.Logger
}

type initRepositoriesApp struct {
	huggingFaceHttpRepo            ocrPort.IHuggingFaceHTTP
	googleaiTextGenerationHTTPRepo ocrPort.IGoogleAIHTTP
}

func initAppRepo(initializeApp *InternalAppStruct) {
	// initializeApp.Repositories.huggingFaceHttpRepo = ocrRepo.NewHuggingFaceHTTP(
	// 	&initializeApp.Config.HuggingFaceAPIConf,
	// 	httpclient.NewRestClient(3*time.Minute,
	// 		initializeApp.Logger))

	initializeApp.Repositories.googleaiTextGenerationHTTPRepo = ocrRepo.NewGoogleAIHTTP(
		&initializeApp.Config.GoogleAIAPIConf,
		httpclient.NewRestClient(3*time.Minute,
			initializeApp.Logger))
}

type initServicesApp struct {
	OCRService ocrPort.IOCRService
}

func initAppService(initializeApp *InternalAppStruct) {
	initializeApp.Services.OCRService = ocrService.NewOCRService(gosseract.NewClient(), initializeApp.Repositories.googleaiTextGenerationHTTPRepo)
}

// HANDLER INIT
type InitHandlerApp struct {
	OCRHandler ocrPort.IOCRHandler
}

func initAppHandler(initializeApp *InternalAppStruct) {
	initializeApp.Handler.OCRHandler = ocrHandler.New(initializeApp.Services.OCRService)
}
