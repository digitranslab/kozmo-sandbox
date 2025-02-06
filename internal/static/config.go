package static

import (
	"os"
	"strconv"
	"strings"

	"github.com/digitranslab/kozmo-sandbox/internal/types"
	"github.com/digitranslab/kozmo-sandbox/internal/utils/log"
	"gopkg.in/yaml.v3"
)

var kozmoSandboxGlobalConfigurations types.KozmoSandboxGlobalConfigurations

func InitConfig(path string) error {
	kozmoSandboxGlobalConfigurations = types.KozmoSandboxGlobalConfigurations{}

	// read config file
	configFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer configFile.Close()

	// parse config file
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&kozmoSandboxGlobalConfigurations)
	if err != nil {
		return err
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err == nil {
		kozmoSandboxGlobalConfigurations.App.Debug = debug
	}

	max_workers := os.Getenv("MAX_WORKERS")
	if max_workers != "" {
		kozmoSandboxGlobalConfigurations.MaxWorkers, _ = strconv.Atoi(max_workers)
	}

	max_requests := os.Getenv("MAX_REQUESTS")
	if max_requests != "" {
		kozmoSandboxGlobalConfigurations.MaxRequests, _ = strconv.Atoi(max_requests)
	}

	port := os.Getenv("SANDBOX_PORT")
	if port != "" {
		kozmoSandboxGlobalConfigurations.App.Port, _ = strconv.Atoi(port)
	}

	timeout := os.Getenv("WORKER_TIMEOUT")
	if timeout != "" {
		kozmoSandboxGlobalConfigurations.WorkerTimeout, _ = strconv.Atoi(timeout)
	}

	api_key := os.Getenv("API_KEY")
	if api_key != "" {
		kozmoSandboxGlobalConfigurations.App.Key = api_key
	}

	python_path := os.Getenv("PYTHON_PATH")
	if python_path != "" {
		kozmoSandboxGlobalConfigurations.PythonPath = python_path
	}

	if kozmoSandboxGlobalConfigurations.PythonPath == "" {
		kozmoSandboxGlobalConfigurations.PythonPath = "/usr/local/bin/python3"
	}

	python_lib_path := os.Getenv("PYTHON_LIB_PATH")
	if python_lib_path != "" {
		kozmoSandboxGlobalConfigurations.PythonLibPaths = strings.Split(python_lib_path, ",")
	}

	if len(kozmoSandboxGlobalConfigurations.PythonLibPaths) == 0 {
		kozmoSandboxGlobalConfigurations.PythonLibPaths = DEFAULT_PYTHON_LIB_REQUIREMENTS
	}

	python_pip_mirror_url := os.Getenv("PIP_MIRROR_URL")
	if python_pip_mirror_url != "" {
		kozmoSandboxGlobalConfigurations.PythonPipMirrorURL = python_pip_mirror_url
	}

	python_deps_update_interval := os.Getenv("PYTHON_DEPS_UPDATE_INTERVAL")
	if python_deps_update_interval != "" {
		kozmoSandboxGlobalConfigurations.PythonDepsUpdateInterval = python_deps_update_interval
	}

	// if not set "PythonDepsUpdateInterval", update python dependencies every 30 minutes to keep the sandbox up-to-date
	if kozmoSandboxGlobalConfigurations.PythonDepsUpdateInterval == "" {
		kozmoSandboxGlobalConfigurations.PythonDepsUpdateInterval = "30m"
	}

	nodejs_path := os.Getenv("NODEJS_PATH")
	if nodejs_path != "" {
		kozmoSandboxGlobalConfigurations.NodejsPath = nodejs_path
	}

	if kozmoSandboxGlobalConfigurations.NodejsPath == "" {
		kozmoSandboxGlobalConfigurations.NodejsPath = "/usr/local/bin/node"
	}

	enable_network := os.Getenv("ENABLE_NETWORK")
	if enable_network != "" {
		kozmoSandboxGlobalConfigurations.EnableNetwork, _ = strconv.ParseBool(enable_network)
	}

	enable_preload := os.Getenv("ENABLE_PRELOAD")
	if enable_preload != "" {
		kozmoSandboxGlobalConfigurations.EnablePreload, _ = strconv.ParseBool(enable_preload)
	}

	allowed_syscalls := os.Getenv("ALLOWED_SYSCALLS")
	if allowed_syscalls != "" {
		strs := strings.Split(allowed_syscalls, ",")
		ary := make([]int, len(strs))
		for i := range ary {
			ary[i], err = strconv.Atoi(strs[i])
			if err != nil {
				return err
			}
		}
		kozmoSandboxGlobalConfigurations.AllowedSyscalls = ary
	}

	if kozmoSandboxGlobalConfigurations.EnableNetwork {
		log.Info("network has been enabled")
		socks5_proxy := os.Getenv("SOCKS5_PROXY")
		if socks5_proxy != "" {
			kozmoSandboxGlobalConfigurations.Proxy.Socks5 = socks5_proxy
		}

		if kozmoSandboxGlobalConfigurations.Proxy.Socks5 != "" {
			log.Info("using socks5 proxy: %s", kozmoSandboxGlobalConfigurations.Proxy.Socks5)
		}

		https_proxy := os.Getenv("HTTPS_PROXY")
		if https_proxy != "" {
			kozmoSandboxGlobalConfigurations.Proxy.Https = https_proxy
		}

		if kozmoSandboxGlobalConfigurations.Proxy.Https != "" {
			log.Info("using https proxy: %s", kozmoSandboxGlobalConfigurations.Proxy.Https)
		}

		http_proxy := os.Getenv("HTTP_PROXY")
		if http_proxy != "" {
			kozmoSandboxGlobalConfigurations.Proxy.Http = http_proxy
		}

		if kozmoSandboxGlobalConfigurations.Proxy.Http != "" {
			log.Info("using http proxy: %s", kozmoSandboxGlobalConfigurations.Proxy.Http)
		}
	}
	return nil
}

// avoid global modification, use value copy instead
func GetKozmoSandboxGlobalConfigurations() types.KozmoSandboxGlobalConfigurations {
	return kozmoSandboxGlobalConfigurations
}

type RunnerDependencies struct {
	PythonRequirements string
}

var runnerDependencies RunnerDependencies

func GetRunnerDependencies() RunnerDependencies {
	return runnerDependencies
}

func SetupRunnerDependencies() error {
	file, err := os.ReadFile("dependencies/python-requirements.txt")
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}

	runnerDependencies.PythonRequirements = string(file)

	return nil
}
