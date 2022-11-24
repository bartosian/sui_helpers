#!/bin/bash
  
echo '-------------------------------------'
echo $(date)
echo 
SUISTART=$(curl --location --request POST https://fullnode.devnet.sui.io:443 \
--header 'Content-Type: application/json' \
--data-raw '{ "jsonrpc":"2.0", "method":"sui_getTotalTransactionNumber","id":1}' 2>/dev/null | jq .result)

NODESTART=$(curl --location --request POST 127.0.0.1:9000 \
--header 'Content-Type: application/json' \
--data-raw '{ "jsonrpc":"2.0", "method":"sui_getTotalTransactionNumber","id":1}' 2>/dev/null | jq .result)

for I in {1..10}; do
  sleep 1
  BAR="$(yes . | head -n ${I} | tr -d '\n')"
  printf "\r[%3d/100] %s" $((I * 10)) ${BAR}
done
printf "\n"

SUIEND=$(curl --location --request POST https://fullnode.devnet.sui.io:443 \
--header 'Content-Type: application/json' \
--data-raw '{ "jsonrpc":"2.0", "method":"sui_getTotalTransactionNumber","id":1}' 2>/dev/null | jq .result)

NODEEND=$(curl --location --request POST 127.0.0.1:9000 \
--header 'Content-Type: application/json' \
--data-raw '{ "jsonrpc":"2.0", "method":"sui_getTotalTransactionNumber","id":1}' 2>/dev/null | jq .result)

SUITPS=$((($SUIEND-$SUISTART)/10))
MYTPS=$((($NODEEND-$NODESTART)/10))

echo 'SUI TPS: '$SUITPS
echo 'NODE TPS: '$MYTPS
echo '-------------------------------------'
