#!/usr/bin/bash

version=$1
if [[ -z "$version" ]]; then
  echo "usage: $0 <version>"
  exit 1
fi

package="github.com/spitzfaust/file-server"
output_folder="build"
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64" "linux/386")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$output_folder'/'$package_name'-'$GOOS'-'$GOARCH'-'$version
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
