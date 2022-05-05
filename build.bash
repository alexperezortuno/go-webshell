#!/usr/bin/env bash

package=""
APP_NAME="go_app"
ARCH=""
PLATFORM=""
ALL=1

RCol='\e[0m'
Yel='\e[033m'
Red='\e[031m'
Gre='\e[032m'
Divider='=============================================================='

usage() {
  echo -e "${Gre}"
  echo -e "API LOGS${RCol}\r"
  echo -e "${Divider}${RCol}\r"
  echo -e "${Yel}"
  echo -e "COMMAND\t\t\t\tDESCRIPTION\n"
  echo -e "${RCol}"
}

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64")

while getopts p:o:h:a: flag
do
    case "${flag}" in
        p) package=${OPTARG};;
        o) APP_NAME=${OPTARG};;
        a) platforms=(${OPTARG});;
        h) HELP=1
            usage
            exit1;;
    esac
done

if [[ -z "$package" ]]; then
  echo "usage: $package <package-name>"
  exit 1
fi

package_split=(${package//\// })
package_name=${package_split[-1]}

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$APP_NAME'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
