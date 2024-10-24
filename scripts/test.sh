#!/usr/bin/env bash

# Exit if any command fails
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m' # No Color


echo -e "${CYAN}Linting...${NC}"
gofmt -s -w .


echo -e "\n${CYAN}Running go vet...${NC}"
go vet ./...


echo -e "\n${CYAN}Running gosec...${NC}"
gosec ./...


echo -e "${CYAN}Running osv-scanner...${NC}"
osv-scanner scan .


echo -e "\n${CYAN}Testing...${NC}"
mkdir -p test
go test -v -coverprofile test/coverage.out -failfast -shuffle on -parallel 4 ./...
go tool cover -html=test/coverage.out -o test/coverage.html

# Consider using gcov2lcov tool / perhaps instead of go tool cover?
#gocov convert coverage.out > coverage.lcov
#genhtml coverage/lcov.info -o coverage/html -t 'StickerDocs API' --no-function-coverage

echo -e "\n${GREEN}All tests OK${NC}"