package compose

const nginxMainConf = `
user www-data;
worker_processes auto;
pid /run/nginx.pid;
include /etc/nginx/modules-enabled/*.conf;

events {
	worker_connections 768;
	# multi_accept on;
}

http {

	##
	# Basic Settings
	##

	sendfile on;
	tcp_nopush on;
	tcp_nodelay on;
	keepalive_timeout 65;
	types_hash_max_size 2048;
  client_max_body_size 20M;
	# server_tokens off;

	# server_names_hash_bucket_size 64;
	# server_name_in_redirect off;

	include /etc/nginx/mime.types;
	default_type application/octet-stream;

	##
	# SSL Settings
	##

	ssl_protocols TLSv1 TLSv1.1 TLSv1.2; # Dropping SSLv3, ref: POODLE
	ssl_prefer_server_ciphers on;

	##
	# Logging Settings
	##

	access_log /var/log/nginx/access.log;
	error_log /var/log/nginx/error.log;

	##
	# Gzip Settings
	##

	gzip on;

	# gzip_vary on;
	# gzip_proxied any;
	# gzip_comp_level 6;
	# gzip_buffers 16 8k;
	# gzip_http_version 1.1;
	# gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

	##
	# Virtual Host Configs
	##

	include /etc/nginx/conf.d/*.conf;
	#include /etc/nginx/sites-enabled/*;
}


#mail {
#	# See sample authentication script at:
#	# http://wiki.nginx.org/ImapAuthenticateWithApachePhpScript
#
#	# auth_http localhost/auth.php;
#	# pop3_capabilities "TOP" "USER";
#	# imap_capabilities "IMAP4rev1" "UIDPLUS";
#
#	server {
#		listen     localhost:110;
#		protocol   pop3;
#		proxy      on;
#	}
#
#	server {
#		listen     localhost:143;
#		protocol   imap;
#		proxy      on;
#	}
#}
`

const selksNginxConfig = `
server {
    listen 80 default_server;
    listen 443 default_server ssl;
    ssl_certificate /etc/nginx/ssl/scirius.crt;
    ssl_certificate_key /etc/nginx/ssl/scirius.key;
    server_name SELKS;
    access_log /var/log/nginx/scirius.access.log;
    error_log /var/log/nginx/scirius.error.log;

    # https://docs.djangoproject.com/en/dev/howto/static-files/#serving-static-files-in-production
    location /static/ { # STATIC_URL
        alias /static/; # STATIC_ROOT
        expires 30d;
    }

    location /media/ { # MEDIA_URL
        alias /static/; # MEDIA_ROOT
        expires 30d;
    }

    location /plugins/ {
        proxy_pass http://kibana:5601/plugins/;
        proxy_redirect off;
    }

    location /dlls/ {
        proxy_pass http://kibana:5601/dlls/;
        proxy_redirect off;
    }

    location /socket.io/ {
        proxy_pass http://kibana:5601/socket.io/;
        proxy_redirect off;
    }

    location /dataset/ {
        proxy_pass http://kibana:5601/dataset/;
        proxy_redirect off;
    }

    location /translations/ {
        proxy_pass http://kibana:5601/translations/;
        proxy_redirect off;
    }

    location ^~ /built_assets/ {
        proxy_pass http://kibana:5601/built_assets/;
        proxy_redirect off;
    }

    location /ui/ {
        proxy_pass http://kibana:5601/ui/;
        proxy_redirect off;
    }

   location /spaces/ {
        proxy_pass http://kibana:5601/spaces/;
        proxy_redirect off;
    }

  location /node_modules/ {
        proxy_pass http://kibana:5601/node_modules/;
        proxy_redirect off;
    }

  location /bootstrap.js {
        proxy_pass http://kibana:5601/bootstrap.js;
        proxy_redirect off;
    }

 location /internal/ {
        proxy_pass http://kibana:5601/internal/;
        proxy_redirect off;
    }

 location ~ "^/([\d]{5}/.*)" {
        proxy_pass http://kibana:5601/$1;
        proxy_redirect off;
    }


 location / {
       proxy_pass http://scirius:8000;
       proxy_read_timeout 600;
       proxy_set_header Host $http_host;
       proxy_set_header X-Forwarded-Proto https;
       proxy_redirect off;
       client_max_body_size 100M;
    }

}
`

const logstashConfig = `
input {
  file {
    path => ["/var/log/suricata/*.json"]
    #sincedb_path => ["/var/lib/logstash/"]
    sincedb_path => ["/usr/share/logstash/since.db"]
    codec =>   json
    type => "SELKS"
  }

}

filter {
  if [type] == "SELKS" {

    date {
      match => [ "timestamp", "ISO8601" ]
    }

    ruby {
      code => "
        if event.get('[event_type]') == 'fileinfo'
          event.set('[fileinfo][type]', event.get('[fileinfo][magic]').to_s.split(',')[0])
        end
      "
    }
    ruby {
      code => "
        if event.get('[event_type]') == 'alert'
          sp = event.get('[alert][signature]').to_s.split(' group ')
          if (sp.length == 2) and /\A\d+\z/.match(sp[1])
            event.set('[alert][signature]', sp[0])
          end
        end
      "
     }

    metrics {
      meter => [ "eve_insert" ]
      add_tag => "metric"
      flush_interval => 30
    }
  }

  if [http] {
    useragent {
       source => "[http][http_user_agent]"
       target => "[http][user_agent]"
    }
  }
  if [src_ip]  {
    geoip {
      source => "src_ip"
      target => "geoip"
      #database => "/opt/logstash/vendor/geoip/GeoLiteCity.dat"
      #add_field => [ "[geoip][coordinates]", "%{[geoip][longitude]}" ]
      #add_field => [ "[geoip][coordinates]", "%{[geoip][latitude]}"  ]
    }
  }
    if [dest_ip]  {
    geoip {
      source => "dest_ip"
      target => "geoip"
      #database => "/opt/logstash/vendor/geoip/GeoLiteCity.dat"
      #add_field => [ "[geoip][coordinates]", "%{[geoip][longitude]}" ]
      #add_field => [ "[geoip][coordinates]", "%{[geoip][latitude]}"  ]
    }
  }
}

output {
  if [event_type] and [event_type] != 'stats' {
    elasticsearch {
      hosts => "elasticsearch"
      index => "logstash-%{event_type}-%{+YYYY.MM.dd}"
      template_overwrite => true
      template => "/usr/share/logstash/config/elasticsearch7-template.json"
    }
  } else {
    elasticsearch {
      hosts => "elasticsearch"
      index => "logstash-%{+YYYY.MM.dd}"
      template_overwrite => true
      template => "/usr/share/logstash/config/elasticsearch7-template.json"
    }
  }
}
`

