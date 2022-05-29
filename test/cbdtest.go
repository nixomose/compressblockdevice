// SPDX-License-Identifier: LGPL-2.1
// Copyright (C) 2021-2022 stu mark

package main

import (
	"github.com/nixomose/compressblockdevice/compressblockdevice/cbdkompressor"
	"github.com/nixomose/compressblockdevice/test/backend"
	"github.com/nixomose/nixomosegotools/tools"
)

func main() {
	backend.Cbd_backend_tst()

	kompressor_test()
}

func kompressor_test() {

	var log *tools.Nixomosetools_logger = tools.New_Nixomosetools_logger(0)

	var comp *cbdkompressor.Compression_pipeline_element = cbdkompressor.New_compression_pipeline(log, "cbdkompressortest.cf")
	var ret = comp.Init()
	if ret != nil {
		return
	}

	var block_size uint32 = 100
	ret = comp.Init_block_size(block_size)
	if ret != nil {
		return
	}
	var test []byte = make([]byte, block_size)
	for i := 0; i < len(test); i++ {
		test[i] = byte(i)
	}

	go_compress(comp, &test)
	go_decompress(comp, &test)

	block_size = 550
	ret = comp.Init_block_size(block_size)
	if ret != nil {
		return
	}

	var test2 []byte = make([]byte, block_size)
	for i := 0; i < len(test2); i++ {
		test2[i] = byte(i)
	}

	go_compress(comp, &test2)
	go_decompress(comp, &test2)
}

func go_compress(comp *cbdkompressor.Compression_pipeline_element, test *[]byte) {

	var compressed_length int
	var ret tools.Ret = comp.Pipe_in(test)
	if ret != nil {
		panic(ret.Get_errmsg())
	}
	compressed_length = len(*test)
	print("length of first compression test: ", compressed_length, "\n")

}

func go_decompress(comp *cbdkompressor.Compression_pipeline_element, test *[]byte) {

	var compressed_length int
	var ret tools.Ret = comp.Pipe_out(test)
	if ret != nil {
		panic(ret.Get_errmsg())
	}
	compressed_length = len(*test)
	print("length of first compression test: ", compressed_length, "\n")
}
