#!/bin/bash

# Check if jq is installed
if ! dpkg-query -W -f='${Status}' jq | grep -q "install ok installed"; then
  echo "jq not found, installing..."
  sudo apt-get update
  sudo apt-get install -y jq
fi

# Check if bc is installed
if ! dpkg-query -W -f='${Status}' bc | grep -q "install ok installed"; then
  echo "bc not found, installing..."
  sudo apt-get update
  sudo apt-get install -y bc
fi

echo '------------------DEVNET TPS-------------------'
echo $(date)
echo

REMOTE_RPC="https://fullnode.devnet.sui.io:443"
LOCAL_RPC="127.0.0.1:9000"

function get_transactions {
  result=$(curl -m 2 --location --request POST $1 \
  --header 'Content-Type: application/json' \
  --data-raw '{ "jsonrpc":"2.0", "method":"sui_getTotalTransactionBlocks","id":1}' 2>/dev/null | jq .result)

    if [ -z "$result" ]; then
      echo "Error: Failed to extract 'result' field from JSON response."
      return 1
    fi

  echo $result
}

function calculate_tps {
  local start=$1
  local end=$2

  if [[ "$start" =~ ^[0-9]+(\.[0-9]+)?$ && "$end" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
    start=$(echo "$start" | bc)
    end=$(echo "$end" | bc)
    tps=$(bc <<< "scale=2; ($end - $start) / 10")
    if [[ $? -ne 0 || -z "$tps" ]]; then
      echo "Failed to calculate TPS: error occurred during calculation or response format has changed"
      exit 1
    fi
    echo $tps
  else
    echo "Failed to calculate TPS: transaction data is not in the expected format or server is not running"
    exit 1
  fi
}

SUISTART=$(get_transactions $REMOTE_RPC)
if [ -z "$SUISTART" ]; then
    echo "Failed to retrieve transactions from $REMOTE_RPC: check if remote devnet RPC is up and running or if the response format has changed"
    exit 1
fi

NODESTART=$(get_transactions $LOCAL_RPC)
if [ -z "$NODESTART" ]; then
  echo "Failed to retrieve transactions from $LOCAL_RPC: check if your node is up and running on port 9000 or if the response format has changed"
  exit 1
fi

for I in {1..10}; do
  sleep 1
  BAR="$(yes . | head -n ${I} | tr -d '\n')"
  printf "\rIN PROGRESS [%3d/100] %s" $((I * 10)) ${BAR}
done

printf "\n\n"

SUIEND=$(get_transactions $REMOTE_RPC)
NODEEND=$(get_transactions $LOCAL_RPC)

if [[ -z "$SUISTART" || -z "$SUIEND" ]]; then
  echo "Failed to calculate SUI TPS: transaction data is missing or response format has changed"
  exit 1
fi

if [[ -z "$NODESTART" || -z "$NODEEND" ]]; then
  echo "Failed to calculate local TPS: transaction data is missing or response format has changed"
  exit 1
fi

SUITPS=$(calculate_tps "${SUISTART//\"/}" "${SUIEND//\"/}")
MYTPS=$(calculate_tps "${NODESTART//\"/}" "${NODEEND//\"/}")

echo 'SUI TPS: '$SUITPS
echo 'NODE TPS: '$MYTPS
echo '-----------------------------------------------'