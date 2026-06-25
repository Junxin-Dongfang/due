package protocol_test

import (
	"bytes"
	"testing"

	"github.com/dobyte/due/v2/core/buffer"
	"github.com/dobyte/due/v2/internal/transporter/internal/codes"
	"github.com/dobyte/due/v2/internal/transporter/internal/protocol"
)

func TestEncodeDeliverReq(t *testing.T) {
	buffer := protocol.EncodeDeliverReq(1, 2, 3, nil, buffer.NewNocopyBuffer([]byte("hello world")))

	t.Log(buffer.Bytes())
}

func TestDecodeDeliverReq(t *testing.T) {
	buffer := protocol.EncodeDeliverReq(1, 2, 3, nil, buffer.NewNocopyBuffer([]byte("hello world")))

	seq, cid, uid, tc, message, err := protocol.DecodeDeliverReq(buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(tc) != 0 {
		t.Fatalf("trace context = %v, want empty", tc)
	}

	t.Logf("seq: %v", seq)
	t.Logf("cid: %v", cid)
	t.Logf("uid: %v", uid)
	t.Logf("message: %v", string(message))
}

func TestDeliverReqTraceContext(t *testing.T) {
	tc := make([]byte, 25)
	for i := range tc {
		tc[i] = byte(i + 1)
	}
	buffer := protocol.EncodeDeliverReq(1, 2, 3, tc, buffer.NewNocopyBuffer([]byte("hello world")))

	_, _, _, gotTC, message, err := protocol.DecodeDeliverReq(buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotTC, tc) {
		t.Fatalf("trace context = %v, want %v", gotTC, tc)
	}
	if string(message) != "hello world" {
		t.Fatalf("message = %q", string(message))
	}
}

func TestDeliverReqInvalidTraceContextUsesOldWire(t *testing.T) {
	before := protocol.EncodeDeliverReq(1, 2, 3, nil, buffer.NewNocopyBuffer([]byte("hello world"))).Bytes()
	after := protocol.EncodeDeliverReq(1, 2, 3, []byte("bad"), buffer.NewNocopyBuffer([]byte("hello world"))).Bytes()
	if !bytes.Equal(after, before) {
		t.Fatalf("wire bytes changed for invalid trace context\nbefore=%v\nafter=%v", before, after)
	}
}

func TestDecodeDeliverReqRejectsTruncatedTraceContext(t *testing.T) {
	tc := make([]byte, 25)
	buffer := protocol.EncodeDeliverReq(1, 2, 3, tc, buffer.NewNocopyBuffer([]byte("hello world")))
	data := buffer.Bytes()
	truncated := data[:54]

	if _, _, _, _, _, err := protocol.DecodeDeliverReq(truncated); err == nil {
		t.Fatal("expected truncated trace context to fail")
	}
}

func TestEncodeDeliverRes(t *testing.T) {
	buffer := protocol.EncodeDeliverRes(1, codes.OK)

	t.Log(buffer.Bytes())
}

func TestDecodeDeliverRes(t *testing.T) {
	buffer := protocol.EncodePushRes(1, codes.OK)

	code, err := protocol.DecodeDeliverRes(buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("code: %v", code)
}
