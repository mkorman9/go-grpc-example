#!/usr/bin/env bash

branch_name="$(git rev-parse --abbrev-ref HEAD)"
last_version="$(git tag --sort=-v:refname | head -1)"

if [[ -z "$last_version" ]]; then
  last_version="v0.0.0"
fi

if [[ "$branch_name" != "master" ]]; then
  echo "release can only be done from the master branch" 1>&2
  exit 1
fi

[[ "$last_version" =~ ^v(.*)\.(.*)\.(.*).*$ ]] && \
  major="${BASH_REMATCH[1]}" && \
  minor="${BASH_REMATCH[2]}" && \
  patch="${BASH_REMATCH[3]}"

new_version="v${major}.$(($minor + 1)).0"

git tag "$new_version"
git push origin --tags
