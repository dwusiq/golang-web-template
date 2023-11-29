#!/bin/bash
SERVER_NAME="trx_dev"
CHAIN_ID="2494104990"

cd /data/walletSyn_$CHAIN_ID/walletsynv2/
git pull ssh://git@gitlab.thyy.pro:5519/blockchain/walletsynv2.git dev
make walletSyn

while read -r pid; do
  if [[ -n "$pid" ]]; then
    kill -9 "$pid"
    echo "PID $pid is killed"
  fi
done <<< $(ps -ef | grep ./build/bin/$SERVER_NAME | grep -v grep | awk '{print $2}')

screen -dmS syn_$CHAIN_ID bash -c "cd /data/walletSyn_$CHAIN_ID/walletsynv2/ && ./build/bin/$SERVER_NAME; exec bash"

screen -wipe