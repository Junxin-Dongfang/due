package msgpack

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

const Name = "msgpack"

var DefaultCodec = &codec{}

type codec struct{}

// Name 编解码器名称
func (codec) Name() string {
	return Name
}

// Marshal 编码
func (codec) Marshal(v any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := msgpack.NewEncoder(&buffer)
	// Preserve the compact integer and floating-point representation emitted by
	// the previous codec so existing peers keep seeing the same MessagePack wire
	// values while using a decoder without the truncated-fixext panic.
	encoder.UseCompactInts(true)
	encoder.UseCompactFloats(true)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Unmarshal 解码
func (codec) Unmarshal(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}

// Marshal 编码
func Marshal(v any) ([]byte, error) {
	return DefaultCodec.Marshal(v)
}

// Unmarshal 解码
func Unmarshal(data []byte, v any) error {
	return DefaultCodec.Unmarshal(data, v)
}
