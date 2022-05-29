// SPDX-License-Identifier: LGPL-2.1
// Copyright (C) 2021-2022 stu mark

// Package backend has a comment
package backend

import (
	"bytes"
	"container/list"

	"github.com/nixomose/blockdevicelib/blockdevicelib"
	"github.com/nixomose/compressblockdevice/compressblockdevice/cbdkompressor"
	"github.com/nixomose/nixomosegotools/tools"
	"github.com/nixomose/stree_v/stree_v_lib/stree_v_lib"

	//	stree_v_lib "github.com/nixomose/stree_v/stree_v_lib/stree_v_interfaces"
	"github.com/nixomose/zosbd2goclient/zosbd2_stree_v_storage_mechanism"
	"github.com/nixomose/zosbd2goclient/zosbd2cmdlib/zosbd2interfaces"
)

func Cbd_backend_tst() {
	Cbd_test_write_discard()
	Cbd_test_write_read()
}

func Cbd_test_write_discard() {
	var ret tools.Ret
	var l *blockdevicelib.Lbd_lib
	ret, l = blockdevicelib.New_blockdevicelib("test_backend")
	if ret != nil {
		return
	}
	var log *tools.Nixomosetools_logger = tools.New_Nixomosetools_logger(tools.DEBUG)
	l.Set_log(log)

	var comp *cbdkompressor.Compression_pipeline_element = cbdkompressor.New_compression_pipeline(l.Get_log(), "cbdkompressortest.cf")
	ret = comp.Init()
	if ret != nil {
		return
	}

	var data_block_size uint32 = 1048576
	var z zosbd2interfaces.Storage_mechanism
	var stree *stree_v_lib.Stree_v = nil

	ret = test_a_size(log, l, comp, data_block_size, &z, &stree)
	if ret != nil {
		panic("test a size failed for" + tools.Uint32tostring(data_block_size))
	}
	ret = z.Discard_block(17395712, 116822016)
	if ret != nil {
		panic("discard failed for" + tools.Uint32tostring(data_block_size))
	}

	ret = z.Discard_block(0, 1048576)
	if ret != nil {
		panic("discard failed for" + tools.Uint32tostring(data_block_size))
	}
	ret = stree.Shutdown()
	if ret != nil {
		log.Error("Unable to shutdown backing storage, error: ", ret.Get_errmsg())
	}

}
func Cbd_test_write_read() {

	var ret tools.Ret
	var l *blockdevicelib.Lbd_lib
	ret, l = blockdevicelib.New_blockdevicelib("test_backend")
	if ret != nil {
		return
	}
	var log *tools.Nixomosetools_logger = tools.New_Nixomosetools_logger(tools.DEBUG)
	l.Set_log(log)

	var comp *cbdkompressor.Compression_pipeline_element = cbdkompressor.New_compression_pipeline(l.Get_log(), "cbdkompressortest.cf")
	ret = comp.Init()
	if ret != nil {
		return
	}

	var data_block_size uint32 = 100
	var z zosbd2interfaces.Storage_mechanism
	var stree *stree_v_lib.Stree_v = nil

	ret = test_a_size(log, l, comp, data_block_size, &z, &stree)
	if ret != nil {
		panic("1m compress failed: " + ret.Get_errmsg())
	}

	ret = stree.Shutdown()
	if ret != nil {
		log.Error("Unable to shutdown backing storage, error: ", ret.Get_errmsg())
	}

	data_block_size = 65536

	ret = test_a_size(log, l, comp, data_block_size, &z, &stree)
	if ret != nil {
		panic("64k compress failed: " + ret.Get_errmsg())
	}
	ret = stree.Shutdown()
	if ret != nil {
		log.Error("Unable to shutdown backing storage, error: ", ret.Get_errmsg())
	}

	data_block_size = 1048576

	ret = test_a_size(log, l, comp, data_block_size, &z, &stree)
	if ret != nil {
		panic("test a size failed for" + tools.Uint32tostring(data_block_size))
	}
	ret = stree.Shutdown()
	if ret != nil {
		log.Error("Unable to shutdown backing storage, error: ", ret.Get_errmsg())
	}

}

func test_a_size(log *tools.Nixomosetools_logger, l *blockdevicelib.Lbd_lib, comp *cbdkompressor.Compression_pipeline_element,
	data_block_size uint32, zout *zosbd2interfaces.Storage_mechanism,
	stree **stree_v_lib.Stree_v) tools.Ret {

	var pipeline list.List
	pipeline.PushBack(comp)

	var calculated_stree_node_size uint32 = 0 // this gets calculated in make_stree
	var additional_nodes_per_block uint32 = 50

	/* the number of bytes we can store in one stree tree node is 1024 * 51 as there is room for 50 offspring.
	we needed to make a large block size because my test data compresses really well (not realistic for normal data)
	and wasn't able to compress to something larger than 1k when we only had 5 additional offspring.*/
	var stree_value_size uint32 = 1024
	var alignment = stree_value_size

	var device = l.New_block_device("test_device", 1048576, "no storage file use ramdisk", false, false,
		alignment, stree_value_size, calculated_stree_node_size, additional_nodes_per_block,
		false, "", false, true) // use the stree ramdisk

	var ret tools.Ret
	ret, *stree = l.Make_stree(device)
	if ret != nil {
		return ret
	}

	/* as this is a ramdisk backing the stree, we must initialize it. */
	ret = (*stree).Init()
	if ret != nil {
		return ret
	}

	ret = l.Start_stree((*stree), false)
	if ret != nil {
		return ret
	}

	var stree_block_size uint32 = stree_value_size * (additional_nodes_per_block + 1)
	var stree_block_size2 uint32 = (*stree).Get_node_size_in_bytes()
	if stree_block_size != stree_block_size2 {
		return tools.Error(log, "block size calculation is off")
	}

	ret = comp.Init_block_size(stree_block_size)
	if ret != nil {
		return ret
	}

	/* now that we have an stree_v, we pass that to the zosbd2_backing_store */
	*zout = zosbd2_stree_v_storage_mechanism.New_zosbd2_storage_mechanism(log, (*stree), &pipeline)

	var testdata []byte = make([]byte, data_block_size)
	for i := 0; i < len(testdata); i++ {
		testdata[i] = byte((i * i) - i + (i % 5))
	}

	ret = (*zout).Write_block(0, data_block_size, testdata)
	if ret != nil {
		panic("write block size failed: " + ret.Get_errmsg())
	}

	// now read it back and make sure it matches
	var returndata []byte = make([]byte, data_block_size)

	ret = (*zout).Read_block(0, data_block_size, returndata)
	if ret != nil {
		panic("read block size failed: " + ret.Get_errmsg())
	}

	var compval int = bytes.Compare(testdata, returndata)
	if compval != 0 {
		panic("data read back in doesn't match data written out.")

	}

	return nil

}
