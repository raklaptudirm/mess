#!/bin/zsh

cd "$(dirname "$0")"
go run . > dirty-data.tune
sort dirty-data.tune | uniq > "$1"
rm dirty-data.tune
wc -l "$1"