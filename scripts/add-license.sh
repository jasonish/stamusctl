#!/bin/bash

# Define the header content
HEADER="/*
    Copyright(C) 2024 Stamus Networks
    Written by Valentin Vivier <vvivier@stamus-networks.com>
    Written by Baptiste Ternoir <bternoir@stamus-networks.com>

    This file is part of Stamus-ctl.

    Stamus-ctl is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Stamus-ctl is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with Stamus-ctl.  If not, see <http://www.gnu.org/licenses/>.
 */
 
"

# Define the directory containing .go files
DIRECTORY="."

# Loop through all .go files in the specified directory
for FILE in "$DIRECTORY"/*.go; do
  if [ -f "$FILE" ]; then
    # Create a temporary file
    TEMP_FILE=$(mktemp)

    # Write the header and the original content to the temporary file
    echo "$HEADER" > "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    cat "$FILE" >> "$TEMP_FILE"

    # Move the temporary file to the original file
    mv "$TEMP_FILE" "$FILE"

    echo "Header added to $FILE"
  fi
done
