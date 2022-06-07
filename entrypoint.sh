#!/bin/sh
set -e

RUN_CMD="lenpaste"


# LENPASTE_ADDRESS
if [[ -n "$LENPASTE_ADDRESS" ]]; then
	RUN_CMD="$RUN_CMD -address $LENPASTE_ADDRESS"
fi


# LENPASTE_DB_DRIVER
if [[ -n "$LENPASTE_DB_DRIVER" ]]; then
	RUN_CMD="$RUN_CMD -db-driver $LENPASTE_DB_DRIVER"
fi


# LENPASTE_DB_SOURCE
if [[ -z "$LENPASTE_DB_DRIVER" || "$LENPASTE_DB_DRIVER" == "sqlite3" ]]; then
	RUN_CMD="$RUN_CMD -db-source /data/lenpaste.db"

else
	RUN_CMD="$RUN_CMD -db-source $LENPASTE_DB_SOURCE"
fi


# LENPASTE_ROBOTS_DISALLOW
if [[ "$LENPASTE_ROBOTS_DISALLOW" == "true" ]]; then
	RUN_CMD="$RUN_CMD -robots-disallow"
	
else
	if [[ "$LENPASTE_ROBOTS_DISALLOW" != "" && "$LENPASTE_ROBOTS_DISALLOW" != "false" ]]; then
		echo "[ENTRYPOINT] Error: unknown: LENPASTE_ROBOTS_DISALLOW = $LENPASTE_ROBOTS_DISALLOW"
		exit 2
	fi
fi


# Run Lenpaste
echo "[ENTRYPOINT] $RUN_CMD"
$RUN_CMD
