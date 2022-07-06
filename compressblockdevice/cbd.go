// SPDX-License-Identifier: LGPL-2.1
// Copyright (C) 2021-2022 stu mark

package main

import (
	"container/list"
	"os"

	blockdevicelib "github.com/nixomose/blockdevicelib/blockdevicelib"
	"github.com/nixomose/compressblockdevice/compressblockdevice/cbdkompressor"
	"github.com/nixomose/nixomosegotools/tools"
	"github.com/spf13/cobra"
)

const TXT_APPLICATION_NAME = "cbd"

const TXT_DEFAULT_CONFIG_FILE = "/etc/compressblockdevice/compressblockdevice.cf"
const TXT_DEFAULT_LOG_FILE = "/var/log/compressblockdevice/compressblockdevice.log"
const TXT_DEFAULT_CATALOG_FILE = "/etc/compressblockdevice/catalog.toml"

/* this is the main binary entry point for local block device that uses stree_v to store the block device data
on a local disk or file. */

func main() {

	/* we have lots of space in the block zero header, and we're going to need to store
	   compression type and block size in there, so we need a way to get a callback about
	   metadata info that we store on init and get on load-startup so the compression stuff
	   can work as it is plugged in, based on metadata in the header.
	   so the header will have to know it has a compression payload, and then call compression to get the
	   specific data for that.
	   that means we need some kind of registry in the header.
	   So the header will now include an array of 'features' in this backing store file, and that list
	   will have an order, and each feature implemented will have a payload after the feature list
	   for it's specific data. So you have to pass in your feature, and a callbacky thing to set
	   and get your header data. */

	var ret tools.Ret
	var root_cmd *cobra.Command
	var l *blockdevicelib.Lbd_lib
	ret, l = blockdevicelib.New_blockdevicelib(TXT_APPLICATION_NAME)
	if ret != nil {
		os.Exit(1)
		return
	}
	ret, root_cmd = l.Startup(TXT_DEFAULT_CONFIG_FILE, TXT_DEFAULT_LOG_FILE,
		TXT_DEFAULT_CATALOG_FILE) // start configuring and make log and stuff.
	if ret != nil {
		os.Exit(1)
		return
	}
	var comp *cbdkompressor.Compression_pipeline_element = cbdkompressor.New_compression_pipeline(l.Get_log(),
		l.Get_config_file())
	ret = comp.Init()
	if ret != nil {
		os.Exit(1)
		return
	}

	var pipeline list.List
	pipeline.PushBack(comp)

	ret = l.Run(root_cmd, &pipeline)
	if ret != nil {
		os.Exit(1)
		return
	}
}
