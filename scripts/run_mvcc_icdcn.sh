#!/bin/bash
set -e

CAMPAIGN="campaign-icdcn"

CONTROLLER_SET=(1 2 4 8)
TRIALS=30

RESOURCES=10
DURATION=20s

mkdir -p experiments/mvcc-conflicts/$CAMPAIGN

for controllers in "${CONTROLLER_SET[@]}"
do
    echo ""
    echo "=================================="
    echo "Controllers = $controllers"
    echo "=================================="

    for trial in $(seq 1 $TRIALS)
    do
        echo ""
        echo "Trial $trial/$TRIALS"

        pkill controller || true
        sleep 2

        etcdctl del /resources --prefix >/dev/null

        rm -f *.jsonl

        sleep 2

        echo "Launching $controllers controllers..."

        controller_pids=()

        for i in $(seq 1 $controllers)
        do
            ./bin/controller &
            controller_pids+=($!)
        done

        sleep 5

        for i in $(seq 0 $((RESOURCES-1)))
        do
            payload="{\"metadata\":{\"name\":\"resource-$i\"},\"spec\":{\"replicas\":3}}"

            etcdctl put \
                "/resources/resource-$i" \
                "$payload" >/dev/null
        done

        sleep ${DURATION%s}

        for pid in "${controller_pids[@]}"
        do
            kill "$pid" 2>/dev/null || true
        done

        run_dir="experiments/mvcc-conflicts/$CAMPAIGN/run-$(date +%s)"

        mkdir -p "$run_dir"

        cp mvcc_conflicts.jsonl "$run_dir/" 2>/dev/null || true
        cp reconcile.jsonl "$run_dir/" 2>/dev/null || true
        cp leader_events.jsonl "$run_dir/" 2>/dev/null || true
        cp state_samples.jsonl "$run_dir/" 2>/dev/null || true

        cat > "$run_dir/metadata.json" <<EOF
{
  "controllers": $controllers,
  "resources": $RESOURCES,
  "duration": "$DURATION",
  "trial": $trial
}
EOF

        sleep 5
    done
done
