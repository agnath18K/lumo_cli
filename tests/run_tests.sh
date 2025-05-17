#!/bin/bash

# Run all tests in the tests directory
# Usage: ./run_tests.sh [test_name]

# Set colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Print header
echo -e "${YELLOW}=======================================${NC}"
echo -e "${YELLOW}       Lumo CLI Test Runner           ${NC}"
echo -e "${YELLOW}=======================================${NC}"

# Change to the tests directory
cd "$(dirname "$0")"

# Check if a specific test was requested
if [ $# -eq 1 ]; then
    TEST_FILE="$1"
    if [[ ! "$TEST_FILE" == *_test.go ]]; then
        TEST_FILE="${TEST_FILE}_test.go"
    fi
    
    if [ -f "$TEST_FILE" ]; then
        echo -e "${YELLOW}Running test: ${TEST_FILE}${NC}"
        go test -v ./$TEST_FILE
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}Test passed!${NC}"
        else
            echo -e "${RED}Test failed!${NC}"
            exit 1
        fi
    else
        echo -e "${RED}Test file not found: ${TEST_FILE}${NC}"
        exit 1
    fi
else
    # Run all tests
    echo -e "${YELLOW}Running all tests...${NC}"
    go test -v ./...
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
    else
        echo -e "${RED}Some tests failed!${NC}"
        exit 1
    fi
fi

echo -e "${YELLOW}=======================================${NC}"
