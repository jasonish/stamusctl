package embeds

import (
	_ "embed"
)

//go:embed selks/docker-compose.yaml
var DockerFile string
