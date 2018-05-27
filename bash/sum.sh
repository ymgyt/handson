#!/bin/bash

readonly script_name=${0##*/}

while IFS= read -r n; do
  if [[ ! $n =~ ^-?[0-9]+$ ]]; then
    printf '%s\n' "${script_name}: '$n': non-integernumber" 1>&2
    exit 1
  fi

  ((result+=n))
done

printf '%s\n' "$result"
