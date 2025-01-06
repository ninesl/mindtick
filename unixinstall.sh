#!/bin/bash

BINARY_NAME="mindtick"
DEST_DIR="/usr/local/bin"

# Ensure the user has write permissions
if [ ! -w $DEST_DIR ]; then
    echo "Error: You need sudo permissions to copy to $DEST_DIR."
    exit 1
fi

# Copy the binary to the destination directory
rm -f $DEST_DIR/$BINARY_NAME
cp $BINARY_NAME $DEST_DIR/

# Make it executable
chmod +x $DEST_DIR/$BINARY_NAME

echo "$BINARY_NAME installed to $DEST_DIR. You can now use it anywhere."