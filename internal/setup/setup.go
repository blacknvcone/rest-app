package setup

import (
	"log/slog"
	"rest-app/config"
)

// BaseURL base url of api
const BaseURL = "/v1/api"

// CloseDB close connection to db
var CloseDB func() error

type SetupData struct {
	ConfigData  config.Config
	InternalApp InternalAppStruct
}

func Init() *SetupData {
	configData := config.GetConfig()

	// LOGGER init
	logger := slog.Default()

	internalAppVar := initInternalApp(logger, configData)

	return &SetupData{
		ConfigData:  configData,
		InternalApp: internalAppVar,
	}
}

func initInternalApp(logger *slog.Logger, conf config.Config) InternalAppStruct {
	var internalAppVar InternalAppStruct

	internalAppVar.Logger = logger
	internalAppVar.Config = conf

	initAppRepo(&internalAppVar)
	initAppService(&internalAppVar)
	initAppHandler(&internalAppVar)

	return internalAppVar
}