const elasticTemplate = `
{
  "template" : "logstash-*",
  "order": 1,
  "version" : 70002,
  "index_patterns": [
    "logstash-*"
  ],
  "settings" : {
        "index": {
      "number_of_replicas": 0,
      "refresh_interval" : "5s",
      "mapping.total_fields.limit": 10000
    },
      "analysis": {
      "analyzer": {
        "sn_analyzer": {
          "type": "custom",
          "tokenizer": "whitespace",
          "char_filter": [
            "sn_lowercase"
          ]
        }
      },
      "char_filter": {
                "sn_lowercase": {
          "type": "mapping",
          "mappings": [
                      "A => a",                      "B => b",                      "C => c",                      "D => d",                      "E => e",                      "F => f",                      "G => g",                      "H => h",                      "I => i",                      "J => j",                      "K => k",                      "L => l",                      "M => m",                      "N => n",                      "O => o",                      "P => p",                      "Q => q",                      "R => r",                      "S => s",                      "T => t",                      "U => u",                      "V => v",                      "W => w",                      "X => x",                      "Y => y",                      "Z => z"                    ]
        }
      }
    }
  },
  "mappings" : {
    "dynamic_templates" : [ {
      "message_field" : {
        "path_match" : "message",
        "match_mapping_type" : "string",
        "mapping" : {
          "type" : "text",
          "norms" : false,
          "analyzer": "sn_analyzer",
          "search_analyzer":"sn_analyzer",
          "search_quote_analyzer":"sn_analyzer"
        }
      }
    }, {
      "string_fields" : {
        "match" : "*",
        "match_mapping_type" : "string",
        "mapping" : {
          "type" : "text", "norms" : false,
          "analyzer": "sn_analyzer",
          "search_analyzer":"sn_analyzer",
          "search_quote_analyzer":"sn_analyzer",
          "fields" : {
            "keyword" : { "type": "keyword", "ignore_above": 256 },
            "raw" : { "type": "keyword", "ignore_above": 256 }
          }
        }
      }
    }, {
      "percentage_fields_long_to_float": {
        "path_match": "*.pct",
        "match_mapping_type": "long",
        "mapping": {
          "type": "float"
        }
      }
    } ],
    "properties" : {
      "@timestamp": { "type": "date" },
      "@version": { "type": "keyword" },
      "geoip"  : {
        "dynamic": true,
        "properties" : {
          "ip": { "type": "ip" },
          "location" : { "type" : "geo_point" },
          "latitude" : { "type" : "half_float" },
          "longitude" : { "type" : "half_float" }
        }
      },
      "discovery"  : {
        "dynamic": true,
        "properties" : {
          "asset": {
            "type": "ip",
            "fields": {
              "raw": {"type": "keyword"},
              "keyword": {"type": "keyword"}
            }
          }
        }
      },
      "dest_ip": {
          "type": "ip",
          "fields": {
              "raw": {"type": "keyword"},
              "keyword": {"type": "keyword"}
           }
      },
      "src_ip": {
          "type": "ip",
          "fields": {
              "raw": {"type": "keyword"},
              "keyword": {"type": "keyword"}
           }
      },
      "cpu": {
        "properties": {
          "system_p": {
            "doc_values": "true",
            "type": "float"
          },
          "user_p": {
            "doc_values": "true",
            "type": "float"
          }
        }
      },
      "fs": {
        "properties": {
          "used_p": {
            "doc_values": "true",
            "type": "float"
          }
        }
      },
      "load": {
        "properties": {
          "load1": {
            "doc_values": "true",
            "type": "float"
          },
          "load15": {
            "doc_values": "true",
            "type": "float"
          },
          "load5": {
            "doc_values": "true",
            "type": "float"
          }
        }
      },
      "mem": {
        "properties": {
          "actual_used_p": {
            "doc_values": "true",
            "type": "float"
          },
          "used_p": {
            "doc_values": "true",
            "type": "float"
          }
        }
      },
      "proc": {
        "properties": {
          "cpu": {
            "properties": {
              "user_p": {
                "doc_values": "true",
                "type": "float"
              }
            }
          },
          "mem": {
            "properties": {
              "rss_p": {
                "doc_values": "true",
                "type": "float"
              }
            }
          }
        }
      },
      "swap": {
        "properties": {
          "used_p": {
            "doc_values": "true",
            "type": "float"
          }
        }
      },
      "ip": {
        "type": "ip"
      },
      "alert": {
        "properties": {
          "source": {
            "properties": {
              "ip": {
                "type": "ip",
                "fields": {
                  "keyword": {"type": "keyword"}
                }
              }
            }
          },
          "target": {
            "properties": {
              "ip": {
                "type": "ip",
                "fields": {
                  "keyword": {"type": "keyword"}
                }
              }
            }
          }
        }
      }
    }
  }
}
`

const cronJobsDailyScirius = `
#!/usr/bin/env bash

echo "Updating Suricata rules from Scirius"
docker exec scirius python /opt/scirius/manage.py updatesuricata && echo "done." || echo "ERROR"
`

const cronJobsDailySuricata = `
#!/usr/bin/env bash
#
# Example of rotating the logs within the Suricata container.
#
# Add -v for verbose output.
# Add -f to force rotation.

echo "Rotating Suricata logs"
docker exec suricata logrotate -v /etc/logrotate.d/suricata $@ && echo "done." || echo "ERROR"
`

