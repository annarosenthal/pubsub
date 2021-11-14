package pubsub

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

type ProtoWriter struct {
	writer io.Writer
}

func NewProtoWriter(w io.Writer) *ProtoWriter {
	return &ProtoWriter{writer: w}
}

//framing
func (pw *ProtoWriter) Write(msg proto.Message) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = pw.writeLength(bytes)
	if err != nil {
		return err
	}

	_, err = pw.writer.Write(bytes)
	return err
}

func (pw *ProtoWriter) writeLength(bytes []byte) error {
	length := len(bytes)
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(length))
	_, err := pw.writer.Write(buffer)
	return err
}
