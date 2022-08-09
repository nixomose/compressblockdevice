# compressblockdevice
create a block device that compresses the data stored in it so it takes up less space that the block device exposes.

# Building
cd into compressblockdevice and 

`go build`

# Dependencies
all the go dependencies should be pulled in automatically upon building

The only other thing you need is zosbd2, the kernel module that creates the block device in the kernel and hands 
the requests to the userspace program to service them.


# Running
