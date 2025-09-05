#!/bin/bash

# SeedFlow Knowledge Management Tool - Backup Script
# This script creates backups of data and configuration

set -e

BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="seedflow_backup_${DATE}.tar.gz"

echo "💾 SeedFlow Knowledge Management Tool - Backup"
echo "============================================="

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

echo "📦 Creating backup: $BACKUP_FILE"

# Create backup excluding logs (they can be regenerated)
tar -czf $BACKUP_DIR/$BACKUP_FILE \
  --exclude='logs/*' \
  --exclude='backups/*' \
  data/ \
  config/ \
  exports/ \
  .env 2>/dev/null || true

if [ -f $BACKUP_DIR/$BACKUP_FILE ]; then
    echo "✅ Backup created successfully: $BACKUP_DIR/$BACKUP_FILE"
    
    # Show backup size
    BACKUP_SIZE=$(du -h $BACKUP_DIR/$BACKUP_FILE | cut -f1)
    echo "📏 Backup size: $BACKUP_SIZE"
    
    # Clean up old backups (keep only last 30)
    echo "🧹 Cleaning up old backups..."
    cd $BACKUP_DIR
    ls -t seedflow_backup_*.tar.gz | tail -n +31 | xargs -r rm --
    cd ..
    
    BACKUP_COUNT=$(ls -1 $BACKUP_DIR/seedflow_backup_*.tar.gz 2>/dev/null | wc -l)
    echo "📊 Total backups: $BACKUP_COUNT"
else
    echo "❌ Backup failed!"
    exit 1
fi

echo ""
echo "📋 Restore instructions:"
echo "   1. Stop SeedFlow: ./scripts/stop.sh"
echo "   2. Extract backup: tar -xzf $BACKUP_DIR/$BACKUP_FILE"
echo "   3. Start SeedFlow: ./scripts/start.sh"