const suricataEtcEntryPoint = `
#!/bin/bash
set -e

fix_perms() {
    if [[ "${PGID}" ]]; then
        groupmod -o -g "${PGID}" suricata
    fi

    if [[ "${PUID}" ]]; then
        usermod -o -u "${PUID}" suricata
    fi

    chown -R suricata:suricata /etc/suricata
    chown -R suricata:suricata /var/lib/suricata
    chown -R suricata:suricata /var/log/suricata
    chown -R suricata:suricata /var/run/suricata
}

for src in /etc/suricata.dist/*; do
    filename=$(basename ${src})
    dst="/etc/suricata/${filename}"
    if ! test -e "${dst}"; then
        echo "Creating ${dst}."
        cp -a "${src}" "${dst}"
    fi
done

mkdir -p /var/log/suricata/fpc/
cat /etc/suricata/suricata.yaml | grep "include: selks6-addin.yaml" || echo "include: selks6-addin.yaml" >> /etc/suricata/suricata.yaml && echo 'suricata.yaml edited'

exec /docker-entrypoint.sh $@
`

const suricataEtcAddin = `
%YAML 1.1
---

# Suricata configuration file SELKS addition.
# This file is added to /etc/suricata/suricata.yaml and overrides
# specific settings

# IP Reputation
reputation-categories-file: /etc/suricata/rules/scirius-categories.txt
default-reputation-path: /etc/suricata/rules
reputation-files:
 - scirius-iprep.list

##
## Configure Suricata to load Suricata-Update managed rules.
##
## If this section is completely commented out move down to the "Advanced rule
## file configuration".
##

#default-rule-path: /etc/suricata/rules
#rule-files:
# - suricata.rules

##
## Advanced rule file configuration.
##
## If this section is completely commented out then your configuration
## is setup for suricata-update as it was most likely bundled and
## installed with Suricata.
##

default-rule-path: /etc/suricata/rules
rule-files:
 - scirius.rules
# - botcc.rules
 ## - botcc.portgrouped.rules
# - ciarmy.rules
# - compromised.rules
# - drop.rules
# - dshield.rules
## - emerging-activex.rules
# - emerging-attack_response.rules
# - emerging-chat.rules
# - emerging-current_events.rules
# - emerging-dns.rules
# - emerging-dos.rules
# - emerging-exploit.rules
# - emerging-ftp.rules
## - emerging-games.rules
## - emerging-icmp_info.rules
## - emerging-icmp.rules
# - emerging-imap.rules
## - emerging-inappropriate.rules
## - emerging-info.rules
# - emerging-malware.rules
# - emerging-misc.rules
# - emerging-mobile_malware.rules
# - emerging-netbios.rules
# - emerging-p2p.rules
# - emerging-policy.rules
# - emerging-pop3.rules
# - emerging-rpc.rules
## - emerging-scada.rules
## - emerging-scada_special.rules
# - emerging-scan.rules
## - emerging-shellcode.rules
# - emerging-smtp.rules
# - emerging-snmp.rules
# - emerging-sql.rules
# - emerging-telnet.rules
# - emerging-tftp.rules
# - emerging-trojan.rules
# - emerging-user_agents.rules
# - emerging-voip.rules
# - emerging-web_client.rules
# - emerging-web_server.rules
## - emerging-web_specific_apps.rules
# - emerging-worm.rules
# - tor.rules
## - decoder-events.rules # available in suricata sources under rules dir
## - stream-events.rules  # available in suricata sources under rules dir
# - http-events.rules    # available in suricata sources under rules dir
# - smtp-events.rules    # available in suricata sources under rules dir
# - dns-events.rules     # available in suricata sources under rules dir
# - tls-events.rules     # available in suricata sources under rules dir
## - modbus-events.rules  # available in suricata sources under rules dir
## - app-layer-events.rules  # available in suricata sources under rules dir
## - dnp3-events.rules       # available in suricata sources under rules dir

classification-file: /etc/suricata/rules/classification.config
reference-config-file: /etc/suricata/reference.config
threshold-file: /etc/suricata/rules/threshold.config

# Daemon working directory
# Suricata will change directory to this one if provided
# Default: "/"
daemon-directory: "/var/log/suricata/core"


##
## Performance tuning and profiling
##

# The detection engine builds internal groups of signatures. The engine
# allow us to specify the profile to use for them, to manage memory on an
# efficient way keeping a good performance. For the profile keyword you
# can use the words "low", "medium", "high" or "custom". If you use custom
# make sure to define the values at "- custom-values" as your convenience.
# Usually you would prefer medium/high/low.
#
# "sgh mpm-context", indicates how the staging should allot mpm contexts for
# the signature groups.  "single" indicates the use of a single context for
# all the signature group heads.  "full" indicates a mpm-context for each
# group head.  "auto" lets the engine decide the distribution of contexts
# based on the information the engine gathers on the patterns from each
# group head.
#
# The option inspection-recursion-limit is used to limit the recursive calls
# in the content inspection code.  For certain payload-sig combinations, we
# might end up taking too much time in the content inspection code.
# If the argument specified is 0, the engine uses an internally defined
# default limit.  On not specifying a value, we use no limits on the recursion.
detect:
  profile: high
  custom-values:
    toclient-groups: 3
    toserver-groups: 25
  sgh-mpm-context: auto
  inspection-recursion-limit: 3000
  # If set to yes, the loading of signatures will be made after the capture
  # is started. This will limit the downtime in IPS mode.
  #delayed-detect: yes

  prefilter:
    # default prefiltering setting. "mpm" only creates MPM/fast_pattern
    # engines. "auto" also sets up prefilter engines for other keywords.
    # Use --list-keywords=all to see which keywords support prefiltering.
    default: auto

  # the grouping values above control how many groups are created per
  # direction. Port whitelisting forces that port to get it's own group.
  # Very common ports will benefit, as well as ports with many expensive
  # rules.
  grouping:
    tcp-whitelist: 53, 80, 139, 443, 445, 1433, 3306, 3389, 6666, 6667, 8080
    udp-whitelist: 53, 135, 5060


# Runmode the engine should use. Please check --list-runmodes to get the available
# runmodes for each packet acquisition method. Defaults to "autofp" (auto flow pinned
# load balancing).
#runmode: autofp
runmode: workers

##
## Step 3: select outputs to enable
##

# The default logging directory.  Any log or output file will be
# placed here if its not specified with a full path name. This can be
# overridden with the -l command line parameter.
default-log-dir: /var/log/suricata/

# global stats configuration
stats:
  enabled: yes
  # The interval field (in seconds) controls at what interval
  # the loggers are invoked.
  interval: 8
  # Add decode events as stats.
  decoder-events: true
  # Decoder event prefix in stats. Has been 'decoder' before, but that leads
  # to missing events in the eve.stats records. See issue #2225.
  decoder-events-prefix: "decoder.event"
  # Add stream events as stats.
  stream-events: true

# Configure the type of alert (and other) logging you would like.
outputs:
  # a line based alerts log similar to Snort's fast.log
  - fast:
      enabled: no
      filename: fast.log
      append: yes
      #filetype: regular # 'regular', 'unix_stream' or 'unix_dgram'

  # Extensible Event Format (nicknamed EVE) event log in JSON format
  - eve-log:
      enabled: yes
      filetype: regular #regular|syslog|unix_dgram|unix_stream|redis
      filename: eve.json
      #prefix: "@cee: " # prefix to prepend to each log entry
      # the following are valid when type: syslog above
      #identity: "suricata"
      #facility: local5
      #level: Info ## possible levels: Emergency, Alert, Critical,
                   ## Error, Warning, Notice, Info, Debug
      #redis:
      #  server: 127.0.0.1
      #  port: 6379
      #  async: true ## if redis replies are read asynchronously
      #  mode: list ## possible values: list|lpush (default), rpush, channel|publish
      #             ## lpush and rpush are using a Redis list. "list" is an alias for lpush
      #             ## publish is using a Redis channel. "channel" is an alias for publish
      #  key: suricata ## key or channel to use (default to suricata)
      # Redis pipelining set up. This will enable to only do a query every
      # 'batch-size' events. This should lower the latency induced by network
      # connection at the cost of some memory. There is no flushing implemented
      # so this setting as to be reserved to high traffic suricata.
      #  pipelining:
      #    enabled: yes ## set enable to yes to enable query pipelining
      #    batch-size: 10 ## number of entry to keep in buffer

      # Include top level metadata. Default yes.
      #metadata: no

      pcap-file: false

      # Community Flow ID
      # Adds a 'community_id' field to EVE records. These are meant to give
      # records a predictable flow ID that can be used to match records to
      # output of other tools such as Zeek (Bro).
      #
      # Takes a 'seed' that needs to be same across sensors and tools
      # to make the id less predictable.

      # enable/disable the community id feature.
      community-id: true
      # Seed value for the ID output. Valid values are 0-65535.
      community-id-seed: 1

      # HTTP X-Forwarded-For support by adding an extra field or overwriting
      # the source or destination IP address (depending on flow direction)
      # with the one reported in the X-Forwarded-For HTTP header. This is
      # helpful when reviewing alerts for traffic that is being reverse
      # or forward proxied.
      xff:
        enabled: yes
        # Two operation modes are available, "extra-data" and "overwrite".
        mode: extra-data
        # Two proxy deployments are supported, "reverse" and "forward". In
        # a "reverse" deployment the IP address used is the last one, in a
        # "forward" deployment the first IP address is used.
        deployment: reverse
        # Header name where the actual IP address will be reported, if more
        # than one IP address is present, the last IP address will be the
        # one taken into consideration.
        header: X-Forwarded-For

      types:
        - alert:
            payload: yes             # enable dumping payload in Base64
            # payload-buffer-size: 4kb # max size of payload buffer to output in eve-log
            payload-printable: yes   # enable dumping payload in printable (lossy) format
            packet: yes              # enable dumping of packet (without stream segments)
            http-body: yes           # enable dumping of http body in Base64
            http-body-printable: yes # enable dumping of http body in printable format
            # metadata: no             # enable inclusion of app layer metadata with alert. Default yes

            # Enable the logging of tagged packets for rules using the
            # "tag" keyword.
            tagged-packets: yes
        - anomaly:
            # Anomaly log records describe unexpected conditions such
            # as truncated packets, packets with invalid IP/UDP/TCP
            # length values, and other events that render the packet
            # invalid for further processing or describe unexpected
            # behavior on an established stream. Networks which
            # experience high occurrences of anomalies may experience
            # packet processing degradation.
            #
            # Anomalies are reported for the following:
            # 1. Decode: Values and conditions that are detected while
            # decoding individual packets. This includes invalid or
            # unexpected values for low-level protocol lengths as well
            # as stream related events (TCP 3-way handshake issues,
            # unexpected sequence number, etc).
            # 2. Stream: This includes stream related events (TCP
            # 3-way handshake issues, unexpected sequence number,
            # etc).
            # 3. Application layer: These denote application layer
            # specific conditions that are unexpected, invalid or are
            # unexpected given the application monitoring state.
            #
            # By default, anomaly logging is enabled. When anomaly
            # logging is enabled, applayer anomaly reporting is
            # also enabled.
            enabled: yes
            #
            # Choose one or more types of anomaly logging and whether to enable
            # logging of the packet header for packet anomalies.
            types:
              decode: no
              stream: no
              applayer: yes
            #packethdr: no
        - http:
            extended: yes     # enable this for extended logging information
            # custom allows additional http fields to be included in eve-log
            # the example below adds three additional fields when uncommented
            #custom: [Accept-Encoding, Accept-Language, Authorization]
            custom: [accept, accept-charset, accept-encoding, accept-language,
            accept-datetime, authorization, cache-control, cookie, from,
            max-forwards, origin, pragma, proxy-authorization, range, te, via,
            x-requested-with, dnt, x-forwarded-proto, accept-range, age,
            allow, connection, content-encoding, content-language,
            content-length, content-location, content-md5, content-range,
            content-type, date, etags, last-modified, link, location,
            proxy-authenticate, referrer, refresh, retry-after, server,
            set-cookie, trailer, transfer-encoding, upgrade, vary, warning,
            www-authenticate, true-client-ip, org-src-ip, x-bluecoat-via]
            # set this value to one among {both, request, response} to dump all
            # http headers for every http request and/or response
            dump-all-headers: [both]
        - dns:
            # This configuration uses the new DNS logging format,
            # the old configuration is still available:
            # http://suricata.readthedocs.io/en/latest/configuration/suricata-yaml.html#eve-extensible-event-format
            # Use version 2 logging with the new format:
            # DNS answers will be logged in one single event
            # rather than an event for each of it.
            # Without setting a version the version
            # will fallback to 1 for backwards compatibility.
            version: 2

            # Enable/disable this logger. Default: enabled.
            #enabled: no

            # Control logging of requests and responses:
            # - requests: enable logging of DNS queries
            # - responses: enable logging of DNS answers
            # By default both requests and responses are logged.
            #requests: no
            #responses: no

            # Format of answer logging:
            # - detailed: array item per answer
            # - grouped: answers aggregated by type
            # Default: all
            #formats: [detailed, grouped]

            # Answer types to log.
            # Default: all
            #types: [a, aaaa, cname, mx, ns, ptr, txt]
        - tls:
            extended: yes     # enable this for extended logging information
            # output TLS transaction where the session is resumed using a
            # session id
            #session-resumption: no
            # custom allows to control which tls fields that are included
            # in eve-log
            #custom: [subject, issuer, session_resumed, serial, fingerprint, sni, version, not_before, not_after, certificate, chain, ja3, ja3s]
        - files:
            force-magic: yes   # force logging magic on all logged files
            # force logging of checksums, available hash functions are md5,
            # sha1 and sha256
            force-hash: [md5, sha1, sha256]
        #- drop:
        #    alerts: yes      # log alerts that caused drops
        #    flows: all       # start or all: 'start' logs only a single drop
        #                     # per flow direction. All logs each dropped pkt.
        - smtp:
            #extended: yes # enable this for extended logging information
            # this includes: bcc, message-id, subject, x_mailer, user-agent
            # custom fields logging from the list:
            #  reply-to, bcc, message-id, subject, x-mailer, user-agent, received,
            #  x-originating-ip, in-reply-to, references, importance, priority,
            #  sensitivity, organization, content-md5, date
            #custom: [received, x-mailer, x-originating-ip, relays, reply-to, bcc]
            custom: [received, x-mailer, x-originating-ip, relays, reply-to, bcc,
            reply-to, bcc, message-id, subject, x-mailer, user-agent, received,
            x-originating-ip, in-reply-to, references, importance, priority,
            sensitivity, organization, content-md5, date]
            # output md5 of fields: body, subject
            # for the body you need to set app-layer.protocols.smtp.mime.body-md5
            # to yes
            md5: [body, subject]
        - dnp3
        - ftp
        - rdp
        - nfs
        - smb
        - tftp
        - ike
        - krb5
        - snmp
        - rfb
        - sip
        - ssh
        - dhcp:
            # DHCP logging requires Rust.
            enabled: yes
            # When extended mode is on, all DHCP messages are logged
            # with full detail. When extended mode is off (the
            # default), just enough information to map a MAC address
            # to an IP address is logged.
            extended: yes
        - stats:
            totals: yes       # stats for all threads merged together
            threads: no       # per thread stats
            deltas: yes        # include delta values
        # bi-directional flows
        - flow
        # uni-directional flows
        #- netflow
        # Metadata event type. Triggered whenever a pktvar is saved
        # and will include the pktvars, flowvars, flowbits and
        # flowints.
        - metadata

  # alert output for use with Barnyard2
  - unified2-alert:
      enabled: no
      filename: unified2.alert

      # File size limit.  Can be specified in kb, mb, gb.  Just a number
      # is parsed as bytes.
      #limit: 32mb

      # By default unified2 log files have the file creation time (in
      # unix epoch format) appended to the filename. Set this to yes to
      # disable this behaviour.
      #nostamp: no

      # Sensor ID field of unified2 alerts.
      #sensor-id: 0

      # Include payload of packets related to alerts. Defaults to true, set to
      # false if payload is not required.
      #payload: yes

      # HTTP X-Forwarded-For support by adding the unified2 extra header or
      # overwriting the source or destination IP address (depending on flow
      # direction) with the one reported in the X-Forwarded-For HTTP header.
      # This is helpful when reviewing alerts for traffic that is being reverse
      # or forward proxied.
      xff:
        enabled: no
        # Two operation modes are available, "extra-data" and "overwrite". Note
        # that in the "overwrite" mode, if the reported IP address in the HTTP
        # X-Forwarded-For header is of a different version of the packet
        # received, it will fall-back to "extra-data" mode.
        mode: extra-data
        # Two proxy deployments are supported, "reverse" and "forward". In
        # a "reverse" deployment the IP address used is the last one, in a
        # "forward" deployment the first IP address is used.
        deployment: reverse
        # Header name where the actual IP address will be reported, if more
        # than one IP address is present, the last IP address will be the
        # one taken into consideration.
        header: X-Forwarded-For

  # a line based log of HTTP requests (no alerts)
  - http-log:
      enabled: no
      filename: http.log
      append: yes
      #extended: yes     # enable this for extended logging information
      #custom: yes       # enabled the custom logging format (defined by customformat)
      #customformat: "%{%D-%H:%M:%S}t.%z %{X-Forwarded-For}i %H %m %h %u %s %B %a:%p -> %A:%P"
      #filetype: regular # 'regular', 'unix_stream' or 'unix_dgram'

  # a line based log of TLS handshake parameters (no alerts)
  - tls-log:
      enabled: no  # Log TLS connections.
      filename: tls.log # File to store TLS logs.
      append: yes
      #filetype: regular # 'regular', 'unix_stream' or 'unix_dgram'
      #extended: yes # Log extended information like fingerprint

  # output module to store certificates chain to disk
  - tls-store:
      enabled: no
      #certs-log-dir: certs # directory to store the certificates files

  # a line based log of DNS requests and/or replies (no alerts)
  - dns-log:
      enabled: no
      filename: dns.log
      append: yes
      #filetype: regular # 'regular', 'unix_stream' or 'unix_dgram'

  # Packet log... log packets in pcap format. 3 modes of operation: "normal"
  # "multi" and "sguil".
  #
  # In normal mode a pcap file "filename" is created in the default-log-dir,
  # or are as specified by "dir".
  # In multi mode, a file is created per thread. This will perform much
  # better, but will create multiple files where 'normal' would create one.
  # In multi mode the filename takes a few special variables:
  # - %n -- thread number
  # - %i -- thread id
  # - %t -- timestamp (secs or secs.usecs based on 'ts-format'
  # E.g. filename: pcap.%n.%t
  #
  # Note that it's possible to use directories, but the directories are not
  # created by Suricata. E.g. filename: pcaps/%n/log.%s will log into the
  # per thread directory.
  #
  # Also note that the limit and max-files settings are enforced per thread.
  # So the size limit when using 8 threads with 1000mb files and 2000 files
  # is: 8*1000*2000 ~ 16TiB.
  #
  # In Sguil mode "dir" indicates the base directory. In this base dir the
  # pcaps are created in th directory structure Sguil expects:
  #
  # $sguil-base-dir/YYYY-MM-DD/$filename.<timestamp>
  #
  # By default all packets are logged except:
  # - TCP streams beyond stream.reassembly.depth
  # - encrypted streams after the key exchange
  #
  - pcap-log:
      enabled: yes
      filename: log-%t-%n.pcap
      #filename: log.pcap

      # File size limit.  Can be specified in kb, mb, gb.  Just a number
      # is parsed as bytes.
      limit: 10mb

      # If set to a value will enable ring buffer mode. Will keep Maximum of "max-files" of size "limit"
      max-files: 20

      mode: multi # normal, multi or sguil.

      # Directory to place pcap files. If not provided the default log
      # directory will be used. Required for "sguil" mode.
      dir: /var/log/suricata/fpc/

      #ts-format: usec # sec or usec second format (default) is filename.sec usec is filename.sec.usec
      use-stream-depth: no #If set to "yes" packets seen after reaching stream inspection depth are ignored. "no" logs all packets
      honor-pass-rules: no # If set to "yes", flows in which a pass rule matched will stopped being logged.

  # a full alerts log containing much information for signature writers
  # or for investigating suspected false positives.
  - alert-debug:
      enabled: no
      filename: alert-debug.log
      append: yes
      #filetype: regular # 'regular', 'unix_stream' or 'unix_dgram'

  # alert output to prelude (http://www.prelude-technologies.com/) only
  # available if Suricata has been compiled with --enable-prelude
  - alert-prelude:
      enabled: no
      profile: suricata
      log-packet-content: no
      log-packet-header: yes

  # Stats.log contains data from various counters of the suricata engine.
  - stats:
      enabled: yes
      filename: stats.log
      totals: yes       # stats for all threads merged together
      threads: no       # per thread stats
      #null-values: yes  # print counters that have value 0

  # a line based alerts log similar to fast.log into syslog
  - syslog:
      enabled: no
      # reported identity to syslog. If ommited the program name (usually
      # suricata) will be used.
      #identity: "suricata"
      facility: local5
      #level: Info ## possible levels: Emergency, Alert, Critical,
                   ## Error, Warning, Notice, Info, Debug

  # a line based information for dropped packets in IPS mode
  - drop:
      enabled: no
      filename: drop.log
      append: yes
      #filetype: regular # 'regular', 'unix_stream' or 'unix_dgram'

  # Output module for storing files on disk. Files are stored in a
  # directory names consisting of the first 2 characters of the
  # SHA256 of the file. Each file is given its SHA256 as a filename.
  #
  # When a duplicate file is found, the existing file is touched to
  # have its timestamps updated.
  #
  # Unlike the older filestore, metadata is not written out by default
  # as each file should already have a "fileinfo" record in the
  # eve.log. If write-fileinfo is set to yes, the each file will have
  # one more associated .json files that consists of the fileinfo
  # record. A fileinfo file will be written for each occurrence of the
  # file seen using a filename suffix to ensure uniqueness.
  #
  # To prune the filestore directory see the "suricatactl filestore
  # prune" command which can delete files over a certain age.
  - file-store:
      version: 2
      enabled: no

      # Set the directory for the filestore. If the path is not
      # absolute will be be relative to the default-log-dir.
      #dir: filestore

      # Write out a fileinfo record for each occurrence of a
      # file. Disabled by default as each occurrence is already logged
      # as a fileinfo record to the main eve-log.
      #write-fileinfo: yes

      # Force storing of all files. Default: no.
      #force-filestore: yes

      # Override the global stream-depth for sessions in which we want
      # to perform file extraction. Set to 0 for unlimited.
      #stream-depth: 0

      # Uncomment the following variable to define how many files can
      # remain open for filestore by Suricata. Default value is 0 which
      # means files get closed after each write
      #max-open-files: 1000

      # Force logging of checksums, available hash functions are md5,
      # sha1 and sha256. Note that SHA256 is automatically forced by
      # the use of this output module as it uses the SHA256 as the
      # file naming scheme.
      #force-hash: [sha1, md5]

  # output module to store extracted files to disk (old style, deprecated)
  #
  # The files are stored to the log-dir in a format "file.<id>" where <id> is
  # an incrementing number starting at 1. For each file "file.<id>" a meta
  # file "file.<id>.meta" is created. Before they are finalized, they will
  # have a ".tmp" suffix to indicate that they are still being processed.
  #
  # If include-pid is yes, then the files are instead "file.<pid>.<id>", with
  # meta files named as "file.<pid>.<id>.meta"
  #
  # File extraction depends on a lot of things to be fully done:
  # - file-store stream-depth. For optimal results, set this to 0 (unlimited)
  # - http request / response body sizes. Again set to 0 for optimal results.
  # - rules that contain the "filestore" keyword.
  - file-store:
      enabled: no       # set to yes to enable
      log-dir: files    # directory to store the files
      force-magic: no   # force logging magic on all stored files
      # force logging of checksums, available hash functions are md5,
      # sha1 and sha256
      #force-hash: [md5]
      force-filestore: no # force storing of all files
      # override global stream-depth for sessions in which we want to
      # perform file extraction. Set to 0 for unlimited.
      #stream-depth: 0
      #waldo: file.waldo # waldo file to store the file_id across runs
      # uncomment to disable meta file writing
      #write-meta: no
      # uncomment the following variable to define how many files can
      # remain open for filestore by Suricata. Default value is 0 which
      # means files get closed after each write
      #max-open-files: 1000
      include-pid: no # set to yes to include pid in file names

  # output module to log files tracked in a easily parsable json format
  - file-log:
      enabled: no
      filename: files-json.log
      append: yes
      #filetype: regular # 'regular', 'unix_stream' or 'unix_dgram'

      force-magic: yes   # force logging magic on all logged files
      # force logging of checksums, available hash functions are md5,
      # sha1 and sha256
      force-hash: [md5, sha1, sha256]

  # Log TCP data after stream normalization
  # 2 types: file or dir. File logs into a single logfile. Dir creates
  # 2 files per TCP session and stores the raw TCP data into them.
  # Using 'both' will enable both file and dir modes.
  #
  # Note: limited by stream.depth
  - tcp-data:
      enabled: no
      type: file
      filename: tcp-data.log

  # Log HTTP body data after normalization, dechunking and unzipping.
  # 2 types: file or dir. File logs into a single logfile. Dir creates
  # 2 files per HTTP session and stores the normalized data into them.
  # Using 'both' will enable both file and dir modes.
  #
  # Note: limited by the body limit settings
  - http-body-data:
      enabled: no
      type: file
      filename: http-data.log

  # Lua Output Support - execute lua script to generate alert and event
  # output.
  # Documented at:
  # https://redmine.openinfosecfoundation.org/projects/suricata/wiki/Lua_Output
  - lua:
      enabled: no
      #scripts-dir: /etc/suricata/lua-output/
      scripts:
      #   - script1.lua

# Logging configuration.  This is not about logging IDS alerts/events, but
# output about what Suricata is doing, like startup messages, errors, etc.
logging:
  # The default log level, can be overridden in an output section.
  # Note that debug level logging will only be emitted if Suricata was
  # compiled with the --enable-debug configure option.
  #
  # This value is overridden by the SC_LOG_LEVEL env var.
  default-log-level: notice

  # The default output format.  Optional parameter, should default to
  # something reasonable if not provided.  Can be overriden in an
  # output section.  You can leave this out to get the default.
  #
  # This value is overridden by the SC_LOG_FORMAT env var.
  #default-log-format: "[%i] %t - (%f:%l) <%d> (%n) -- "

  # A regex to filter output.  Can be overridden in an output section.
  # Defaults to empty (no filter).
  #
  # This value is overridden by the SC_LOG_OP_FILTER env var.
  default-output-filter:

  # Define your logging outputs.  If none are defined, or they are all
  # disabled you will get the default - console output.
  outputs:
  - console:
      enabled: yes
      # type: json
  - file:
      enabled: yes
      level: info
      filename: /var/log/suricata/suricata.log
      # type: json
  - syslog:
      enabled: no
      facility: local5
      format: "[%i] <%d> -- "
      # type: json


##
## Step 4: configure common capture settings
##
## See "Advanced Capture Options" below for more options, including NETMAP
## and PF_RING.
##


##
## Step 5: App Layer Protocol Configuration
##

# Configure the app-layer parsers. The protocols section details each
# protocol.
#
# The option "enabled" takes 3 values - "yes", "no", "detection-only".
# "yes" enables both detection and the parser, "no" disables both, and
# "detection-only" enables protocol detection only (parser disabled).
app-layer:
  protocols:
    rfb:
      enabled: yes
      detection-ports:
        dp: 5900, 5901, 5902, 5903, 5904, 5905, 5906, 5907, 5908, 5909
    krb5:
      enabled: yes
    snmp:
      enabled: yes
    ikev2:
      enabled: yes
    tls:
      enabled: yes
      detection-ports:
        dp: 443

      # Generate JA3 fingerprint from client hello
      ja3-fingerprints: yes

      # What to do when the encrypted communications start:
      # - default: keep tracking TLS session, check for protocol anomalies,
      #            inspect tls_* keywords. Disables inspection of unmodified
      #            'content' signatures.
      # - bypass:  stop processing this flow as much as possible. No further
      #            TLS parsing and inspection. Offload flow bypass to kernel
      #            or hardware if possible.
      # - full:    keep tracking and inspection as normal. Unmodified content
      #            keyword signatures are inspected as well.
      #
      # For best performance, select 'bypass'.
      #
      #encrypt-handling: default

    dcerpc:
      enabled: yes
    ftp:
      enabled: yes
      # memcap: 64mb
    rdp:
      enabled: yes
    ssh:
      enabled: yes
    smtp:
      enabled: yes
      # Configure SMTP-MIME Decoder
      mime:
        # Decode MIME messages from SMTP transactions
        # (may be resource intensive)
        # This field supercedes all others because it turns the entire
        # process on or off
        decode-mime: yes

        # Decode MIME entity bodies (ie. base64, quoted-printable, etc.)
        decode-base64: yes
        decode-quoted-printable: yes

        # Maximum bytes per header data value stored in the data structure
        # (default is 2000)
        header-value-depth: 2000

        # Extract URLs and save in state data structure
        extract-urls: yes
        # Set to yes to compute the md5 of the mail body. You will then
        # be able to journalize it.
        body-md5: no
      # Configure inspected-tracker for file_data keyword
      inspected-tracker:
        content-limit: 100000
        content-inspect-min-size: 32768
        content-inspect-window: 4096
    imap:
      enabled: detection-only
    msn:
      enabled: detection-only
    # Note: --enable-rust is required for full SMB1/2 support. W/o rust
    # only minimal SMB1 support is available.
    smb:
      enabled: yes
      detection-ports:
        dp: 139, 445
    # Note: NFS parser depends on Rust support: pass --enable-rust
    # to configure.
    nfs:
      enabled: yes
    tftp:
      enabled: yes
    sip:
      enabled: yes
    dhcp:
      enabled: yes
    dns:
      # memcaps. Globally and per flow/state.
      #global-memcap: 16mb
      #state-memcap: 512kb

      # How many unreplied DNS requests are considered a flood.
      # If the limit is reached, app-layer-event:dns.flooded; will match.
      #request-flood: 500

      tcp:
        enabled: yes
        detection-ports:
          dp: 53
      udp:
        enabled: yes
        detection-ports:
          dp: 53
    http:
      enabled: yes
      # memcap: 64mb

      # default-config:           Used when no server-config matches
      #   personality:            List of personalities used by default
      #   request-body-limit:     Limit reassembly of request body for inspection
      #                           by http_client_body & pcre /P option.
      #   response-body-limit:    Limit reassembly of response body for inspection
      #                           by file_data, http_server_body & pcre /Q option.
      #   double-decode-path:     Double decode path section of the URI
      #   double-decode-query:    Double decode query section of the URI
      #   response-body-decompress-layer-limit:
      #                           Limit to how many layers of compression will be
      #                           decompressed. Defaults to 2.
      #
      # server-config:            List of server configurations to use if address matches
      #   address:                List of IP addresses or networks for this block
      #   personalitiy:           List of personalities used by this block
      #   request-body-limit:     Limit reassembly of request body for inspection
      #                           by http_client_body & pcre /P option.
      #   response-body-limit:    Limit reassembly of response body for inspection
      #                           by file_data, http_server_body & pcre /Q option.
      #   double-decode-path:     Double decode path section of the URI
      #   double-decode-query:    Double decode query section of the URI
      #
      #   uri-include-all:        Include all parts of the URI. By default the
      #                           'scheme', username/password, hostname and port
      #                           are excluded. Setting this option to true adds
      #                           all of them to the normalized uri as inspected
      #                           by http_uri, urilen, pcre with /U and the other
      #                           keywords that inspect the normalized uri.
      #                           Note that this does not affect http_raw_uri.
      #                           Also, note that including all was the default in
      #                           1.4 and 2.0beta1.
      #
      #   meta-field-limit:       Hard size limit for request and response size
      #                           limits. Applies to request line and headers,
      #                           response line and headers. Does not apply to
      #                           request or response bodies. Default is 18k.
      #                           If this limit is reached an event is raised.
      #
      # Currently Available Personalities:
      #   Minimal, Generic, IDS (default), IIS_4_0, IIS_5_0, IIS_5_1, IIS_6_0,
      #   IIS_7_0, IIS_7_5, Apache_2
      libhtp:
         default-config:
           personality: IDS

           # Can be specified in kb, mb, gb.  Just a number indicates
           # it's in bytes.
           request-body-limit: 100kb
           response-body-limit: 100kb

           # inspection limits
           request-body-minimal-inspect-size: 32kb
           request-body-inspect-window: 4kb
           response-body-minimal-inspect-size: 40kb
           response-body-inspect-window: 16kb

           # response body decompression (0 disables)
           response-body-decompress-layer-limit: 2

           # auto will use http-body-inline mode in IPS mode, yes or no set it statically
           http-body-inline: auto

           # Decompress SWF files.
           # 2 types: 'deflate', 'lzma', 'both' will decompress deflate and lzma
           # compress-depth:
           # Specifies the maximum amount of data to decompress,
           # set 0 for unlimited.
           # decompress-depth:
           # Specifies the maximum amount of decompressed data to obtain,
           # set 0 for unlimited.
           swf-decompression:
             enabled: yes
             type: both
             compress-depth: 0
             decompress-depth: 0

           # Take a random value for inspection sizes around the specified value.
           # This lower the risk of some evasion technics but could lead
           # detection change between runs. It is set to 'yes' by default.
           #randomize-inspection-sizes: yes
           # If randomize-inspection-sizes is active, the value of various
           # inspection size will be choosen in the [1 - range%, 1 + range%]
           # range
           # Default value of randomize-inspection-range is 10.
           #randomize-inspection-range: 10

           # decoding
           double-decode-path: no
           double-decode-query: no

         server-config:

           #- apache:
           #    address: [192.168.1.0/24, 127.0.0.0/8, "::1"]
           #    personality: Apache_2
           #    # Can be specified in kb, mb, gb.  Just a number indicates
           #    # it's in bytes.
           #    request-body-limit: 4096
           #    response-body-limit: 4096
           #    double-decode-path: no
           #    double-decode-query: no

           #- iis7:
           #    address:
           #      - 192.168.0.0/24
           #      - 192.168.10.0/24
           #    personality: IIS_7_0
           #    # Can be specified in kb, mb, gb.  Just a number indicates
           #    # it's in bytes.
           #    request-body-limit: 4096
           #    response-body-limit: 4096
           #    double-decode-path: no
           #    double-decode-query: no

    # Note: Modbus probe parser is minimalist due to the poor significant field
    # Only Modbus message length (greater than Modbus header length)
    # And Protocol ID (equal to 0) are checked in probing parser
    # It is important to enable detection port and define Modbus port
    # to avoid false positive
    modbus:
      # How many unreplied Modbus requests are considered a flood.
      # If the limit is reached, app-layer-event:modbus.flooded; will match.
      #request-flood: 500

      enabled: yes
      detection-ports:
        dp: 502
      # According to MODBUS Messaging on TCP/IP Implementation Guide V1.0b, it
      # is recommended to keep the TCP connection opened with a remote device
      # and not to open and close it for each MODBUS/TCP transaction. In that
      # case, it is important to set the depth of the stream reassembling as
      # unlimited (stream.reassembly.depth: 0)

      # Stream reassembly size for modbus. By default track it completely.
      stream-depth: 0

    # DNP3
    dnp3:
      enabled: yes
      detection-ports:
        dp: 20000

    # SCADA EtherNet/IP and CIP protocol support
    enip:
      enabled: yes
      detection-ports:
        dp: 44818
        sp: 44818

    # Note: parser depends on experimental Rust support
    ntp:
      enabled: yes

# Limit for the maximum number of asn1 frames to decode (default 256)
asn1-max-frames: 256
`
