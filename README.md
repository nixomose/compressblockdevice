# compressblockdevice
create a block device that compresses the data stored in it so it takes up less space that the block device exposes.

# Building
cd into compressblockdevice and 

`time go build compressblockdevice/cbd.go`

# Dependencies
all the go dependencies should be pulled in automatically upon building

The only other thing you need is zosbd2, the kernel module that creates the block device in the kernel and hands 
the requests to the userspace program to service them.

https://github.com/nixomose/zosbd2 

Instructions to build the kernel module are all there.


# Running

the build produces a cbd binary which you can put in your path so you can just run cbd from anywhere...
At the moment you have to be root to run cbd mostly because of interacting with the kernel module block block and the control device.
On my list of things to do is make it not require root, that should be possible as long as you can convey the user you want to use it as to the kernel module the first time. chicken and egg problem for another day.

```
./cbd  --help
local block device allows you to define a catalog of block devices defining their size and backing store and lets you easily start up and shut down these block devices. requires zosbd2 - https
://github.com/nixomose/zosbd2

Usage:
  cbd [command]

Available Commands:
  catalog             list one or all of the devices defined in the catalog
  completion          Generate the autocompletion script for the specified shell
  destroy-all-devices destroy all block devices
  destroy-device      destroy a block device by name
  device-status       display status of all block devices
  diag                diagnostic tools
  help                Help about any command
  storage-status      display definition of backing storage

Flags:
  -c, --config-file string   configuration file (default "/etc/compressblockdevice/compressblockdevice.cf")
  -h, --help                 help for cbd
  -l, --log-file string      log file (default "/var/log/compressblockdevice/compressblockdevice.log")
  -v, --log-level uint32     log level: 0=debug 200=info 500=error (default 200)

```
