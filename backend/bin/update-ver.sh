#!/bin/sh

# update version type (major, minor, patch, prerelease)
update_version_type="prerelease"
if [ -n "$1" ]; then
	update_version_type="$1"
fi

# current version
current_version=$(grep -oEi 'version: v([0-9\.?]+)' conf/config.yml | sed -E 's/version: v//g')

# next version
next_version=$(./bin/semver.sh bump $update_version_type $current_version)

# update next version to conf/config.yml
sed -i '' "s/version: v$current_version/version: v$next_version/g" conf/config.yml

