package main

type Service struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Command string `json:"command"`
	// Env is a map of environment iables to set for the command. which can be nil.
	Env map[string]interface{} `json:"env"` // optional environment variables to set
}

type Docker struct {
	// Name is the name of the container
	Name string `json:"name"`
	// Image is the name of the image to use for the container
	Image string `json:"image"`
	// Command is the command to run in the container
	Command string `json:"command"`
	// Env is a map of environment iables to set for the command. which can be nil.
	Env     map[string]interface{} `json:"env"` // optional environment variables to set
	Compose bool                   `json:"compose"`
}

type TOMLRead struct {
	Version     string                 `toml:"version"`
	Services    map[string]Service     `toml:"services"`
	Containers  map[string]Docker      `toml:"containers"`
	Environment map[string]interface{} `toml:"environment"`
	GlobalEnv   []string               // this will be set by the program and not by the config file
} // this struct is used to read the config file

type TOML struct {
	Version     string                 `toml:"version"`
	Services    []Service              `toml:"services"`
	Docker      []Docker               `toml:"docker"`
	Environment map[string]interface{} `toml:"environment"`
	GlobalEnv   []string               // this will be set by the program and not by the config fileş
}
