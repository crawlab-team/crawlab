#!/bin/sh

IFS=$'\n'
pattern=^replace
content=$(cat ./backend/go.mod)
for line in $content
do
	if [[ $line =~ $pattern ]]; then
		echo "Invalid ./backend/go.mod, which should not contain \"^replace\""
		exit 1
	fi
done
