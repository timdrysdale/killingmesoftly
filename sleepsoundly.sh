#!/bin/bash
# https://www.cyberciti.biz/faq/unix-linux-shell-scripting-disable-controlc/
# Signal 2 is Ctrl+C
# Disable it:
trap 'sleep 10' 2  
sleep 2
echo "finished sleeping" > sleep.log
# Enable Ctrl+C
trap 2
