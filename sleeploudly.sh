#!/bin/bash
trap '' 2
echo "$$ ${BASHPID} sleeping loudly for ${1} seconds" >> sleep.log
sleep $1
echo "$$ ${BASHPID} finished sleeping loudly" >> sleep.log
