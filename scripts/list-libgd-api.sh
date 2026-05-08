#!/bin/sh
set -eu

headers="${LIBGD_HEADERS:-/opt/homebrew/include/gd.h /opt/homebrew/include/gdfx.h /usr/local/include/gd.h /usr/local/include/gdfx.h /usr/include/gd.h /usr/include/gdfx.h}"

for header in $headers; do
	if [ -f "$header" ]; then
		awk '
			/BGD_DECLARE/ {
				line=$0
				while (line !~ /;/ && getline nextline) {
					line=line " " nextline
				}
				if (match(line, /(gd[A-Za-z0-9_]+)[[:space:]]*\(/)) {
					name = substr(line, RSTART, RLENGTH - 1)
					gsub(/[[:space:]]+$/, "", name)
					print name
				}
			}
		' "$header"
	fi
done | sort -u
