package compose

type Parameters struct {
	InterfacesList string

	Registry    string
	DebugMode   bool
	RestartMode string

	ElkVersion          string
	ArkimeviewerVersion string
	SciriusVersion      string

	SciriusToken string

	VolumeDataPath string
	ElasticPath    string

	ElasticMemory  string
	LogstashMemory string

	MLEnabled bool

	NginxExec string
}
