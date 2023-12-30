#!/bin/bash

# ANSI color codes for colored output
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RESET='\033[0m'

# Build all packages in cmd/ directory
CNT=0
for cmd_pkg in ./cmd/*; do
  if [ -d "$cmd_pkg" ]; then

    # If the modify date of bin/cmd_pkg is older than the modify date of cmd_pkg,
    # then rebuild cmd_pkg.

    # Get the last modified date of bin/cmd_pkg
    bin_last_modified=$(stat -c %Y "./bin/$(basename "$cmd_pkg")" 2>/dev/null)
    cmd_pkg_last_modified=$(stat -c %Y "$cmd_pkg")

    # Loop over each file in the pkg folder
    for file in "$cmd_pkg"/*; do
      if [ -f "$file" ]; then
        # Get the last modified date of each file in pkg
        file_last_modified=$(stat -c %Y "$file")

        # Check if the file is newer than the binary
        if [ "$file_last_modified" -gt "$bin_last_modified" ] || [ "$cmd_pkg_last_modified" -gt "$bin_last_modified" ]; then
          CNT=$((CNT + 1))
          echo -e "${COLOR_YELLOW}[INFO]${COLOR_RESET} - Started building \"$(basename "$cmd_pkg")\"..."
          go build -o "./bin/$(basename "$cmd_pkg")" "$cmd_pkg"
          if [ $? -eq 0 ]; then
            echo -e "${COLOR_GREEN}[SUCCESS]${COLOR_RESET} - Successfully built \"$(basename "$cmd_pkg")\"."
          else
            echo -e "${COLOR_RED}[ERROR]${COLOR_RESET} - Failed to build \"$(basename "$cmd_pkg")\"."
            exit 1
          fi
          break  # No need to check other files if we've rebuilt the binary
        fi
      fi
    done
  fi
done

if [ "$CNT" -eq "0" ]; then
  echo -e "${COLOR_YELLOW}[INFO]${COLOR_RESET} - No packages were built."
fi

# Check if there are any files in the bin directory
if [ -n "$(find ./bin -maxdepth 1 -type f)" ]; then
  sudo cp ./bin/* /usr/local/bin/
else
  echo -e "${COLOR_RED}[ERROR]${COLOR_RESET} - No files found in the bin directory."
  exit 1
fi
