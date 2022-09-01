#!/usr/bin/env bash

current_dir="$(dirname -- "$0")"

commit_id="$(git rev-parse HEAD)"
tag_name="$(git describe --contains $commit_id 2> /dev/null)"
branch_name="$(git rev-parse --abbrev-ref HEAD)"

if [[ ! -z "$tag_name" ]]; then
  echo "$tag_name"
  exit 0
fi

last_version="$(git tag --sort=-v:refname | head -1)"

if [[ -z "$last_version" ]]; then
  last_version="v0.0.0"
fi

version="$("${current_dir}/upgrade_patch.sh" $last_version)"

if [[ "$branch_name" != "master" ]]; then
  version="${version}-${branch_name}"
fi

version="${version}-SNAPSHOT"

echo "$version"
