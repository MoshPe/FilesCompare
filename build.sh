#!/usr/bin/env bash
package_name="FileCompare"

build_folder=""

# Check operating system


platforms=("windows/amd64" "windows/386" "linux/386" "linux/amd64" "linux/arm" "linux/amd64")

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}

	if [ "$GOOS" == "windows" ]; then
      build_folder="windows-build"
  else
      build_folder="linux-build"
  fi

  # Create build folder if it doesn't exist
  if [ ! -d "$build_folder" ]; then
      mkdir -p "$build_folder"
  fi

	output_name=$package_name'-'$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

	env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
	mv $output_name $build_folder/$output_name

	# shellcheck disable=SC2181
	if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done

# Zip the build folders
zip_name="FileCompare.zip"
zip -r "$zip_name" "windows-build" "linux-build"
echo "Builds have been zipped to $zip_name"
