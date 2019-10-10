#!/bin/bash

# variables
package_name="builder"
target="./src/builder/main.go"
output="./dist/"
conf="./src/builder/conf"

if [ -d "$output" ]; then
    rm -rf $output
fi
mkdir -p $output

# cross platform
platforms=(
"darwin/amd64"
"linux/amd64"
"linux/arm64"
"windows/amd64" )

# cross compile
for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo 'Compiling ---> '$GOOS'-'$GOARCH
    env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -a --ldflags=--s -o $output$output_name $target
    if [ $GOOS = "linux" -a $GOARCH = "amd64" ]; then
        cp -f $output$output_name $package_name
    fi

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

# default config
cp -rf $conf $output