package compose

import (
	"bytes"
	"text/template"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/Masterminds/sprig/v3"
)

func GenerateComposeFile(params Parameters) string {
	var out bytes.Buffer

	tmpl, err := template.New("Dockerfile").Funcs(sprig.FuncMap()).Parse(dockerFile)
	if err != nil {
		logging.Sugar.Fatalw("Template rendering failed", "error", err)
	}
	err = tmpl.Execute(&out, params)
	if err != nil {
		logging.Sugar.Fatalw("Template rendering failed", "error", err)
	}

	return out.String()
}
