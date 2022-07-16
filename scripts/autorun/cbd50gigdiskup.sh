#!/bin/bash

DEVICENAME=cbd50gigdisk
STORAGEFILE=/opt/cbd/$DEVICENAME

CBDPATH=/opt/cbd
CBDBINPATH=$CBDPATH/compressblockdevice
CBDBIN=$CBDBINPATH/cbd

SERVICEPATH=/etc/systemd/system/
SERVICEFILE=cbd50gigdisk.service

UPFILE=cbd50gigdiskup.sh
DOWNFILE=cbd50gigdiskdown.sh
ZOSBD2FILE=zosbd2startup.sh

WHOAMI=`whoami`
if [ "$WHOAMI" != "root" ]; 
then echo "you must be root.";
exit;
fi

d=`dirname "$0"`
d=`cd $d && pwd`

# go to where the script and all the other files to setup are.
cd $d

if [ ! -d $CBDPATH ]; then
  mkdir $CBDPATH
fi

FIRSTTIMEIN=0
# set up for systemd startup and shutdown
if [ ! -f $SERVICEPATH$SERVICEFILE ]; then
  echo cp $SERVICEFILE $SERVICEPATH$SERVICEFILE
  cp $SERVICEFILE $SERVICEPATH$SERVICEFILE
  systemctl enable $SERVICEFILE 
  systemctl daemon-reload
  FIRSTTIMEIN=1
fi

if [ ! -f /usr/local/bin/$UPFILE ]; then
  cp $UPFILE /usr/local/bin/$UPFILE
fi
chmod 755 /usr/local/bin/$UPFILE

if [ ! -f /usr/local/bin/$DOWNFILE ]; then
  cp $DOWNFILE /usr/local/bin/$DOWNFILE
fi
chmod 755 /usr/local/bin/$DOWNFILE

if [ ! -f /usr/local/bin/$ZOSBD2FILE ]; then
  cp $ZOSBD2FILE /usr/local/bin/$ZOSBD2FILE
fi
chmod 755 /usr/local/bin/$ZOSBD2FILE

# see if it's already running

lsblk -o name | grep $DEVICENAME
if [ $? == 0 ]; then
  echo "$DEVICENAME already running, nothing to do"
  exit 0
fi

# setup the kernel module
$d/zosbd2startup.sh

if [ $? != 0 ]; then
  echo "unable to start up zosbd2"
  exit 1
fi

cd $CBDPATH
if [ $? != 0 ]; then
  echo "unable to cd to $CBDPATH"
  exit 1
fi

if [ ! -f $CBDBIN ]; then
  git clone http://github.com/nixomose/compressblockdevice
  git clone http://github.com/nixomose/blockdevicelib
  git clone http://github.com/nixomose/compressblockdevice
  git clone http://github.com/nixomose/nixomosegotools
  git clone http://github.com/nixomose/stree_v
  cd compressblockdevice
  ./build.sh
  if [ ! -f $CBDBIN ]; then
    echo "unable to build cbd"
    exit 1
  fi
fi

mkdir -p /etc/compressblockdevice

if [ ! -f  /etc/compressblockdevice/compressblockdevice.cf ]; then
  cp -r $CBDBINPATH/compressblockdevice.cf /etc/compressblockdevice/compressblockdevice.cf
fi 

$CBDBIN catalog list | jq . | grep $STORAGEFILE
if [ $? != 0 ]; then
  $CBDBIN catalog add -p 6 -d $DEVICENAME -t $STORAGEFILE -s $(( 1024*1024*1024*50 )) -e 16384
fi

if [ "$FIRSTTIMEIN" == "1" ]; then
  systemctl start cbd50gigdisk
else
  $CBDBIN catalog start -d $DEVICENAME $1
  partprobe /dev/$DEVICENAME
  # so nixo user programs can access it
fi

