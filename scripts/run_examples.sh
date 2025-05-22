#!/bin/bash
# ==========================================
# run_examples.sh - Run debuggo examples and capture their output
# This script executes the examples with different DEBUG settings
# and captures their output for documentation purposes.
# ==========================================
set -e  # Exit on error

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Banner
echo -e "${GREEN}==================================================${NC}"
echo -e "${GREEN}   debuggo Examples Runner${NC}"
echo -e "${GREEN}==================================================${NC}"

# Function to handle errors gracefully
handle_error() {
	echo -e "${YELLOW}Warning: An error occurred but continuing...${NC}"
}

# Ensure script continues even if some commands fail
trap handle_error ERR

# Get the absolute path to the project root directory
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
echo -e "${YELLOW}Project root: ${ROOT_DIR}${NC}"

# Clean up any previous output files
echo -e "${YELLOW}Cleaning up previous output files...${NC}"
rm -f "${ROOT_DIR}"/example_outputs/*.stdout "${ROOT_DIR}"/example_outputs/*.stderr 2>/dev/null || true

# Create temporary directory
TMPDIR=$(mktemp -d)
echo -e "${YELLOW}Using temporary directory: ${TMPDIR}${NC}"

# Clean up on exit
trap "echo 'Cleaning up temporary files...'; rm -rf $TMPDIR" EXIT

# Create directories for running examples
mkdir -p "${TMPDIR}/basic"
mkdir -p "${TMPDIR}/advanced"

# Copy the examples to the temporary directory
cp "${ROOT_DIR}/examples/basic/main.go" "${TMPDIR}/basic/"
cp "${ROOT_DIR}/examples/advanced/main.go" "${TMPDIR}/advanced/"

# Create go.mod files for both examples with the local replace directive
cat > "${TMPDIR}/basic/go.mod" << EOF
module basicexample

go 1.19

replace github.com/GeoffreyPlitt/debuggo => ${ROOT_DIR}
require github.com/GeoffreyPlitt/debuggo v0.0.0-00000000000000-000000000000
EOF

cat > "${TMPDIR}/advanced/go.mod" << EOF
module advancedexample

go 1.19

replace github.com/GeoffreyPlitt/debuggo => ${ROOT_DIR}
require github.com/GeoffreyPlitt/debuggo v0.0.0-00000000000000-000000000000
EOF

# Output directory - ensure it exists with absolute path
OUTDIR="${ROOT_DIR}/example_outputs"
mkdir -p "${OUTDIR}"

echo -e "${YELLOW}Output will be saved to ${OUTDIR}${NC}"

# Run examples and capture output
echo -e "${YELLOW}=== Running basic example with DEBUG=* ===${NC}"
(cd "${TMPDIR}/basic" && go mod tidy && DEBUG="*" go run . > "${OUTDIR}/basic_all_output.stdout" 2> "${OUTDIR}/basic_all_output.stderr")

echo ""
echo -e "${YELLOW}=== Running basic example with DEBUG=db ===${NC}"
(cd "${TMPDIR}/basic" && DEBUG="db" go run . > "${OUTDIR}/basic_db_output.stdout" 2> "${OUTDIR}/basic_db_output.stderr")

echo ""
echo -e "${YELLOW}=== Running advanced example with DEBUG=app:server:* ===${NC}"
(cd "${TMPDIR}/advanced" && go mod tidy && DEBUG="app:server:*" go run . > "${OUTDIR}/advanced_server_output.stdout" 2> "${OUTDIR}/advanced_server_output.stderr")

echo ""
echo -e "${YELLOW}=== Running advanced example with DEBUG=* ===${NC}"
(cd "${TMPDIR}/advanced" && DEBUG="*" go run . > "${OUTDIR}/advanced_all_output.stdout" 2> "${OUTDIR}/advanced_all_output.stderr")

echo ""
echo -e "${GREEN}All outputs captured in ${OUTDIR}/${NC}"
echo -e "${GREEN}You can use these files to update the README.md or documentation.${NC}"

# Print a summary of what was captured
echo ""
echo -e "${GREEN}=== Basic Example (DEBUG=*) ===${NC}"
echo -e "${YELLOW}STDOUT:${NC}"
cat "${OUTDIR}/basic_all_output.stdout" 2>/dev/null || echo "No stdout captured"
echo -e "${YELLOW}STDERR:${NC}"
cat "${OUTDIR}/basic_all_output.stderr" 2>/dev/null || echo "No stderr captured"

echo ""
echo -e "${GREEN}=== Basic Example (DEBUG=db) ===${NC}"
echo -e "${YELLOW}STDOUT:${NC}"
cat "${OUTDIR}/basic_db_output.stdout" 2>/dev/null || echo "No stdout captured"
echo -e "${YELLOW}STDERR:${NC}"
cat "${OUTDIR}/basic_db_output.stderr" 2>/dev/null || echo "No stderr captured"

echo ""
echo -e "${GREEN}=== Advanced Example (DEBUG=app:server:*) ===${NC}"
echo -e "${YELLOW}STDOUT:${NC}"
cat "${OUTDIR}/advanced_server_output.stdout" 2>/dev/null || echo "No stdout captured"
echo -e "${YELLOW}STDERR:${NC}"
cat "${OUTDIR}/advanced_server_output.stderr" 2>/dev/null || echo "No stderr captured"

echo ""
echo -e "${GREEN}=== Advanced Example (DEBUG=*) ===${NC}"
echo -e "${YELLOW}STDOUT:${NC}"
cat "${OUTDIR}/advanced_all_output.stdout" 2>/dev/null || echo "No stdout captured"
echo -e "${YELLOW}STDERR:${NC}"
cat "${OUTDIR}/advanced_all_output.stderr" 2>/dev/null || echo "No stderr captured"

echo ""
echo -e "${GREEN}==================================================${NC}"
echo -e "${GREEN}   Script completed successfully${NC}"
echo -e "${GREEN}==================================================${NC}"

# Instructions for next steps
echo ""
echo -e "${YELLOW}To use these outputs in the README.md:${NC}"
echo "1. Copy the content from the output files"
echo "2. Update the README.md with the actual output examples"
echo ""
echo -e "${YELLOW}To clean up all generated files:${NC}"
echo "  rm -rf ${OUTDIR}"

# Always exit with success
exit 0 