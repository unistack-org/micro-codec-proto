// Package proto provides a proto codec
package proto

import (
	"io"

	"github.com/unistack-org/micro/v3/codec"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
	"google.golang.org/protobuf/proto"
)

type protoCodec struct{}

const (
	flattenTag = "flatten"
)

func (c *protoCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	case proto.Message:
		if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
			if nm, ok := nv.(proto.Message); ok {
				m = nm
			}
		}
		return proto.Marshal(m)
	}
	return nil, codec.ErrInvalidMessage
}

func (c *protoCodec) Unmarshal(d []byte, v interface{}) error {
	if len(d) == 0 {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = d
	case proto.Message:
		if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
			if nm, ok := nv.(proto.Message); ok {
				m = nm
			}
		}
		return proto.Unmarshal(d, m)
	}
	return codec.ErrInvalidMessage
}

func (c *protoCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *protoCodec) ReadBody(conn io.Reader, v interface{}) error {
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		buf, err := io.ReadAll(conn)
		if err != nil {
			return err
		} else if len(buf) == 0 {
			return nil
		}
		m.Data = buf
		return nil
	case proto.Message:
		buf, err := io.ReadAll(conn)
		if err != nil {
			return err
		} else if len(buf) == 0 {
			return nil
		}
		if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
			if nm, ok := nv.(proto.Message); ok {
				m = nm
			}
		}
		return proto.Unmarshal(buf, m)
	}
	return codec.ErrInvalidMessage
}

func (c *protoCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		_, err := conn.Write(m.Data)
		return err
	case proto.Message:
		if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
			if nm, ok := nv.(proto.Message); ok {
				m = nm
			}
		}

		buf, err := proto.Marshal(m)
		if err != nil {
			return err
		} else if len(buf) == 0 {
			return codec.ErrInvalidMessage
		}
		_, err = conn.Write(buf)
		return err
	}
	return codec.ErrInvalidMessage
}

func (c *protoCodec) String() string {
	return "proto"
}

func NewCodec() codec.Codec {
	return &protoCodec{}
}
