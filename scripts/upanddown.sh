#!/bin/bash

while true; do
echo "writing zeroes..."
dd if=/dev/zero of=testfile bs=1M count=500 status=progress conv=notrunc; sync
ls -latrhs testfile /home/nixo/testzosdevicecompress
echo
echo "writing urandom..."
dd if=/dev/urandom of=testfile bs=1M count=500 status=progress conv=notrunc; sync
ls -latrhs testfile /home/nixo/testzosdevicecompress
echo
echo
echo "sleeping..."
sleep 10
done
