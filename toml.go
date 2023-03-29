package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func readTomlConfigFile(configFile string) TOML {

	var configs TOMLRead
	_, err := toml.DecodeFile(configFile, &configs)
	if err != nil {
		log.Fatal(err)
	}

	version := configs.Version
	checkVersion(version) // check if version is greater than or equal to 1.0.0 (the current version)
	var output TOML = TOML{Version: version}
	if configs.Services != nil {
		for k, v := range configs.Services {
			output.Services = append(output.Services, Service{Name: k, Path: v.Path, Command: v.Command, Env: v.Env})

		}

	}
	if configs.Containers != nil {
		for k, v := range configs.Containers {
			output.Docker = append(output.Docker, Docker{Name: k, Image: v.Image, Command: v.Command, Env: v.Env, Compose: v.Compose})

		}

	}

	if configs.Environment != nil {

		for k, v := range configs.Environment {
			// convert v to string and append to globalEnv (but it can be bool, int, float, string, array, or map)
			switch v.(type) {
			case bool:
				globalEnv = append(configs.GlobalEnv, k+"="+strconv.FormatBool(v.(bool)))
			case int:
				configs.GlobalEnv = append(configs.GlobalEnv, k+"="+strconv.Itoa(v.(int)))

			case float64:
				configs.GlobalEnv = append(configs.GlobalEnv, k+"="+strconv.FormatFloat(v.(float64), 'f', -1, 64))
			case string:
				configs.GlobalEnv = append(configs.GlobalEnv, k+"="+v.(string))
			case []interface{}:
				configs.GlobalEnv = append(configs.GlobalEnv, k+"="+strings.Join(v.([]string), ","))
			case map[string]interface{}:
				panic("Error: environment variables cannot be a map")
			}

		}
	}

	return output

}

func writeTomlConfigFile(config TOML, configFile string) {
	// write TOML config file
	var configs TOMLRead = TOMLRead{Version: config.Version}
	for _, v := range config.Services {
		configs.Services[v.Name] = Service{Path: v.Path, Command: v.Command, Env: v.Env}

	}
	for _, v := range config.Docker {
		configs.Containers[v.Name] = Docker{Image: v.Image, Command: v.Command, Env: v.Env, Compose: v.Compose}

	}
	configs.GlobalEnv = globalEnv
	var buf bytes.Buffer
	err := toml.NewEncoder(&buf).Encode(configs)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(configFile, buf.Bytes(), 0644)
}
