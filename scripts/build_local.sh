#!/usr/bin/env bash

# ensure we are in the repo root
DIR=$(realpath "$(dirname "$0")")
cd "$DIR/.." || exit

if [[ -z "${WSL_DISTRO_NAME}" ]] || [[ ! $DIR =~ ^/mnt/.* ]] ; then
  goreleaser --snapshot --skip-publish --rm-dist
else
  # on WSL we have to move to a non-windows file system for file permissions to be set correctly
  tmp_dir=$(mktemp -d -t nfetch-XXXXXXXXXX)
  echo "building in '$tmp_dir'"

  echo "cleaning old assets"
  rm -r ./dist

  echo "copying to tempdir"
  cp -R . "$tmp_dir"

  cd "$tmp_dir" || exit

  echo "setting correct permissions"
  chmod -R 644 ./completions/*

  echo "building assets"
  goreleaser --snapshot --skip-publish --rm-dist

  echo "copying resulting assets"
  cp -R ./dist "$DIR/.."
  cp -R ./completions/nfetch.sh "$DIR/../completions/nfetch.sh"

  echo "cleaning up"
  cd "$DIR/.." || exit
  rm -rf "$tmp_dir"
fi
