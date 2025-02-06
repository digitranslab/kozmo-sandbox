package server

import (
	"fmt"
	"time"

	"github.com/digitranslab/kozmo-sandbox/internal/controller"
	"github.com/digitranslab/kozmo-sandbox/internal/core/runner/python"
	"github.com/digitranslab/kozmo-sandbox/internal/static"
	"github.com/digitranslab/kozmo-sandbox/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func initConfig() {
	// auto migrate database
	err := static.InitConfig("conf/config.yaml")
	if err != nil {
		log.Panic("failed to init config: %v", err)
	}
	log.Info("config init success")

	err = static.SetupRunnerDependencies()
	if err != nil {
		log.Error("failed to setup runner dependencies: %v", err)
	}
	log.Info("runner dependencies init success")
}

func initServer() {
	config := static.GetKozmoSandboxGlobalConfigurations()
	if !config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Logger())
	}

	controller.Setup(r)

	r.Run(fmt.Sprintf(":%d", config.App.Port))
}

func initDependencies() {
	log.Info("installing python dependencies...")
	dependencies := static.GetRunnerDependencies()
	err := python.InstallDependencies(dependencies.PythonRequirements)
	if err != nil {
		log.Panic("failed to install python dependencies: %v", err)
	}
	log.Info("python dependencies installed")

	log.Info("initializing python dependencies sandbox...")
	err = python.PreparePythonDependenciesEnv()
	if err != nil {
		log.Panic("failed to initialize python dependencies sandbox: %v", err)
	}
	log.Info("python dependencies sandbox initialized")

	// start a ticker to update python dependencies to keep the sandbox up-to-date
	go func() {
		updateInterval := static.GetKozmoSandboxGlobalConfigurations().PythonDepsUpdateInterval
		tickerDuration, err := time.ParseDuration(updateInterval)
		if err != nil {
			log.Error("failed to parse python dependencies update interval, skip periodic updates: %v", err)
			return
		}
		ticker := time.NewTicker(tickerDuration)
		for range ticker.C {
			if err := updatePythonDependencies(dependencies); err != nil {
				log.Error("Failed to update Python dependencies: %v", err)
			}
		}
	}()
}

func updatePythonDependencies(dependencies static.RunnerDependencies) error {
	log.Info("Updating Python dependencies...")
	if err := python.InstallDependencies(dependencies.PythonRequirements); err != nil {
		log.Error("Failed to install Python dependencies: %v", err)
		return err
	}
	if err := python.PreparePythonDependenciesEnv(); err != nil {
		log.Error("Failed to prepare Python dependencies environment: %v", err)
		return err
	}
	log.Info("Python dependencies updated successfully.")
	return nil
}

func Run() {
	// init config
	initConfig()
	// init dependencies, it will cost some times
	go initDependencies()

	initServer()
}
