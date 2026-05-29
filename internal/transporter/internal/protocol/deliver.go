package protocol

import (
	"encoding/binary"
	"io"

	"github.com/dobyte/due/v2/core/buffer"
	"github.com/dobyte/due/v2/errors"
	"github.com/dobyte/due/v2/internal/transporter/internal/route"
	duetrace "github.com/dobyte/due/v2/trace"
)

const (
	deliverReqBytes      = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + b64 + b64
	deliverTraceReqBytes = deliverReqBytes + duetrace.ContextLength
	deliverResBytes      = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + defaultCodeBytes
)

// EncodeDeliverReq 编码投递消息请求
// 协议：size + header + route + seq + cid + uid + [trace context] + <message packet>
func EncodeDeliverReq(seq uint64, cid int64, uid int64, tc []byte, buf buffer.Buffer) *buffer.NocopyBuffer {
	header := dataBit
	reqBytes := deliverReqBytes
	if len(tc) == duetrace.ContextLength {
		header |= traceBit
		reqBytes = deliverTraceReqBytes
	}

	writer := buffer.MallocWriter(reqBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(reqBytes-defaultSizeBytes+buf.Len()))
	writer.WriteUint8s(header)
	writer.WriteUint8s(route.Deliver)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteInt64s(binary.BigEndian, cid, uid)
	if header&traceBit != 0 {
		writer.WriteBytes(tc...)
	}

	return buffer.NewNocopyBuffer(writer, buf)
}

// DecodeDeliverReq 解码投递消息请求
func DecodeDeliverReq(data []byte) (seq uint64, cid int64, uid int64, tc []byte, message []byte, err error) {
	reader := buffer.NewReader(data)

	if _, err = reader.Seek(defaultSizeBytes, io.SeekStart); err != nil {
		return
	}

	var header uint8
	if header, err = reader.ReadUint8(); err != nil {
		return
	}

	if _, err = reader.Seek(defaultRouteBytes, io.SeekCurrent); err != nil {
		return
	}

	if seq, err = reader.ReadUint64(binary.BigEndian); err != nil {
		return
	}
	if cid, err = reader.ReadInt64(binary.BigEndian); err != nil {
		return
	}
	if uid, err = reader.ReadInt64(binary.BigEndian); err != nil {
		return
	}

	messageOffset := deliverReqBytes
	if header&traceBit != 0 {
		if len(data) < deliverTraceReqBytes {
			err = errors.ErrInvalidMessage
			return
		}
		tc = append([]byte(nil), data[deliverReqBytes:deliverTraceReqBytes]...)
		messageOffset = deliverTraceReqBytes
	}
	message = data[messageOffset:]

	return
}

// EncodeDeliverRes 编码投递消息响应
// 协议：size + header + route + seq + code
func EncodeDeliverRes(seq uint64, code uint16) *buffer.NocopyBuffer {
	writer := buffer.MallocWriter(deliverResBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(deliverResBytes-defaultSizeBytes))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.Deliver)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint16s(binary.BigEndian, code)

	return buffer.NewNocopyBuffer(writer)
}

// DecodeDeliverRes 解码投递消息响应
// 协议：size + header + route + seq + code
func DecodeDeliverRes(data []byte) (code uint16, err error) {
	if len(data) != deliverResBytes {
		err = errors.ErrInvalidMessage
		return
	}

	reader := buffer.NewReader(data)

	if _, err = reader.Seek(-defaultCodeBytes, io.SeekEnd); err != nil {
		return
	}

	if code, err = reader.ReadUint16(binary.BigEndian); err != nil {
		return
	}

	return
}
