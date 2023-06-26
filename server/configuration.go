package server

import _ "embed"

//go:embed default-config.yaml
var DefaultConfiguration []byte

type Configuration struct {
	Env          string // environment mode ("dev", "prod", etc)
	Port         int    // server port
	InitialToDos []Todo // an initial set of todos
}
