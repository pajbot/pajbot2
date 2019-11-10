#!/usr/bin/env bash
PBUID=$(id -u pajbot)
PBGID=$(id -g pajbot)

if [[ -z "${PBUID}${PBGID}" ]]; then
  echo 'pajbot user not detected.'
  exit 1
fi

if [ ! -f /opt/pajbot/configs/pajbot2.json ]; then
    echo "No config file /opt/pajbot/configs/pajbot2.json found."
    exit 1
fi

docker run \
--name pajbot2 \
--network host \
--restart unless-stopped \
-d \
-v /opt/pajbot/configs/pajbot2.json:/app/cmd/bot/config.json \
-v /var/run/postgresql:/var/run/postgresql:ro \
-v /etc/localtime:/etc/localtime:ro \
-u "$PBUID":"$PBGID" \
pajbot2:latest
