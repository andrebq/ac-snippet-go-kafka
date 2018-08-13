#!/bin/bash

messageID=$(curl -X POST -F "Text='$3'"  -u "$1:$1" http://localhost:9099/channels/$2/message 2>/dev/null)
echo "$1@$2 [$messageID]: $3"