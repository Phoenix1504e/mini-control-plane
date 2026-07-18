#!/usr/bin/env bash
set -e

CAMPAIGN_NAME="campaign-$1"
MITIGATION_STRATEGY=$1

if [ -z "$MITIGATION_STRATEGY" ]; then
    echo "Usage: $0 <mitigation-strategy>"
    exit 1
fi

echo "========================================="
echo "Starting MVCC Experimental Campaign"
echo "Campaign: $CAMPAIGN_NAME | Strategy: $MITIGATION_STRATEGY"
echo "========================================="

# 1. Automate upfront project compilation into dedicated target spaces
echo "Compiling system control plane components cleanly..."
mkdir -p bin tools/bin

go build -o ./bin/apiserver ./cmd/apiserver/*.go 2>/dev/null || go build -o ./bin/apiserver ./cmd/apiserver/main.go
go build -o ./bin/controller ./cmd/controller/*.go 2>/dev/null || go build -o ./bin/controller ./cmd/controller/main.go
go build -o ./tools/bin/runner ./tools/runner/main.go

chmod +x ./bin/apiserver ./bin/controller ./tools/bin/runner
echo "Compilation complete. Executables verified."

# Define evaluation scaling dimensions for concurrent conflict simulation
CONTROLLER_SCENARIOS=(1 2 4 8)
TRIALS_PER_SCENARIO=5

for c_count in "${CONTROLLER_SCENARIOS[@]}"; do
    echo ""
    echo "#########################################"
    echo "Controllers = $c_count"
    echo "#########################################"
    
    for trial in $(seq 1 $TRIALS_PER_SCENARIO); do
        echo "-----------------------------------------"
        echo "Trial $trial / $TRIALS_PER_SCENARIO | Controllers = $c_count"
        echo "-----------------------------------------"
        
        echo "Cleaning background runtime processes and stale data..."
        pkill -9 etcd || true
        pkill -9 apiserver || true
        pkill -9 controller || true
        rm -rf default.etcd
        
        echo "Clean state prepared. Launching runner..."
        # Execute the cleanly compiled orchestration runner artifact directly
        ./tools/bin/runner \
            -scenario "$CAMPAIGN_NAME" \
            -controllers "$c_count" \
            -resources 20 \
            -duration 15s \
            -mitigation "$MITIGATION_STRATEGY"
            
        echo "Trial $trial finished successfully."
        sleep 2
    done
done

echo ""
echo "========================================="
echo "MVCC Experimental Campaign Complete!"
echo "Data archived cleanly inside: experiments/$CAMPAIGN_NAME"
echo "========================================="
