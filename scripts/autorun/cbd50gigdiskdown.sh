#!/bin/bash

# check the status of the device and if it's there shut it down and wait for it to go away

logger systemd shutdown cbd50gigdisk service starting, waiting for cbd50gigdisk to free up

while true; do 
    lsof | grep /dev/cbd50gigdisk1   > /dev/null  
    ret=$?
    if [ "$ret" == "1" ]; then  
      break
    fi
    logger found somebody using /dev/cbd50gigdisk1, waiting and checking again
    sleep 1
done

logger cbd50gigdisk is freed up

/opt/cbd/compressblockdevice/cbd catalog stop -d cbd50gigdisk 2>&1 | logger

logger systemd shutdown cbd50gigdisk service completed

