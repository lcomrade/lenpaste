#!/bin/sh
set -e

fatal_comp() {
	>&2 echo "$1"
	>&2 echo ""
	>&2 echo "Yes, I know that breaking backward compatibility is a bad idea. I'll try not to do it again:)"
	>&2 echo ""
	>&2 echo "More on what other backward compatibility was broken in Lenpaste v1.4 can be found here: https://lcomrade.su/rfQCH"
	exit 1
}

fatal_comp_env() {
	fatal_comp "$1 environment variable is no longer supported."
}

fatal_comp_l10n() {
	fatal_comp "A single \"$1\" file is no longer supported. The \"$1\" file must be split into several files for different locales: \"$1/en.txt\", \"$1/ru.txt\" and etc."
}

# LENPASTE_ADDRESS
if [ -n "$LENPASTE_ADDRESS" ]; then
	fatal_comp_env "LENPASTE_ADDRESS"
fi


# LENPASTE_DB_DRIVER
if [ -n "$LENPASTE_DB_DRIVER" ]; then
	fatal_comp_env "LENPASTE_DB_DRIVER"
fi


# LENPASTE_DB_SOURCE
if [ -n "$LENPASTE_DB_SOURCE" ]; then
	fatal_comp_env "LENPASTE_DB_SOURCE"
fi


# LENPASTE_DB_MAX_OPEN_CONNS
if [ -n "$LENPASTE_DB_MAX_OPEN_CONNS" ]; then
	fatal_comp_env "LENPASTE_DB_MAX_OPEN_CONNS"
fi


# LENPASTE_DB_MAX_IDLE_CONNS
if [ -n "$LENPASTE_DB_MAX_IDLE_CONNS" ]; then
	fatal_comp_env "LENPASTE_DB_MAX_IDLE_CONNS"
fi


# LENPASTE_DB_CLEANUP_PERIOD
if [ -n "$LENPASTE_DB_CLEANUP_PERIOD" ]; then
	fatal_comp_env "LENPASTE_DB_CLEANUP_PERIOD"
fi

# LENPASTE_ROBOTS_DISALLOW
if [ -n "$LENPASTE_ROBOTS_DISALLOW" ]; then
	fatal_comp_env "LENPASTE_ROBOTS_DISALLOW"
fi


# LENPASTE_TITLE_MAX_LENGTH
if [ -n "$LENPASTE_TITLE_MAX_LENGTH" ]; then
	fatal_comp_env "LENPASTE_TITLE_MAX_LENGTH"
fi


# LENPASTE_BODY_MAX_LENGTH
if [ -n "$LENPASTE_BODY_MAX_LENGTH" ]; then
	fatal_comp_env "LENPASTE_BODY_MAX_LENGTH"
fi


# LENPASTE_MAX_PASTE_LIFETIME
if [ -n "$LENPASTE_MAX_PASTE_LIFETIME" ]; then
	fatal_comp_env "LENPASTE_MAX_PASTE_LIFETIME"
fi

# Rate limits to get
if [ -n "$LENPASTE_GET_PASTES_PER_5MIN" ]; then
	fatal_comp_env "LENPASTE_GET_PASTES_PER_5MIN"
fi

if [ -n "$LENPASTE_GET_PASTES_PER_15MIN" ]; then
	fatal_comp_env "LENPASTE_GET_PASTES_PER_15MIN"
fi

if [ -n "$LENPASTE_GET_PASTES_PER_1HOUR" ]; then
	fatal_comp_env "LENPASTE_GET_PASTES_PER_1HOUR"
fi


# Rate limits to create
if [ -n "$LENPASTE_NEW_PASTES_PER_5MIN" ]; then
	fatal_comp_env "LENPASTE_NEW_PASTES_PER_5MIN"
fi

if [ -n "$LENPASTE_NEW_PASTES_PER_15MIN" ]; then
	fatal_comp_env "LENPASTE_NEW_PASTES_PER_15MIN"
fi

if [ -n "$LENPASTE_NEW_PASTES_PER_1HOUR" ]; then
	fatal_comp_env "LENPASTE_NEW_PASTES_PER_1HOUR"
fi



# Server about
if [ -f "/data/about" ]; then
	fatal_comp_l10n "about"
fi


# Server rules
if [ -f "/data/rules" ]; then
	fatal_comp_l10n "rules"
fi


# Server terms of use
if [ -f "/data/terms" ]; then
	fatal_comp_l10n "terms"
fi


# LENPASTE_ADMIN_NAME
if [ -n "$LENPASTE_ADMIN_NAME" ]; then
	fatal_comp_env "LENPASTE_ADMIN_NAME"
fi


# LENPASTE_ADMIN_MAIL
if [ -n "$LENPASTE_ADMIN_MAIL" ]; then
	fatal_comp_env "LENPASTE_ADMIN_MAIL"
fi


# LENPASTE_UI_DEFAULT_LIFETIME
if [ -n "$LENPASTE_UI_DEFAULT_LIFETIME" ]; then
	fatal_comp_env "LENPASTE_UI_DEFAULT_LIFETIME"
fi


# LENPASTE_UI_DEFAULT_THEME
if [ -n "$LENPASTE_UI_DEFAULT_THEME" ]; then
	fatal_comp_env "LENPASTE_UI_DEFAULT_THEME"
fi


# External UI themes
# if [ -d "/data/themes" ]; then
# 	# Themes must work normal.
# 	# No alert needed.
# fi


# Lenpsswd file
if [ -f "/data/lenpasswd" ]; then
	fatal_comp "The \"lenpasswd\" file is no longer supported."
fi


# Run Lenpaste
lenpaste --cfg-dir /etc/lenpaste/ run
