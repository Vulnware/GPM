package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func readTomlConfigFile(configFile *string) []project {

	configs := make(map[string]interface{})

	_, err := toml.DecodeFile(*configFile, &configs)
	if err != nil {
		log.Fatal(err)
	}
	// check config schema
	if configs["config"] != nil {
		setGlobalConfigs(configs["config"].(map[string]interface{}))
	}
	version := configs["version"].(string)
	checkVersion(version) // check if version is greater than or equal to 1.0.0 (the current version)
	// this code checks if the config has an environment section, then loops through the variables in the environment section
	// and appends them to the globalEnv array. The config section is parsed as a map[string]interface{}, so the values of
	// the map are of type interface{}. The type of the value depends on what is in the yaml config file, but it can be
	// a bool, int, float, string, array, or map. So, the type of the value is checked and then converted to a string and
	// appended to the globalEnv array.

	if configs["enviroment"] != nil {

		for k, v := range configs["environment"].(map[string]interface{}) {
			// convert v to string and append to globalEnv (but it can be bool, int, float, string, array, or map)
			switch v.(type) {
			case bool:
				globalEnv = append(globalEnv, k+"="+strconv.FormatBool(v.(bool)))
			case int:
				globalEnv = append(globalEnv, k+"="+strconv.Itoa(v.(int)))

			case float64:
				globalEnv = append(globalEnv, k+"="+strconv.FormatFloat(v.(float64), 'f', -1, 64))
			case string:
				globalEnv = append(globalEnv, k+"="+v.(string))
			case []interface{}:
				globalEnv = append(globalEnv, k+"="+strings.Join(v.([]string), ","))
			case map[string]interface{}:
				panic("Error: environment variables cannot be a map")
			}

		}
	}
	if configs["services"] != nil {

		for k, v := range configs["services"].(map[string]interface{}) {
			p := v.(map[string]interface{})
			// log p to see what it looks like
			if p["env"] == nil {
				p["env"] = make(map[string]interface{})
			}

			projects = append(projects, project{
				Name:    k,
				Path:    p["path"].(string),
				Command: p["command"].(string),
				Env:     p["env"].(map[string]interface{}),
			})
		}
	} else {
		log.Fatal("Error: no services found in config file")
	}
	if configs["containers"] != nil {
		for k, v := range configs["containers"].(map[string]interface{}) {
			p := v.(map[string]interface{})
			// log p to see what it looks like
			if p["env"] == nil {
				p["env"] = make(map[string]interface{})
			}
			var (
				Image   string
				Command string
				Env     map[string]interface{}
				Compose bool
			)

			if p["image"] != nil {
				Image = p["image"].(string)
			}
			if p["command"] != nil {
				Command = p["command"].(string)
			}
			if p["env"] != nil {
				Env = p["env"].(map[string]interface{})
			}
			if p["compose"] != nil {
				Compose = p["compose"].(bool)
			}

			containers = append(containers, Docker{
				Name:    k,
				Image:   Image,
				Command: Command,
				Env:     Env,
				Compose: Compose,
			})
		}
	}

	return projects

}
