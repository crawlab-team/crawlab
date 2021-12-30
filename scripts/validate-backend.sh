#!/bin/sh

IFS=$'\n'
pattern=^replace
content=$(cat ./backend/go.mod)
for line in $content
do
	if [[ $line =~ $pattern ]]; then
		exit 1
	fi
done
