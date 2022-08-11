#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
BIN=$SCRIPT_DIR/$(basename "$SCRIPT_DIR")
LOG_PATH=$SCRIPT_DIR/../../logs/$(basename "$SCRIPT_DIR").log

if [ -n "$1" ]; then
	$BIN $1 $2 >> $LOG_PATH 2>&1
else
	$BIN >> $LOG_PATH 2>&1
fi

