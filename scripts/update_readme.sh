#!/bin/bash
# ==========================================
# update_readme.sh - Update README.md with current example code
# This script extracts code from example files and injects it into README.md
# ==========================================
set -e  # Exit on error

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Banner
echo -e "${GREEN}==================================================${NC}"
echo -e "${GREEN}   debuggo README Updater${NC}"
echo -e "${GREEN}==================================================${NC}"

# Paths
README="README.md"
TEMPLATE="scripts/README_template.md"
README_BACKUP="README.md.bak"
BASIC_EXAMPLE="examples/basic/main.go"
ADVANCED_EXAMPLE="examples/advanced/main.go"

# Check if files exist
if [ ! -f "$BASIC_EXAMPLE" ] || [ ! -f "$ADVANCED_EXAMPLE" ]; then
  echo -e "${YELLOW}Error: Example files not found${NC}"
  exit 1
fi

# Create README template if it doesn't exist
if [ ! -f "$TEMPLATE" ]; then
  echo -e "${YELLOW}Creating README template from current README...${NC}"
  cp "$README" "$TEMPLATE"
  
  # Mark sections for automatic replacement
  sed -i.bak -e '/```go.*Basic example of using debuggo/,/```/ c\```go\n<!-- BASIC_EXAMPLE -->\n```' "$TEMPLATE"
  sed -i.bak -e '/```go.*Advanced example of using debuggo/,/```/ c\```go\n<!-- ADVANCED_EXAMPLE -->\n```' "$TEMPLATE"
  
  echo -e "${YELLOW}Template created. Please edit ${TEMPLATE} to mark sections for replacement.${NC}"
  echo -e "${YELLOW}Use <!-- BASIC_EXAMPLE --> and <!-- ADVANCED_EXAMPLE --> as placeholders.${NC}"
  rm -f "${TEMPLATE}.bak"
else
  echo -e "${YELLOW}Using existing README template${NC}"
fi

# Back up current README
cp "$README" "$README_BACKUP"
echo -e "${YELLOW}Backed up README to ${README_BACKUP}${NC}"

# Extract example code
BASIC_CODE=$(cat "$BASIC_EXAMPLE")
ADVANCED_CODE=$(cat "$ADVANCED_EXAMPLE")

# Generate new README
echo -e "${YELLOW}Generating new README...${NC}"
cp "$TEMPLATE" "$README"

# Replace placeholders with actual code
sed -i.bak "s|<!-- BASIC_EXAMPLE -->|$BASIC_CODE|" "$README"
sed -i.bak "s|<!-- ADVANCED_EXAMPLE -->|$ADVANCED_CODE|" "$README"

# Clean up
rm -f "${README}.bak"

echo -e "${GREEN}README.md updated successfully!${NC}"
echo -e "${YELLOW}Please verify the changes to make sure the formatting is correct.${NC}" 