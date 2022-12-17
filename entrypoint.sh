#!/bin/sh
set -e

RUN_CMD="lenpaste"


# LENPASTE_ADDRESS
if [ -n "$LENPASTE_ADDRESS" ]; then
	RUN_CMD="$RUN_CMD -address $LENPASTE_ADDRESS"
fi


# LENPASTE_DB_DRIVER
if [ -n "$LENPASTE_DB_DRIVER" ]; then
	RUN_CMD="$RUN_CMD -db-driver '$LENPASTE_DB_DRIVER'"
fi


# LENPASTE_DB_SOURCE
if [ -z "$LENPASTE_DB_DRIVER" ] || [ "$LENPASTE_DB_DRIVER" = "sqlite3" ]; then
	RUN_CMD="$RUN_CMD -db-source /data/lenpaste.db"

else
	RUN_CMD="$RUN_CMD -db-source '$LENPASTE_DB_SOURCE'"
fi


# LENPASTE_DB_MAX_OPEN_CONNS
if [ -n "$LENPASTE_DB_MAX_OPEN_CONNS" ]; then
	RUN_CMD="$RUN_CMD -db-max-open-conns '$LENPASTE_DB_MAX_OPEN_CONNS'"
fi


# LENPASTE_DB_MAX_IDLE_CONNS
if [ -n "$LENPASTE_DB_MAX_IDLE_CONNS" ]; then
	RUN_CMD="$RUN_CMD -db-max-idle-conns '$LENPASTE_DB_MAX_IDLE_CONNS'"
fi


# LENPASTE_DB_CLEANUP_PERIOD
if [ -n "$LENPASTE_DB_CLEANUP_PERIOD" ]; then
	RUN_CMD="$RUN_CMD -db-cleanup-period '$LENPASTE_DB_CLEANUP_PERIOD'"
fi

# LENPASTE_ROBOTS_DISALLOW
if [ "$LENPASTE_ROBOTS_DISALLOW" = "true" ]; then
	RUN_CMD="$RUN_CMD -robots-disallow"
	
else
	if [ "$LENPASTE_ROBOTS_DISALLOW" != "" ] && [ "$LENPASTE_ROBOTS_DISALLOW" != "false" ]; then
		echo "[ENTRYPOINT] Error: unknown: LENPASTE_ROBOTS_DISALLOW = $LENPASTE_ROBOTS_DISALLOW"
		exit 2
	fi
fi


# LENPASTE_TITLE_MAX_LENGTH
if [ -n "$LENPASTE_TITLE_MAX_LENGTH" ]; then
	RUN_CMD="$RUN_CMD -title-max-length '$LENPASTE_TITLE_MAX_LENGTH'"
fi


# LENPASTE_BODY_MAX_LENGTH
if [ -n "$LENPASTE_BODY_MAX_LENGTH" ]; then
	RUN_CMD="$RUN_CMD -body-max-length '$LENPASTE_BODY_MAX_LENGTH'"
fi


# LENPASTE_MAX_PASTE_LIFETIME
if [ -n "$LENPASTE_MAX_PASTE_LIFETIME" ]; then
	RUN_CMD="$RUN_CMD -max-paste-lifetime '$LENPASTE_MAX_PASTE_LIFETIME'"
fi

# Rate limits to get
if [ -n "$LENPASTE_GET_PASTES_PER_5MIN" ]; then
	RUN_CMD="$RUN_CMD -get-pastes-per-5min '$LENPASTE_GET_PASTES_PER_5MIN'"
fi

if [ -n "$LENPASTE_GET_PASTES_PER_15MIN" ]; then
	RUN_CMD="$RUN_CMD -get-pastes-per-15min '$LENPASTE_GET_PASTES_PER_15MIN'"
fi

if [ -n "$LENPASTE_GET_PASTES_PER_1HOUR" ]; then
	RUN_CMD="$RUN_CMD -get-pastes-per-1hour '$LENPASTE_GET_PASTES_PER_1HOUR'"
fi


# Rate limits to create
if [ -n "$LENPASTE_NEW_PASTES_PER_5MIN" ]; then
	RUN_CMD="$RUN_CMD -new-pastes-per-5min '$LENPASTE_NEW_PASTES_PER_5MIN'"
fi

if [ -n "$LENPASTE_NEW_PASTES_PER_15MIN" ]; then
	RUN_CMD="$RUN_CMD -new-pastes-per-15min '$LENPASTE_NEW_PASTES_PER_15MIN'"
fi

if [ -n "$LENPASTE_NEW_PASTES_PER_1HOUR" ]; then
	RUN_CMD="$RUN_CMD -new-pastes-per-1hour '$LENPASTE_NEW_PASTES_PER_1HOUR'"
fi



# Server about
if [ -f "/data/about" ]; then
	RUN_CMD="$RUN_CMD -server-about /data/about"
fi


# Server rules
if [ -f "/data/rules" ]; then
	RUN_CMD="$RUN_CMD -server-rules /data/rules"
fi


# Server terms of use
if [ -f "/data/terms" ]; then
	RUN_CMD="$RUN_CMD -server-terms /data/terms"
fi


# LENPASTE_ADMIN_NAME
if [ -n "$LENPASTE_ADMIN_NAME" ]; then
	RUN_CMD="$RUN_CMD -admin-name '$LENPASTE_ADMIN_NAME'"
fi


# LENPASTE_ADMIN_MAIL
if [ -n "$LENPASTE_ADMIN_MAIL" ]; then
	RUN_CMD="$RUN_CMD -admin-mail '$LENPASTE_ADMIN_MAIL'"
fi


# LENPASTE_UI_DEFAULT_LIFETIME
if [ -n "$LENPASTE_UI_DEFAULT_LIFETIME" ]; then
	RUN_CMD="$RUN_CMD -ui-default-lifetime '$LENPASTE_UI_DEFAULT_LIFETIME'"
fi


# LENPASTE_UI_DEFAULT_THEME
if [ -n "$LENPASTE_UI_DEFAULT_THEME" ]; then
	RUN_CMD="$RUN_CMD -ui-default-theme $LENPASTE_UI_DEFAULT_THEME"
fi


# External UI themes
if [ -d "/data/themes" ]; then
	RUN_CMD="$RUN_CMD -ui-themes-dir /data/themes"
fi


# Lenpsswd file
if [ -f "/data/lenpasswd" ]; then
	RUN_CMD="$RUN_CMD -lenpasswd-file /data/lenpasswd"
fi


# Run Lenpaste
echo "[ENTRYPOINT] $RUN_CMD"
sh -c "$RUN_CMD"
