package compose

import (
	"bytes"
	"fmt"
	"html/template"

	"stamus-ctl/internal/logging"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/viper"
)

type Parameters struct {
	InterfacesList      string `mapstructure:"suricata.interfaces, squash"`
	SciriusToken        string `mapstructure:"scirius.token, squash"`
	DebugMode           bool   `mapstructure:"scirius.debugMode, squash"`
	SciriusVersion      string `mapstructure:"scirius.version, squash"`
	ArkimeviewerVersion string `mapstructure:"arkimeviewer.version, squash"`
	ElkVersion          string `mapstructure:"elk.version, squash"`
	ElasticPath         string `mapstructure:"elk.elastic.path, squash"`
	ElasticMemory       string `mapstructure:"elk.elastic.memory, squash"`
	MLEnabled           bool   `mapstructure:"elk.elastic.ml, squash"`
	LogstashMemory      string `mapstructure:"elk.logstash.memory, squash"`
	VolumeDataPath      string `mapstructure:"global.volumes.path, squash"`
	RestartMode         string `mapstructure:"global.restartMode, squash"`
	Registry            string `mapstructure:"global.registry, squash"`
	OutputFile          string `mapstructure:"config.outputFile, squash"`
	NginxExec           string `mapstructure:"nginx.exec, squash"`
}

type LoggerFunc func(string) string

func NewParameters() *Parameters {
	return &Parameters{}
}

func NewParametersFromEnv(v *viper.Viper) *Parameters {
	return &Parameters{
		InterfacesList:      v.GetString("suricata.interfaces"),
		SciriusToken:        v.GetString("scirius.token"),
		DebugMode:           v.GetBool("scirius.debugMode"),
		SciriusVersion:      v.GetString("scirius.version"),
		ArkimeviewerVersion: v.GetString("arkimeviewer.version"),
		ElkVersion:          v.GetString("elk.version"),
		ElasticPath:         v.GetString("elk.elastic.path"),
		ElasticMemory:       v.GetString("elk.elastic.memory"),
		MLEnabled:           v.GetBool("elk.elastic.ml"),
		LogstashMemory:      v.GetString("elk.logstash.memory"),
		VolumeDataPath:      v.GetString("global.volumes.path"),
		RestartMode:         v.GetString("global.restartMode"),
		Registry:            v.GetString("global.registry"),
		OutputFile:          v.GetString("config.outputFile"),
		NginxExec:           v.GetString("nginx.exec"),
	}
}

func (params *Parameters) Logs(loggerFunc LoggerFunc) string {
	return loggerFunc(fmt.Sprintf("suricata.interfaces: %s\n", params.InterfacesList)) +
		loggerFunc(fmt.Sprintf("scirius.debugMode: %t\n", params.DebugMode)) +
		loggerFunc(fmt.Sprintf("scirius.version: %s\n", params.SciriusVersion)) +
		loggerFunc(fmt.Sprintf("arkimeviewer.version: %s\n", params.ArkimeviewerVersion)) +
		loggerFunc(fmt.Sprintf("elk.version: %s\n", params.ElkVersion)) +
		loggerFunc(fmt.Sprintf("elk.elastic.path: %s\n", params.ElasticPath)) +
		loggerFunc(fmt.Sprintf("elk.elastic.memory: %s\n", params.ElasticMemory)) +
		loggerFunc(fmt.Sprintf("elk.elastic.ml: %t\n", params.MLEnabled)) +
		loggerFunc(fmt.Sprintf("elk.logstash.memory: %s\n", params.LogstashMemory)) +
		loggerFunc(fmt.Sprintf("global.volumes.path: %s\n", params.VolumeDataPath)) +
		loggerFunc(fmt.Sprintf("global.restartMode: %s\n", params.RestartMode)) +
		loggerFunc(fmt.Sprintf("global.registry: %s\n", params.Registry)) +
		loggerFunc(fmt.Sprintf("config.outputFile: %s\n", params.OutputFile)) +
		loggerFunc(fmt.Sprintf("nginx.exec: %s\n", params.NginxExec))
}

func (params *Parameters) Format(format string) string {
	var out bytes.Buffer

	tmpl, err := template.New("Parameters").Funcs(sprig.FuncMap()).Parse(format)
	if err != nil {
		logging.Sugar.Fatalw("Template rendering failed", "error", err)
	}
	err = tmpl.Execute(&out, params)
	if err != nil {
		logging.Sugar.Fatalw("Template rendering failed", "error", err)
	}

	return out.String()
}
