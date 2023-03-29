package main

import "log"

func setGlobalConfigs(config map[string]interface{}) {
	// check if config is nil
	if config == nil {
		return
	}
	// validate config with globalConfigSchema
	for k, v := range config {
		if globalConfigSchema[k] == nil {
			log.Fatal("Error: invalid config key: ", k)
		}
		switch v.(type) {
		case bool:
			if globalConfigSchema[k] != "bool" {
				log.Fatalf("Error: invalid config value type for key it must be a %s: %s\n", globalConfigSchema[k], k)
			}
		case int:
			if globalConfigSchema[k] != "int" {
				log.Fatalf("Error: invalid config value type for key it must be a %s: %s\n", globalConfigSchema[k], k)

			}
		case float64:
			if globalConfigSchema[k] != "float" {
				log.Fatalf("Error: invalid config value type for key it must be a %s: %s\n", globalConfigSchema[k], k)

			}
		case string:
			if globalConfigSchema[k] != "string" {
				log.Fatalf("Error: invalid config value type for key it must be a %s: %s\n", globalConfigSchema[k], k)

			}
		case []interface{}:
			if globalConfigSchema[k] != "array" {
				log.Fatalf("Error: invalid config value type for key it must be a %s: %s\n", globalConfigSchema[k], k)
			}
		case map[string]interface{}:
			if globalConfigSchema[k] != "map" {
				log.Fatalf("Error: invalid config value type for key it must be a %s: %s\n", globalConfigSchema[k], k)

			}
		case int64:
			if globalConfigSchema[k] != "int" {
				log.Fatalf("Error: invalid config value type for key it must be a %s: %s\n", globalConfigSchema[k], k)

			}
		case float32:
			if globalConfigSchema[k] != "float" {
				log.Fatalf("Error: invalid config value type for key it must be a %s: %s\n", globalConfigSchema[k], k)
			}

		default:
			log.Fatal("Error: invalid config value type for key it must be a bool, int, float, string, array, or map: ", k)
		}
		CONFIGS[k] = v
		log.Printf("Loaded config: %s = %v (%T)\n", k, v, v)
	}

}
