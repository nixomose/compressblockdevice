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
At the moment you have to be root to run cbd mostly because of interacting with the kernel module block device and the control device.
On my list of things to do is make it not require root, that should be possible as long as you can convey the user you want to use it as to the kernel module the first time. chicken and egg problem for another day.

```
./cbd  --help
cbde allows you to define a catalog of block devices defining their size and backing store and lets you 
easily start up and shut down these block devices. 
requires zosbd2 - https://github.com/nixomose/zosbd2

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
  
  
```diag```

diagnostic tools

dump the header of the stree backing the block device.

```
./cbd diag dump header -d z5
    0: 5a454e53 54354b41 00000062 c2010000 00000004 00010030 0062af80 00000000  | ZENST5KA░░░b░░░░░░░░░░░0░b░░░░░░
   32: 00000001 00010030 00000001                                               | ░░░░░░░0░░░░

{
 "0001_key": "    0: 5a454e53 54354b41  | ZENST5KA\n",
 "0002_store_size_in_bytes": "424,161,640,448 0x00000062c2010000",
 "0003_nodes_per_block": "4 0x00000004",
 "0004_block_size": "65,584 0x00010030",
 "0005_block_count": "6,467,456 0x0062af80",
 "0006_root_node": "0 0x00000000",
 "0007_free_position": "1 0x00000001",
 "0008_alignment": "65,584 0x00010030",
 "0009_dirty": "1 0x00000001"
}
```




```storage-status```

display definition of backing storage
  
``` ./cbd storage-status -t /home/nixo/testzosz5 
{
 "backing_storage": "/home/nixo/testzosz5",
 "dirty": "1",
 "inital_block_size_in_bytes": "65,584",
 "inital_nodes_per_block": "4",
 "inital_store_size_in_bytes": "424,161,640,448",
 "node_size_in_bytes": "327,920",
 "number_of_blocks_available_in_backing_store": "6,467,456",
 "number_of_physical_bytes_used_for_a_block": "65,584",
 "physical_store_block_alignment": "65,584",
 "total_bytes_wasted_due_to_alignment_padding": "0",
 "total_waste_percent": "0",
 "wasted_bytes_per_block": "0"
}
```  
  
  
