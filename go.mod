module github.com/nixomose/compressblockdevice

replace github.com/nixomose/blockdevicelib => ../blockdevicelib

replace github.com/nixomose/stree_v => ../stree_v

replace github.com/nixomose/zosbd2goclient => ../zosbd2goclient

replace github.com/nixomose/nixomosegotools => ../nixomosegotools

go 1.17

require (
	github.com/nixomose/blockdevicelib v0.0.0-20220529234758-de0e6ba9bdbe
	github.com/nixomose/nixomosegotools v0.0.0-20220529231952-c38fcdca5407
)

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/ncw/directio v1.0.5 // indirect
	github.com/nixomose/stree_v v0.0.0-20220529234753-57276d134ab1
	github.com/nixomose/zosbd2goclient v0.0.0-20220529234745-b751c582f005
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)
