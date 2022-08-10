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

# catalog

the catalog is just a list of definitions about a block device tied to a name so that you can define the block device characteristics once, and from then on refer to the block device definition by name. most regular usage of compressbd is through the catalog commands.

```
  add         add a catalog entry with this block device definition
  delete      delete the specified catalog entry and its backing store
  list        list one or all of the devices defined in the catalog
  set         set a configuration or catalog entry option
  start       create a block device with the definition specified by the device name in the catalog
  stop        cleanly shutdown a currently running block device specified by the device name
```

# lower level access


  ```destroy-device``` 
this command tells the kernel module to cleanly end a block device. it will attempt to flush any in-flight requests but if they take too long to process, the block device will be pulled out from under you. so you should only call this once you've synced and unmounted any filesystems mounted to the block device.

  
  ```destroy-all-devices```  

same as above, but you don't specify a device name, it just goes through and destroys all the block devices, attempting to flush data first.
  
  
 
  
  ```device-status```
  
display status of all existing block devices.
 
example output:
  
```
./cbd device-status 
2022/08/09 20:51:11 {"msg":"getting status for all block devices."}
2022/08/09 20:51:11 {"msg":"getting status for handle_id 2"}
{
 "z5": {
  "size_in_bytes": 10737418240,
  "number_of_blocks": 2621440,
  "kernel_block_size_in_bytes": 4096,
  "max_segments_per_request": 256,
  "timeout_in_milliseconds": 1200000,
  "handle_id": 2,
  "device_name": "z5"
 }
}
```  
  
  
  
  
  diag                diagnostic tools
  help                Help about any command
  storage-status      display definition of backing storage
