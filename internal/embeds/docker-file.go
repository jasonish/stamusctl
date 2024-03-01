package embeds

import _ "embed"

//go:embed docker-compose.yaml
var DockerFile string
