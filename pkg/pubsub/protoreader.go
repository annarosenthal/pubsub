package pubsub

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

type ProtoReader struct {
	reader io.Reader
}

func NewProtoReader(r io.Reader) *ProtoReader {
	return &ProtoReader{reader: r}
}

func (pr *ProtoReader) Read (msg proto.Message) error {
	length, err := pr.readLength()
	if err != nil {
		return err
	}
	buffer := make([]byte, length)
	_, err = pr.reader.Read(buffer)
	if err != nil {
		return err
	}
	return proto.Unmarshal(buffer, msg)
}

func (pr *ProtoReader) readLength() (int, error) {
	buffer := make([]byte, 4)
	_, err := pr.reader.Read(buffer)
	if err != nil {
		return 0, err
	}
	length := binary.LittleEndian.Uint32(buffer)
	return int(length), nil
}