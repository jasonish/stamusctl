package embeds

import (
	_ "embed"
)

//go:embed config/docker-compose.yaml
var DockerFile string
