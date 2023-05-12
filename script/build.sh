#!/usr/bin/env bash

set -e

package_name='protocoll'
platforms=("windows/amd64" "darwin/amd64" "linux/amd64")
out='./out/'

if [[ ! -f $out ]]; then
  mkdir -p $out
fi

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name=$package_name'-'$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

	env GOOS=$GOOS GOARCH=$GOARCH go build -o $out$output_name ./cmd
	if [ $? -ne 0 ]; then
   		echo "Could not build for ${platform}"
		exit 1
	fi
done
