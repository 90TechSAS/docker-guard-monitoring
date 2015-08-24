#!/bin/bash

# Configuration
HOOK_URL="https://hooks.slack.com/services/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"
CHANNEL="#mychannel"
USERNAME="DGS"
ICON=":squirrel:"

# Check arguments
if [ $# -lt 5 ] || [ $# -gt 5 ] ; then
        echo "Error, usage: slack.sh <severity> <type> <target> <target_probe> <data>"
        echo "severity	Severity levels:"
        echo "				0: Notice"
        echo "				1: Warning"
        echo "				2: Critical"
        echo 
        echo "Type 		Alert types:"
        echo "				DiskSpaceLimitReached"
        echo "				MemorySpaceLimitReached"
        echo "				ContainerStarted"
        echo "				ContainerStopped"
        echo "				ContainerRemoved"
        echo "				DiskIOOverload"
        echo "				NetBandwithOverload"
        echo "				CPUUsageOverload"
        echo 
        echo "target		Targeted system(s)"
        echo 
        echo "target_probe	Name of the target's probe"
        echo
        echo "data		Additional data"
        echo
        echo "Example: slack.sh 3 DiskSpaceLimitReached db-1 production '20GB/20GB'"
        echo
        echo "Used: $0 $@"
        exit 1
fi

# Set variables
COLOR=""
SEVERITY=""
COLOR_NOTICE="#05c1ff"
COLOR_WARNING="#ffff00"
COLOR_CRITICAL="#ff0000"

case "$1" in
	0)
		SEVERITY="NOTICE"
		COLOR=$COLOR_NOTICE		
		;;
	1)
		SEVERITY="WARNING"
		COLOR=$COLOR_WARNING
		;;
	2)
		SEVERITY="CRITICAL"
		COLOR=$COLOR_CRITICAL
		;;
	*)
		echo "Error: Severity unknow"
		exit 1
		;;
esac

TEXT="\"attachments\": [{
            \"pretext\": \"New $SEVERITY $2 alert:\",
            \"text\": \"\",
            \"fields\": [
                {
                    \"title\": \"Target(s)\",
                    \"value\": \"$3\"
                },
                {
                    \"title\": \"Probe\",
                    \"value\": \"$4\"
                },
                {
                    \"title\": \"Additional data\",
                    \"value\": \"$5\"
                }
            ],
            \"color\": \"$COLOR\"
        }
    ]"

PAYLOAD="{\"channel\":\"$CHANNEL\",\"username\":\"$USERNAME\",$TEXT,\"icon_emoji\":\"$ICON\"}"

# Do the HTTP POST request to slack
curl -X POST \
	--data-urlencode "payload=$PAYLOAD" \
	$HOOK_URL
