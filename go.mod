module github.com/nixomose/compressblockdevice

replace github.com/nixomose/blockdevicelib => ../blockdevicelib

replace github.com/nixomose/stree_v => ../stree_v

replace github.com/nixomose/zosbd2goclient => ../zosbd2goclient

replace github.com/nixomose/nixomosegotools => ../nixomosegotools

go 1.17

require (
	github.com/nixomose/blockdevicelib v0.0.0-20220520234636-cf4290d3318c
	github.com/nixomose/nixomosegotools v0.0.0-20220511015301-07588dc19640
)

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/ncw/directio v1.0.5 // indirect
	github.com/nixomose/stree_v v0.0.0-20220517112730-e00d5fe38978 // indirect
	github.com/nixomose/zosbd2goclient v0.0.0-20220511020108-ecfb34c03dfa
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)
