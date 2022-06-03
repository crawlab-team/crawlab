#!/bin/bash

function install_plugin() {
  # plugins executables directory
  local bin_path="/app/plugins/bin"
  if [ -d $bin_path ]; then
    :
  else
    mkdir -p "$bin_path"
  fi

  # plugin name
  local name=$1
  local url="https://github.com/crawlab-team/${name}"
  local repo_path=""/app/plugins/${name}
  git clone "$url" "$repo_path"
  cd "$repo_path" && go build -o "${bin_path}/${name}"
  chmod +x "${bin_path}/${name}"
}

plugin_names="plugin-dependency plugin-notification plugin-spider-assistant"

for name in $plugin_names; do
  install_plugin "$name"
done
