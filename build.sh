#!/bin/bash

# ANSI color codes for colored output
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RESET='\033[0m'

# Build all packages in cmd/ directory
for pkg in ./cmd/*; do
  if [ -d "$pkg" ]; then

    # If the modify date of bin/pkg is older than the modify date of pkg, then
    # rebuild pkg.
    if [ ! -f "./bin/$(basename "$pkg")" ] || [ "./bin/$(basename "$pkg")" -ot "$pkg" ]; then
      echo -e "${COLOR_YELLOW}[INFO]${COLOR_RESET} - Started building \"$(basename "$pkg")\"..."
      go build -o "./bin/$(basename "$pkg")" "$pkg"
      if [ $? -eq 0 ]; then
      echo -e "${COLOR_GREEN}[SUCCESS]${COLOR_RESET} - Successfully built \"$(basename "$pkg")\"."
      else
      echo -e "${COLOR_RED}[ERROR]${COLOR_RESET} - Failed to build \"$(basename "$pkg")\"."
      exit 1
      fi
    fi
  fi
done

sudo cp ./bin/* /usr/local/bin/

