#!/bin/bash
setsid sh -c './sleeploudly.sh ${1} & ./sleeploudly.sh ${1}& ./sleeploudly.sh ${1}&'
