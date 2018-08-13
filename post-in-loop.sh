#!/bin/bash

for i in $(seq 1 $3); do

    messageID=$(curl -X POST -F "Text='Message no.: $i'"  -u "$1:$1" http://localhost:9099/channels/$2/message 2>/dev/null)
    echo "$1@$2 [$messageID] Message no.: $i"

done