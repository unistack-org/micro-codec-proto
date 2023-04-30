// Package proto provides a proto codec
package proto // import "go.unistack.org/micro-codec-proto/v4"

import (
	"io"

	pb "go.unistack.org/micro-proto/v4/codec"
	"go.unistack.org/micro/v4/codec"
	rutil "go.unistack.org/micro/v4/util/reflect"
	"google.golang.org/protobuf/proto"
)

type protoCodec struct {
	opts codec.Options
}

var _ codec.Codec = &protoCodec{}

const (
	flattenTag = "flatten"
)

func (c *protoCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	}

	if _, ok := v.(proto.Message); !ok {
		return nil, codec.ErrInvalidMessage
	}

	return proto.Marshal(v.(proto.Message))
}

func (c *protoCodec) Unmarshal(d []byte, v interface{}, opts ...codec.Option) error {
	if v == nil || len(d) == 0 {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = d
		return nil
	case *pb.Frame:
		m.Data = d
		return nil
	}

	if _, ok := v.(proto.Message); !ok {
		return codec.ErrInvalidMessage
	}

	return proto.Unmarshal(d, v.(proto.Message))
}

func (c *protoCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *protoCodec) ReadBody(conn io.Reader, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}
	return c.Unmarshal(buf, v)
}

func (c *protoCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := c.Marshal(v)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return codec.ErrInvalidMessage
	}
	_, err = conn.Write(buf)
	return err
}

func (c *protoCodec) String() string {
	return "proto"
}

func NewCodec(opts ...codec.Option) codec.Codec {
	return &protoCodec{opts: codec.NewOptions(opts...)}
}
