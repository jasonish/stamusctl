package compose

import (
	"bytes"
	"text/template"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/Masterminds/sprig/v3"
)

const dockerFile = `
# Copyright(C) 2021, Stamus Networks
# Written by Raphaël Brogat <rbrogat@stamus-networks.com>
#
# This file comes with ABSOLUTELY NO WARRANTY!
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.


version: '3.4'

networks:
  network:

volumes:
  elastic-data:  #for ES data persistency
  suricata-rules: #for suricata rules transfer between scirius and suricata and for persistency
  scirius-data: #for scirius data persistency
  scirius-static: #statics files to be served by nginx
  suricata-run: #path where the suricata socket resides
  suricata-logs:
  suricata-logrotate:
    driver_opts:
      type: none
      o: bind
      device: {{.VolumeDataPath | default "./containers-data" }}/suricata/logrotate
  logstash-sincedb: #where logstash stores it's state so it doesn't re-ingest
  arkime-logs:
  arkime-pcap:
  arkime-config:


services:

  elasticsearch:
    container_name: elasticsearch
    image: elastic/elasticsearch:{{.ElkVersion | default "7.16.1"}}
    restart: {{.RestartMode | default "unless-stopped"}}
    healthcheck:
      test: ["CMD-SHELL", "curl --silent --fail localhost:9200/_cluster/health || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 30s
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - xpack.ml.enabled=${ML_ENABLED:-true}
      - ingest.geoip.downloader.enabled=false
    volumes:
      - {{.ElasticPath | default "elastic-data"}}:/usr/share/elasticsearch/data
    mem_limit: ${ELASTIC_MEMORY:-3G}
    ulimits:
      memlock:
        soft: -1
        hard: -1
    networks:
      network:

  kibana:
    container_name: kibana
    image:  elastic/kibana:{{.ElkVersion | default "7.16.1"}}
    restart: {{.RestartMode | default "unless-stopped"}}
    healthcheck:
      test: ["CMD-SHELL", "curl --silent --fail localhost:5601 || exit 1"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 30s
    networks:
      network:

  logstash:
    container_name: logstash
    image:  elastic/logstash:{{.ElkVersion | default "7.16.1"}}
    depends_on:
      scirius:
        condition: service_healthy #because we need to wait for scirius to populate ILM policy
    restart: {{.RestartMode | default "unless-stopped"}}
    healthcheck:
      test: ["CMD-SHELL", "curl --silent --fail localhost:9600 || exit 1"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 30s
    mem_limit: ${LOGSTASH_MEMORY:-2G}
    volumes:
      - logstash-sincedb:/since.db
      - {{.VolumeDataPath | default "./containers-data" }}/suricata/logs:/var/log/suricata:ro
      - {{.VolumeDataPath | default "./containers-data" }}/logstash/conf.d/logstash.conf:/usr/share/logstash/pipeline/logstash.conf
      - {{.VolumeDataPath | default "./containers-data" }}/logstash/templates/elasticsearch7-template.json:/usr/share/logstash/config/elasticsearch7-template.json
    networks:
      network:

  suricata:
    container_name: suricata
    image: jasonish/suricata:master-amd64
    entrypoint: /etc/suricata/new_entrypoint.sh
    restart: {{.RestartMode | default "unless-stopped"}}
    depends_on:
      scirius:
        condition: service_healthy
    environment:
      - SURICATA_OPTIONS={{.InterfacesList}} -vvv --set sensor-name=suricata
    cap_add:
      - NET_ADMIN
      - SYS_NICE
    network_mode: host
    volumes:
       - {{.VolumeDataPath | default "./containers-data" }}/suricata/logs:/var/log/suricata
       - suricata-rules:/etc/suricata/rules
       - suricata-run:/var/run/suricata/
       - {{.VolumeDataPath | default "./containers-data" }}/suricata/etc:/etc/suricata
       - suricata-logrotate:/etc/logrotate.d/

  scirius:
    container_name: scirius
    image: {{.Registry | default "ghcr.io/stamusnetworks"}}/scirius:${SCIRIUS_VERSION:-selks}
    restart: {{.RestartMode | default "unless-stopped"}}
    environment:
      - SECRET_KEY={{.SciriusToken}}
      - DEBUG=${SCIRIUS_DEBUG:-False}
      - SCIRIUS_IN_SELKS=True
      - USE_ELASTICSEARCH=True
      - ELASTICSEARCH_ADDRESS=elasticsearch:9200 #Default
      - USE_KIBANA=True
      - KIBANA_URL=http://kibana:5601 #Default
      - KIBANA_PROXY=True #Get kibana proxied by Scirius
      - ALLOWED_HOSTS=* #allow connexions from anywhere
      - KIBANA7_DASHBOARDS_PATH=/opt/selks/kibana7-dashboards #where to find kibana dashboards
      - SURICATA_UNIX_SOCKET=/var/run/suricata/suricata-command.socket #socket to control suricata
      - USE_EVEBOX=True #gives access to evebox in the top menu
      - EVEBOX_ADDRESS=evebox:5636 #Default
      - USE_SURICATA_STATS=True #display more informations on the suricata page
      - USE_MOLOCH=True
      - MOLOCH_URL=http://arkime:8005

    volumes:
      - scirius-static:/static/
      - scirius-data:/data/
      - {{.VolumeDataPath | default "./containers-data" }}/scirius/logs/:/logs/
      - suricata-rules:/rules
      - suricata-run:/var/run/suricata
      - {{.VolumeDataPath | default "./containers-data" }}/suricata/logs:/var/log/suricata:ro

    networks:
      network:

  evebox:
    container_name: evebox
    image: jasonish/evebox:master
    command: ["-e", "http://elasticsearch:9200"]
    restart: {{.RestartMode | default "unless-stopped"}}
    environment:
      - EVEBOX_HTTP_TLS_ENABLED=false
      - EVEBOX_AUTHENTICATION_REQUIRED=false
    networks:
      network:

  nginx:
    container_name: nginx
    image: nginx
    command: ['${NGINX_EXEC:-nginx}', '-g', 'daemon off;']
    restart: {{.RestartMode | default "unless-stopped"}}
    volumes:
      - scirius-static:/static/:ro
      - {{.VolumeDataPath | default "./containers-data" }}/nginx/conf.d/:/etc/nginx/conf.d/:ro
      - {{.VolumeDataPath | default "./containers-data" }}/nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - {{.VolumeDataPath | default "./containers-data" }}/nginx/ssl:/etc/nginx/ssl:ro
    ports:
      - 443:443
    networks:
      network:

  cron:
    # This containers handles crontabs for the other containers, following the 1 task per container principle.
    # It is based on  'docker:latest' image, wich is an alpine image with docker binary
    container_name: cron
    image: docker:latest
    command: [sh, -c, "echo '*	*	 *	*	 *	run-parts /etc/periodic/1min' >> /etc/crontabs/root && crond -f -l 8"]
    restart: {{.RestartMode | default "unless-stopped"}}
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock # This bind-mout allows using the hosts docker deamon instead of created one inside the container

      # Those volumes will contain the cron jobs
      - {{.VolumeDataPath | default "./containers-data" }}/cron-jobs/1min:/etc/periodic/1min/:ro
      - {{.VolumeDataPath | default "./containers-data" }}/cron-jobs/15min:/etc/periodic/15min/:ro
      - {{.VolumeDataPath | default "./containers-data" }}/cron-jobs/daily:/etc/periodic/daily/:ro
      - {{.VolumeDataPath | default "./containers-data" }}/cron-jobs/hourly:/etc/periodic/hourly/:ro
      - {{.VolumeDataPath | default "./containers-data" }}/cron-jobs/monthly:/etc/periodic/monthly/:ro
      - {{.VolumeDataPath | default "./containers-data" }}/cron-jobs/weekly:/etc/periodic/weekly/:ro


  arkime:
    container_name: arkime
    image: {{.Registry | default "ghcr.io/stamusnetworks"}}/arkimeviewer:${ARKIMEVIEWER_VERSION:-master} ## Repo will need to be changed to stamusnetwork once image built
    restart: {{.RestartMode | default "no"}}
    volumes:
      - {{.VolumeDataPath | default "./containers-data" }}/suricata/logs:/suricata-logs:ro
      - arkime-config:/data/config
      - arkime-pcap:/data/pcap
      - arkime-logs:/data/logs
    networks:
      network:
`

type Parameters struct {
	InterfacesList string
	RestartMode    string
	ElasticPath    string
	VolumeDataPath string
	Registry       string
	SciriusToken   string
	ElkVersion     string
}

func GenerateComposeFile(netInterface, restart, elasticPath, dataPath, registry, token, elkVersion string) string {
	var out bytes.Buffer
	params := Parameters{
		InterfacesList: netInterface,
		RestartMode:    restart,
		ElasticPath:    elasticPath,
		VolumeDataPath: dataPath,
		Registry:       registry,
		SciriusToken:   token,
		ElkVersion:     elkVersion,
	}

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
