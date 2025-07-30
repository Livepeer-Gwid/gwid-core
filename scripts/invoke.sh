#!/bin/bash

grafana_url="https://gwid.io/grafana"

json_response() {
	echo "{\"status\": \"$1\", \"message\": \"$2\", \"data\": \"$3\"}"
}

error_message() {
	json_response "error" "server init failed."
	exit 1
}

error_message

# success_message() {
# 	json_response "success" "Grafana dashboard URL: $grafana_url" $grafana_url
# }
#
# success_message

#
#
# json_response "success" "Deployment process completed."
