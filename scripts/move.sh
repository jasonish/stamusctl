#!/bin/bash

# Define source and destination directories
SOURCE_DIR="/.test/config"
DEST_DIR="/internal/embeds/selks"

# Check if the source directory exists
if [ -d "$SOURCE_DIR" ]; then
    echo "Source directory exists: $SOURCE_DIR"
else
    echo "Source directory does not exist: $SOURCE_DIR"
    exit 1
fi

# Create the destination directory if it does not exist
if [ ! -d "$DEST_DIR" ]; then
    echo "Destination directory does not exist. Creating: $DEST_DIR"
    mkdir -p "$DEST_DIR"
fi

# Move the folder from source to destination
mv "$SOURCE_DIR" "$DEST_DIR"

# Check if the move was successful
if [ $? -eq 0 ]; then
    echo "Successfully moved $SOURCE_DIR to $DEST_DIR"
else
    echo "Failed to move $SOURCE_DIR to $DEST_DIR"
    exit 1
fi
