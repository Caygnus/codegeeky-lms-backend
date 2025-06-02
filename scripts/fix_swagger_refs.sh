#!/bin/bash

# Script to replace allOf with direct $ref in swagger.json
# This improves type definitions in generated documentation

# Configuration
SWAGGER_FILE="docs/swagger/swagger.json"
BACKUP_FILE="${SWAGGER_FILE}.bak"

# Check if the swagger file exists
if [ ! -f "$SWAGGER_FILE" ]; then
    echo "Error: Swagger file not found at $SWAGGER_FILE"
    exit 1
fi

# Create a backup of the original file
cp "$SWAGGER_FILE" "$BACKUP_FILE"
echo "Created backup at $BACKUP_FILE"

# Create a temporary Python script to handle the JSON processing
cat > /tmp/fix_swagger_refs.py << 'EOF'
import json
import re
import sys

swagger_file = sys.argv[1]

# Read the JSON file
with open(swagger_file, 'r') as f:
    content = f.read()

# Count original references
original_ref_count = content.count('"$ref"')

# Define the pattern to match - "allOf": [ { "$ref": "#/definitions/Type" } ]
pattern = r'"allOf":\s*\[\s*{\s*"\$ref":\s*"(#/definitions/[^"]+)"\s*}\s*\]'

# Function to replace matched pattern with direct reference
def replace_allof(match):
    ref_path = match.group(1)
    return '"$ref": "' + ref_path + '"'

# Perform the replacement
modified_content = re.sub(pattern, replace_allof, content)

# Count new references
new_ref_count = modified_content.count('"$ref"')
replaced_count = new_ref_count - original_ref_count

# Write the modified content back to the file
with open(swagger_file, 'w') as f:
    f.write(modified_content)

print(f'Replaced {replaced_count} allOf patterns with direct $ref references')
print(f'Total $ref count: {new_ref_count}')
EOF

# Run the Python script
python3 /tmp/fix_swagger_refs.py "$SWAGGER_FILE" || {
    echo "Error: Failed to process the file with Python"
    # Restore the backup
    cp "$BACKUP_FILE" "$SWAGGER_FILE"
    exit 1
}

# Clean up
rm /tmp/fix_swagger_refs.py
rm "$BACKUP_FILE"

echo "Processed $SWAGGER_FILE"
echo "Backup file stored at $BACKUP_FILE"
echo "Done! The swagger.json file has been updated."