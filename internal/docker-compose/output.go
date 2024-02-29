package compose

import (
	"bytes"
	"os"
	"path"
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

func writeConfGeneric(filePath, outputFile, data string) {
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		logging.Sugar.Errorw("cannot create config folder.", "error", err, "path", filePath)
		return
	}

	f, err := os.OpenFile(path.Join(filePath, outputFile), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logging.Sugar.Errorw("cannot create output file", "error", err, "path", filePath, "outputFile", outputFile)
	}

	defer f.Close()

	f.WriteString(data)
}

func WriteConfigFiles(volumePath string) {

	writeConfGeneric(path.Join(volumePath, "nginx"), "nginx.conf", nginxMainConf)
	writeConfGeneric(path.Join(volumePath, "nginx", "conf.d"), "selks6.conf", selksNginxConfig)

	writeConfGeneric(path.Join(volumePath, "logstash", "conf.d"), "logstash.conf", logstashConfig)
	writeConfGeneric(path.Join(volumePath, "logstash", "templates"), "elasticsearch7-template.json", elasticTemplate)

	writeConfGeneric(path.Join(volumePath, "cron-jobs", "daily"), "scirius-update-suri-rules.sh", cronJobsDailyScirius)
	writeConfGeneric(path.Join(volumePath, "cron-jobs", "daily"), "suricata-logrotate.sh", cronJobsDailySuricata)

	writeConfGeneric(path.Join(volumePath, "suricata", "etc"), "new_entrypoint.sh", suricataEtcEntryPoint)
	writeConfGeneric(path.Join(volumePath, "suricata", "etc"), "selks6-addin.yaml", suricataEtcAddin)

	os.MkdirAll(path.Join(volumePath, "cron-jobs", "1min"), os.ModePerm)
	os.MkdirAll(path.Join(volumePath, "cron-jobs", "15min"), os.ModePerm)
	os.MkdirAll(path.Join(volumePath, "cron-jobs", "hourly"), os.ModePerm)
	os.MkdirAll(path.Join(volumePath, "cron-jobs", "daily"), os.ModePerm)
	os.MkdirAll(path.Join(volumePath, "cron-jobs", "weekly"), os.ModePerm)
	os.MkdirAll(path.Join(volumePath, "cron-jobs", "monthly"), os.ModePerm)

	os.MkdirAll(path.Join(volumePath, "suricata", "logs"), os.ModePerm)
	os.MkdirAll(path.Join(volumePath, "suricata", "logrotate"), os.ModePerm)
}
