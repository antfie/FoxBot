#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color

if [ $# -eq 0 ]
  then
    echo -e "${RED}Error: No version number specified. Try something like \"release.sh x.y\".${NC}"
    exit
fi

export VERSION=$1


echo -e "${CYAN}Releasing v${VERSION}...${NC}"


echo
./scripts/build.sh


echo -e "\n${CYAN}Creating Docker image v${VERSION}...${NC}"
docker pull alpine
docker build -t antfie/foxbot . --no-cache


echo -e "\n${CYAN}Publishing Docker image v${VERSION}...${NC}"
docker push antfie/foxbot


echo -e "\n${GREEN}Release Success${NC}"