#!/usr/bin/env bash

current_version=$1
[[ "$current_version" =~ ^v(.*)\.(.*)\.(.*).*$ ]] && \
  major="${BASH_REMATCH[1]}" && \
  minor="${BASH_REMATCH[2]}" && \
  patch="${BASH_REMATCH[3]}"

if [[ -z "$major" || -z "$minor" || -z "$patch" ]]; then
  echo "invalid version" 1>&2
  exit 1
fi

echo "v${major}.${minor}.$(($patch + 1))"
