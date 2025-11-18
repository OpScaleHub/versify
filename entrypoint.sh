#!/bin/sh -l

# This script is the entrypoint for the Docker container.
# It reads the inputs from the environment variables and executes the versify binary.

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Helper Functions ---
log() {
  echo "INFO: $1"
}

# --- Main Execution ---
log "Starting Versify Action..."

# Build the command-line arguments for the versify binary.
# The inputs from action.yml are passed as environment variables with an 'INPUT_' prefix.
ARGS=""
if [ -n "$INPUT_PREFIX" ]; then
  ARGS="$ARGS --prefix=$INPUT_PREFIX"
fi

if [ -n "$INPUT_BASELINE" ]; then
  ARGS="$ARGS --baseline=$INPUT_BASELINE"
fi

if [ "$INPUT_ADD_SUFFIX" = "true" ]; then
  ARGS="$ARGS --add-suffix"
fi

if [ -n "$INPUT_SUFFIX_FORMAT" ]; then
  ARGS="$ARGS --suffix-format=$INPUT_SUFFIX_FORMAT"
fi

log "Running versify with arguments: $ARGS"

# Execute the versify binary and capture the output.
# The stderr is redirected to the current process's stderr.
NEW_VERSION=$(/versify $ARGS)

# Determine the last version to compare against.
# This logic helps determine if a version bump actually occurred.
LAST_VERSION=""
if [ -n "$INPUT_BASELINE" ]; then
    LAST_VERSION="$INPUT_BASELINE"
else
    # Attempt to get the latest tag. Ignore errors if no tags exist.
    LAST_TAG=$(git describe --tags --abbrev=0 --match "$INPUT_PREFIX[0-9]*.[0-9]*.[0-9]*" 2>/dev/null || echo "")
    if [ -n "$LAST_TAG" ]; then
        LAST_VERSION=$(echo "$LAST_TAG" | sed "s/^$INPUT_PREFIX//")
    fi
fi

# Determine if a version bump happened.
SHOULD_BUMP="false"
# A bump occurred if the new version is different from the last known version.
# Also handles the case where there was no last version.
if [ "$NEW_VERSION" != "$LAST_VERSION" ]; then
  SHOULD_BUMP="true"
fi

# Set the action outputs.
# These are written to the file specified by the GITHUB_OUTPUT environment variable.
log "Setting outputs..."
log "  new_version: $NEW_VERSION"
log "  should_bump: $SHOULD_BUMP"

echo "new_version=$NEW_VERSION" >> "$GITHUB_OUTPUT"
echo "should_bump=$SHOULD_BUMP" >> "$GITHUB_OUTPUT"

log "Versify Action finished."
