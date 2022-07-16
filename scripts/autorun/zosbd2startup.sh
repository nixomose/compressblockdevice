#!/bin/bash

ZOSBD2PATH=/opt/zosbd2

CLONEZOSBD2="git clone http://github.com/nixomose/zosbd2"
ZOSBD2SRC=zosbd2/src

lsmod | grep zosbd2
if [ $? == 0 ]; then
  echo "zosbd2 is loaded. nothing to do."
  exit 0
fi

WHOAMI=`whoami`
if [ "$WHOAMI" != "root" ]; 
then echo "you must be root.";
exit;
fi


d=`dirname "$0"`
d=`cd $d && pwd`

# go to where the script is
cd $d


if [ ! -d $ZOSBD2PATH ]; then
  mkdir $ZOSBD2PATH
fi

cd $ZOSBD2PATH
if [ $? != 0 ]; then
  echo "unable to cd to $ZOSBD2PATH"
  exit 1
fi

test -d $ZOSBD2SRC
if [ $? != 0 ]; then
  $CLONEZOSBD2
  if [ $? != 0 ]; then
    echo "$CLONEZOSBD2 failed"
    exit 1
  fi
fi

cd $ZOSBD2SRC
if [ $? != 0 ]; then
  echo "unable to cd to $ZOSBD2SRC"
  exit 1
fi

# this might fail because its not there or because 
# the kernel was updated
insmod zosbd2.ko
if [ $? != 0 ]; then
  make
  if [ $? != 0 ]; then
    echo "build zosbd2 failed"
    exit 1
  fi
  insmod zosbd2.ko
  if [ $? != 0 ]; then
    echo "insmod zosbd2 failed"
    exit 1
  fi
fi

lsmod | grep zosbd2
if [ $? != 0 ]; then
  echo "zosbd2 not running"
  exit 1
fi
