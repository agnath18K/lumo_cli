#!/bin/bash
# Script to reconcile divergent Git branches
# Usage: ./scripts/reconcile_branches.sh <source_branch> <target_branch>

set -e

# Display help if no arguments provided
if [ $# -lt 2 ]; then
    echo "Usage: $0 <source_branch> <target_branch> [--force]"
    echo ""
    echo "This script helps reconcile divergent Git branches by:"
    echo "  1. Checking the status of both branches"
    echo "  2. Identifying divergent commits"
    echo "  3. Providing options to merge, rebase, or cherry-pick"
    echo ""
    echo "Arguments:"
    echo "  source_branch   The branch containing changes you want to incorporate"
    echo "  target_branch   The branch you want to update"
    echo "  --force         Skip safety checks (use with caution)"
    echo ""
    echo "Examples:"
    echo "  $0 feature/new-feature main"
    echo "  $0 main feature/outdated-branch"
    echo "  $0 release/v1.0 main --force"
    exit 1
fi

SOURCE_BRANCH=$1
TARGET_BRANCH=$2
FORCE_MODE=false

# Check for force flag
if [ $# -eq 3 ] && [ "$3" == "--force" ]; then
    FORCE_MODE=true
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if a branch exists
branch_exists() {
    git show-ref --verify --quiet refs/heads/$1
    return $?
}

# Function to check if there are uncommitted changes
has_uncommitted_changes() {
    if ! git diff-index --quiet HEAD --; then
        return 0 # Has changes
    else
        return 1 # No changes
    fi
}

# Function to count commits ahead/behind
count_commits_diff() {
    local source=$1
    local target=$2
    
    # Get ahead/behind counts
    local ahead=$(git rev-list --count $target..$source)
    local behind=$(git rev-list --count $source..$target)
    
    echo "$ahead $behind"
}

# Function to show commits that are in source but not in target
show_unique_commits() {
    local source=$1
    local target=$2
    
    echo -e "${BLUE}Commits in $source that are not in $target:${NC}"
    git log --oneline --graph --decorate $target..$source
}

# Function to perform the reconciliation
perform_reconciliation() {
    local method=$1
    local source=$2
    local target=$3
    
    # Save the current branch
    local current_branch=$(git symbolic-ref --short HEAD)
    
    case $method in
        merge)
            echo -e "${BLUE}Performing merge from $source into $target...${NC}"
            git checkout $target
            git merge $source
            ;;
        rebase)
            echo -e "${BLUE}Performing rebase of $target onto $source...${NC}"
            git checkout $target
            git rebase $source
            ;;
        cherry-pick)
            echo -e "${BLUE}Performing cherry-pick of commits from $source to $target...${NC}"
            git checkout $target
            # Get the list of commit hashes to cherry-pick
            local commits=$(git rev-list $target..$source --reverse)
            for commit in $commits; do
                echo -e "${BLUE}Cherry-picking commit ${commit:0:8}...${NC}"
                if ! git cherry-pick $commit; then
                    echo -e "${RED}Cherry-pick conflict detected.${NC}"
                    echo -e "${YELLOW}Resolve conflicts and then run:${NC}"
                    echo "git cherry-pick --continue"
                    echo -e "${YELLOW}Or to abort:${NC}"
                    echo "git cherry-pick --abort"
                    return 1
                fi
            done
            ;;
        *)
            echo -e "${RED}Invalid reconciliation method: $method${NC}"
            return 1
            ;;
    esac
    
    # Return to the original branch
    git checkout $current_branch
    return 0
}

# Main script execution starts here
echo -e "${BLUE}=== Git Branch Reconciliation Tool ===${NC}"
echo -e "${BLUE}Source branch: ${GREEN}$SOURCE_BRANCH${NC}"
echo -e "${BLUE}Target branch: ${GREEN}$TARGET_BRANCH${NC}"
echo ""

# Check if branches exist
if ! branch_exists $SOURCE_BRANCH; then
    echo -e "${RED}Error: Source branch '$SOURCE_BRANCH' does not exist.${NC}"
    exit 1
fi

if ! branch_exists $TARGET_BRANCH; then
    echo -e "${RED}Error: Target branch '$TARGET_BRANCH' does not exist.${NC}"
    exit 1
fi

# Check for uncommitted changes
if has_uncommitted_changes && [ "$FORCE_MODE" != "true" ]; then
    echo -e "${RED}Error: You have uncommitted changes in your working directory.${NC}"
    echo -e "${YELLOW}Please commit or stash your changes before reconciling branches.${NC}"
    echo -e "${YELLOW}Or use --force to proceed anyway (not recommended).${NC}"
    exit 1
fi

# Get the commit difference counts
read ahead behind <<< $(count_commits_diff $SOURCE_BRANCH $TARGET_BRANCH)

echo -e "${BLUE}Branch Status:${NC}"
echo -e "  ${SOURCE_BRANCH} is ${GREEN}$ahead commits ahead${NC} and ${YELLOW}$behind commits behind${NC} ${TARGET_BRANCH}"

# Show the unique commits
echo ""
show_unique_commits $SOURCE_BRANCH $TARGET_BRANCH
echo ""
show_unique_commits $TARGET_BRANCH $SOURCE_BRANCH
echo ""

# Provide reconciliation options
echo -e "${BLUE}Reconciliation Options:${NC}"
echo -e "  ${GREEN}1. Merge${NC} - Merge $SOURCE_BRANCH into $TARGET_BRANCH"
echo -e "  ${GREEN}2. Rebase${NC} - Rebase $TARGET_BRANCH onto $SOURCE_BRANCH"
echo -e "  ${GREEN}3. Cherry-pick${NC} - Cherry-pick commits from $SOURCE_BRANCH to $TARGET_BRANCH"
echo -e "  ${RED}4. Abort${NC} - Exit without making changes"
echo ""

# Get user choice
read -p "Select an option (1-4): " choice

case $choice in
    1)
        perform_reconciliation "merge" $SOURCE_BRANCH $TARGET_BRANCH
        ;;
    2)
        perform_reconciliation "rebase" $SOURCE_BRANCH $TARGET_BRANCH
        ;;
    3)
        perform_reconciliation "cherry-pick" $SOURCE_BRANCH $TARGET_BRANCH
        ;;
    4)
        echo -e "${YELLOW}Operation aborted. No changes were made.${NC}"
        exit 0
        ;;
    *)
        echo -e "${RED}Invalid option. Exiting.${NC}"
        exit 1
        ;;
esac

# Final status message
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Branch reconciliation completed successfully.${NC}"
else
    echo -e "${RED}Branch reconciliation encountered issues.${NC}"
    echo -e "${YELLOW}Please resolve any conflicts and complete the operation manually.${NC}"
fi
