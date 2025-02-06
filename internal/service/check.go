package service

import (
	"errors"

	"github.com/digitranslab/kozmo-sandbox/internal/core/runner/types"
	"github.com/digitranslab/kozmo-sandbox/internal/static"
)

var (
	ErrNetworkDisabled = errors.New("network is disabled, please enable it in the configuration")
)

func checkOptions(options *types.RunnerOptions) error {
	configuration := static.GetKozmoSandboxGlobalConfigurations()

	if options.EnableNetwork && !configuration.EnableNetwork {
		return ErrNetworkDisabled
	}

	return nil
}
