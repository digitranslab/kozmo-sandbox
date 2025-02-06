package integrationtests_test

import (
	"github.com/digitranslab/kozmo-sandbox/internal/core/runner/python"
	"github.com/digitranslab/kozmo-sandbox/internal/static"
	"github.com/digitranslab/kozmo-sandbox/internal/utils/log"
)

func init() {
	static.InitConfig("conf/config.yaml")

	err := python.PreparePythonDependenciesEnv()
	if err != nil {
		log.Panic("failed to initialize python dependencies sandbox: %v", err)
	}
}
