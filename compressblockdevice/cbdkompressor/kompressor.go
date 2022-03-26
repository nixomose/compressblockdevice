// Package cbdkompressor comment
package cbdkompressor

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/nixomose/nixomosegotools/tools"
	"github.com/nixomose/zosbd2goclient/zosbd2cmdlib/zosbd2interfaces"
	"github.com/spf13/cobra"
)

/* a compression pipeline */

type Compression_pipeline_element struct {
	log         *tools.Nixomosetools_logger
	config_file string

	max_node_size_in_bytes uint32 // the number of bytes we expect to be compressing and the byte size we will return when decompressing

	conf *Compressor_config

	context *Compression_pipeline_element_context
}

type Compression_pipeline_element_context struct {
	// stuff int
}

// verify that compression pipeline_element implements data_pipeline_element
var _ zosbd2interfaces.Data_pipeline_element = &Compression_pipeline_element{}
var _ zosbd2interfaces.Data_pipeline_element = (*Compression_pipeline_element)(nil)

// verify that compression pipeline_element_context implements data_pipeline_element_context
var _ zosbd2interfaces.Data_pipeline_element_context = &Compression_pipeline_element_context{}
var _ zosbd2interfaces.Data_pipeline_element_context = (*Compression_pipeline_element_context)(nil)

const COMPRESSION_TYPE_NONE = "none"
const COMPRESSION_TYPE_GZIP = "gzip"

// config settings for the compression pipeline

type Compressor_config struct {
	Compression_settings compression_settings
}
type compression_settings struct {
	Compression_level int
	Compression_type  string
}

func New_compression_pipeline(log *tools.Nixomosetools_logger, config_file string) *Compression_pipeline_element {
	return &Compression_pipeline_element{log: log, config_file: config_file, max_node_size_in_bytes: 0}
}

func (this *Compression_pipeline_element) Process_parameters(params *cobra.Command) tools.Ret {
	// take any command line params and add them to the context
	return nil
}

func (this *Compression_pipeline_element) Process_device(device zosbd2interfaces.Device_interface) tools.Ret {
	// take the block size from the device and store it so we know if things compressed or not.

	this.max_node_size_in_bytes = device.Get_node_size_in_bytes()
	return nil
}

func (this *Compression_pipeline_element) Init() tools.Ret {
	var ret = this.parse_config_file()

	return ret
}

func (this *Compression_pipeline_element) Init_block_size(max_node_size_in_bytes uint32) tools.Ret {
	this.max_node_size_in_bytes = max_node_size_in_bytes
	return nil
}

func (this *Compression_pipeline_element) Pipe_in(data_in_out *[]byte) tools.Ret {
	/* we only compress whole blocks, because we don't want to store metadata about the original size per block.
	   if it doesn't compress to smaller, than store the original size, and if it's the original size
		 then it's not compressed. we update the callers slice so any avoid-copying optimizations
		 can be made in the future. */

	if this.max_node_size_in_bytes == 0 {
		return tools.Error(this.log, "kompressor not initialized correctly")
	}

	// https://stackoverflow.com/questions/19197874/how-can-i-use-gzip-on-a-string
	var b bytes.Buffer
	var original_length uint32 = uint32(len(*data_in_out))
	if original_length != this.max_node_size_in_bytes {
		return tools.Error(this.log, "incoming block is ", original_length, " not ", this.max_node_size_in_bytes, " bytes")
	}

	var gz, err = gzip.NewWriterLevel(&b, gzip.DefaultCompression)
	if err != nil {
		return tools.Error(this.log, "invalid compression level, error: ", err)
	}
	var bytes_processed int
	if bytes_processed, err = gz.Write(*data_in_out); err != nil {
		return tools.Error(this.log, "error compressing data, error: ", err)
	}
	// close causes flush, otherwise gz.Flush()
	if err := gz.Close(); err != nil {
		return tools.Error(this.log, "error compressing block, error: ", err)
	}
	if uint32(bytes_processed) != original_length {
		return tools.Error(this.log, "error compressing block, only ", bytes_processed, " of ", original_length, "compressed")
	}
	if b.Len() >= len(*data_in_out) {
		// compression made it worse, just return the original length, data in place
		// if it's equal size, also don't compress because equal size means no compression
		return nil
	}
	// put a magic header saying how this was encrypted
	var count int = copy(*data_in_out, b.Bytes())
	if count != b.Len() {
		return tools.Error(this.log, "didn't copy entire compressed buffer, expected: ", b.Len(),
			" only copied: ", count)
	}
	// set caller's slice to new length
	*data_in_out = (*data_in_out)[0:count]
	return nil
}

func (this *Compression_pipeline_element) Pipe_out(data_in_out *[]byte) tools.Ret {
	// decompress/decrypt

	if this.max_node_size_in_bytes == 0 {
		return tools.Error(this.log, "kompressor not initialized correctly")
	}

	if uint32(len(*data_in_out)) == this.max_node_size_in_bytes {
		// if incoming block is the correct size, then it is not compressed
		return nil
	}

	// reader := bytes.NewReader(data_in_out)
	var byte_reader = bytes.NewReader(*data_in_out)
	gz, err := gzip.NewReader(byte_reader)
	if err != nil {
		return tools.Error(this.log, "error decompressing data, error: ", err)
	}

	output, err := ioutil.ReadAll(gz)
	if err != nil {
		return tools.Error(this.log, "error reading decompressed data, error: ", err)
	}

	if err := gz.Close(); err != nil {
		return tools.Error(this.log, "error closing gzip decompressor, error: ", err)
	}

	if uint32(len(output)) != this.max_node_size_in_bytes {
		return tools.Error(this.log, "decompressed data is not correct size, expected ", this.max_node_size_in_bytes,
			" got ", len(output))
	}

	// if it's not big enough make it bigger
	(*data_in_out) = (*data_in_out)[0:this.max_node_size_in_bytes]

	var count int = copy(*data_in_out, output)
	if count != len(output) {
		return tools.Error(this.log, "didn't copy entire decompressed buffer, expected: ", len(output),
			" only copied: ", count)
	}
	return nil
}

func (this *Compression_pipeline_element) Get_context() zosbd2interfaces.Data_pipeline_element_context {
	var ret = this.context
	return ret
}
func (this *Compression_pipeline_element) Set_context(context zosbd2interfaces.Data_pipeline_element_context) {
	this.context = context.(*Compression_pipeline_element_context)
}

func (this *Compression_pipeline_element) parse_config_file() tools.Ret {
	_, err := toml.DecodeFile(this.config_file, &this.conf)
	if err != nil {
		return tools.Error(this.log, "Unable to parse config file: ", this.config_file, ", err: ", err)
	}
	if this.conf.Compression_settings.Compression_level < 0 ||
		this.conf.Compression_settings.Compression_level > 10 {
		return tools.Error(this.log, "invalid compression level: ", this.conf.Compression_settings.Compression_level,
			" in ", this.config_file)
	}

	if this.conf.Compression_settings.Compression_type != COMPRESSION_TYPE_NONE &&
		this.conf.Compression_settings.Compression_type != COMPRESSION_TYPE_GZIP {
		return tools.Error(this.log, " invalid compression type: ", this.conf.Compression_settings.Compression_type,
			" in ", this.config_file)
	}

	return nil
}

//we need to make a context, and we need to put the magic for this application in there for the stree

func (this *Compression_pipeline_element_context) Create() tools.Ret {
	return nil
}

func (this *Compression_pipeline_element_context) Get_context() zosbd2interfaces.Data_pipeline_element_context {
	return this
}
