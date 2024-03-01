package embeds

import _ "embed"

//go:embed cron-jobs/daily/scirius-update-suri-rules.sh
var CronJobsDailyScirius string

//go:embed cron-jobs/daily/suricata-logrotate.sh
var CronJobsDailySuricata string

//go:embed nginx/conf.d/selks6.conf
var SelksNginxConfig string

//go:embed nginx/nginx.conf
var NginxMainConf string

//go:embed logstash/templates/elasticsearch7-template.json
var ElasticTemplate string

//go:embed logstash/conf.d/logstash.conf
var LogstashConfig string

//go:embed suricata/etc/new_entrypoint.sh
var SuricataEtcEntryPoint string

//go:embed suricata/etc/selks6-addin.yaml
var SuricataEtcAddin string
