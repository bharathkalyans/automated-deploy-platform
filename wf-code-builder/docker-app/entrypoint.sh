#!/bin/bash

_error() {
    # Error Format to make parsing easier
    # ERROR::ERROR_CODE::REASON
    echo "ERROR::$1::$2"
    exit 1
}
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--project-github-url)
            PROJECT_GITHUB_URL="$2"
            shift 
            shift
            ;;
        -b|--build-command)
            BUILD_COMMAND="$2"
            shift
            shift
            ;;
        -o|--build-out-dir)
            BUILD_OUT_DIR="$2"
            shift
            shift
            ;;
        *)
            # Need not to give the exact error that occurs here to user, but for the
            # sake of this assignment we can do it. 
            _error 'DATA_PROCESSING_ERROR' $1
            ;;
    esac
done
# Pre processing
# Validate package manager
# First value after splitting the build command gives us the package manager
read -a BUILD_COMMAND_ARRAY <<< "$BUILD_COMMAND"
PKG_MGR="${BUILD_COMMAND_ARRAY[0]}"

if ! [ "$PKG_MGR" = "npm" ] && ! [ "$PKG_MGR" = "yarn" ]; then
    _error 'UNSUPPORTED_PKG_MGR' "This (${PKG_MGR}) package manager is not supported"
fi

# Clone the repository
git clone "$PROJECT_GITHUB_URL" /app || _error 'REPO_INACCESSIBLE' "Failed to clone the repository at ${PROJECT_GITHUB_URL}"

# Install dependencies based on the choice of package manager
$PKG_MGR install || _error 'PKG_INSTALL_FAILED' "Failed to install dependencies"

# Build the app
$BUILD_COMMAND || _error 'BUILD_COMMAND_FAILED' "The build command (${BUILD_COMMAND}) failed to execute"

# Check if out dir mentioned is valid

OUT="/app/${BUILD_OUT_DIR}"
if [ -d "$OUT" ]; then 
    # Check if the directory contains an index.html file
    if ! [ -f "$OUT/index.html" ]; then
        _error 'INVALID_OUT_DIR' 'Output dir does not contain an index.html file'
    fi
else
    _error 'INVALID_OUT_DIR' 'Specified out dir was not generated'
fi

# Copy the out dir to wf_storage volume
cp -r $OUT /wf/storage/$BUILD_ID/