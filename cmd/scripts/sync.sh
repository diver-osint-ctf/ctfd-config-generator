#!/bin/bash

genres=("geo" "sns" "crypto" "transportation" "darkweb" "history" "company" "misc" "hardware" "military")
for genre in "${genres[@]}"; do
    if [ ! -d "../$genre" ]; then
        continue
    fi
    sorted_genres=$(find "../$genre" | grep "challenge.yaml$" | sort)

    for g in $sorted_genres; do
        python -m ctfcli challenge sync $g;
        sleep 1; # delay for server
    done
done