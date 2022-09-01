#!/usr/bin/env bash

current_dir="$(dirname -- "$0")"

commit_id="$(git rev-parse HEAD)"
tag_name="$(git describe --contains $commit_id 2> /dev/null)"
branch_name="$(git rev-parse --abbrev-ref HEAD)"

version="v0.0.1-SNAPSHOT"
is_snapshot=0

if [[ "$tag_name" == *"~"* ]]; then
  version="${tag_name%~*}"
  is_snapshot=1
else
  if [[ ! -z "$tag_name" ]]; then
    version="$tag_name"
  fi
fi

if [[ "$branch_name" != "master" ]]; then
  version="${version}-${branch_name}"
fi

if [[ $is_snapshot == 1 ]]; then
  version="$($current_dir/upgrade_patch.sh $version)-SNAPSHOT"
fi

echo "$version"
