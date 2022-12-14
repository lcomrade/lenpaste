#!/bin/sh
set -e

readonly VERSION="1.0, 14.12.2022"

printHelp(){
	echo "Usage: $0 [OPTION]... [FILES]"
	echo "Check locale files."
	echo ""
	echo " -s, --slien     print only errors"
	echo "     --rm-empty  remove empty locale files"
	echo ""
	echo " -h, --help      show help and exit"
	echo " -v, --version   show version and exit"
}

# CLI args
ARG_SLIENT=false
ARG_RM_EMPTY=false

if [ -z "$1" ]; then
	echo "error: not files specified" 1>&2
	exit 1
fi

while [ -n "$1" ]; do
	case "$1" in
		-s|--slient)
		ARG_SLIENT=true
		;;

		--rm-empty)
		ARG_RM_EMPTY=true
		;;

		-h|--help)
		printHelp
		exit 0
		;;

		-v|--version)
		echo "$VERSION"
		exit 0
		;;

		*)
		break
		;;
	esac

	shift
done

# RUN
#find ./ -type f -name "*.locale" | while read -r file; do
#	if ! grep -q -Ev '^$' "$file"; then
#		echo "$file"
#	fi
#done

while [ -n "$1" ]; do
	# If empty
	if ! grep -q -Ev '^$' "$1"; then
		# If need remove empty files
		if [ $ARG_RM_EMPTY = true ]; then
			rm "$1"
			
			if [ $ARG_SLIENT = false ]; then
				echo "remove empty file: $1"
			fi

		# If only print errors
		else
			if [ $ARG_SLIENT = false ]; then
				echo "empty: $1"
			fi
		fi

	# If not ok
	else
		if [ $ARG_SLIENT = false ]; then
			echo "ok $1"
		fi
	fi	

	shift
done
