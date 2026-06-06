#!/bin/bash

set -e

CAMPAIGN="campaign-final"

CONTROLLER_SET=(1 2 4 8)

TRIALS=5

RESOURCES=10

DURATION=20s

echo ""
echo "========================================="
echo "Starting MVCC Experimental Campaign"
echo "Campaign: $CAMPAIGN"
echo "========================================="
echo ""

mkdir -p experiments/mvcc-conflicts/$CAMPAIGN

for controllers in "${CONTROLLER_SET[@]}"
do

    echo ""
    echo "#########################################"
    echo "Controllers = $controllers"
    echo "#########################################"
    echo ""

    for trial in $(seq 1 $TRIALS)
    do

        echo ""
        echo "-----------------------------------------"
        echo "Trial $trial / $TRIALS"
        echo "Controllers = $controllers"
        echo "-----------------------------------------"
        echo ""

        # Stop old processes safely
	pkill controller || true
	sleep 2

	# Clean old resources from etcd
	etcdctl del /resources --prefix

        # Remove old telemetry
        rm -f *.jsonl

        sleep 2

	echo "Clean state prepared"

	# Run experiment with watchdog timeout
	timeout 90s go run ./tools/runner \
    		--scenario mvcc-conflicts/$CAMPAIGN \
    		--controllers $controllers \
    		--resources $RESOURCES \
    		--duration $DURATION
	pkill controller || true

	sleep 5
        echo ""
        echo "Trial completed"
        echo ""

        # cooldown between trials
        sleep 5

    done

done

echo ""
echo "========================================="
echo "All experimental trials completed"
echo "========================================="
echo ""
