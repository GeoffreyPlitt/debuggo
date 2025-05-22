#!/bin/bash
# run_examples.sh - Run debuggo examples and capture their output
# This script executes the examples with different DEBUG settings
# and captures their output for documentation purposes.
set -e  # Exit on error

# Function to handle errors gracefully
handle_error() {
	echo "Warning: An error occurred but continuing..."
}

# Ensure script continues even if some commands fail
trap handle_error ERR

# Get the absolute path to the project root directory
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
echo "Project root: ${ROOT_DIR}"

# Create temporary directory
TMPDIR=$(mktemp -d)
echo "Using temporary directory: ${TMPDIR}"

# Clean up on exit
trap "echo 'Cleaning up temporary files...'; rm -rf $TMPDIR" EXIT

# Create directories for running examples
mkdir -p "${TMPDIR}/basic"
mkdir -p "${TMPDIR}/advanced"

# Copy the examples to the temporary directory
cp "${ROOT_DIR}/examples/basic/main.go" "${TMPDIR}/basic/"
cp "${ROOT_DIR}/examples/advanced/main.go" "${TMPDIR}/advanced/"

# Create go.mod file for an example
create_gomod() {
    local dir=$1
    local module=$2
    
    cat > "${dir}/go.mod" << EOF
module ${module}

go 1.19

replace github.com/GeoffreyPlitt/debuggo => ${ROOT_DIR}
require github.com/GeoffreyPlitt/debuggo v0.0.0-00000000000000-000000000000
EOF
}

# Create go.mod files for both examples with the local replace directive
create_gomod "${TMPDIR}/basic" "basicexample"
create_gomod "${TMPDIR}/advanced" "advancedexample"

# Output directory - ensure it exists with absolute path
OUTDIR="${ROOT_DIR}/example_outputs"
mkdir -p "${OUTDIR}"

echo "Output will be saved to ${OUTDIR}"

# Run an example and capture its output
run_example() {
    local example_type=$1    # basic or advanced
    local debug_value=$2     # DEBUG environment variable value
    local output_prefix=$3   # prefix for output files
    
    echo "Running ${example_type} example with DEBUG=${debug_value}"
    
    # Run tidy only on first run for each example type
    if [[ "$4" == "tidy" ]]; then
        (cd "${TMPDIR}/${example_type}" && go mod tidy && DEBUG="${debug_value}" go run . > "${OUTDIR}/${output_prefix}.stdout" 2> "${OUTDIR}/${output_prefix}.stderr")
    else
        (cd "${TMPDIR}/${example_type}" && DEBUG="${debug_value}" go run . > "${OUTDIR}/${output_prefix}.stdout" 2> "${OUTDIR}/${output_prefix}.stderr")
    fi
    
    echo "Generated ${output_prefix}.stdout and ${output_prefix}.stderr"
}

# Display output summary
display_output() {
    local description=$1
    local prefix=$2
    
    echo ""
    echo "${description}"
    echo "STDOUT:"
    cat "${OUTDIR}/${prefix}.stdout" 2>/dev/null || echo "No stdout captured"
    echo "STDERR:"
    cat "${OUTDIR}/${prefix}.stderr" 2>/dev/null || echo "No stderr captured"
}

# Run examples and capture output
run_example "basic" "*" "basic_all_output" "tidy"
run_example "basic" "db" "basic_db_output"
run_example "advanced" "app:server:*" "advanced_server_output" "tidy"
run_example "advanced" "*" "advanced_all_output"

echo "All outputs captured in ${OUTDIR}"
echo "You can use these files to update the README.md or documentation."

# Print a summary of what was captured
display_output "Basic Example (DEBUG=*)" "basic_all_output"
display_output "Basic Example (DEBUG=db)" "basic_db_output"
display_output "Advanced Example (DEBUG=app:server:*)" "advanced_server_output"
display_output "Advanced Example (DEBUG=*)" "advanced_all_output"

echo "Script completed successfully"

# Instructions for next steps
echo ""
echo "To use these outputs in the README.md:"
echo "1. Copy the content from the output files"
echo "2. Update the README.md with the actual output examples"
echo ""
echo "To clean up all generated files:"
echo "  rm -rf ${OUTDIR}"

# Always exit with success
exit 0 