// Package proto provides a proto codec
package proto

import (
	pb "go.unistack.org/micro-proto/v5/codec"
	"go.unistack.org/micro/v5/codec"
	rutil "go.unistack.org/micro/v5/util/reflect"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/mem"
	"google.golang.org/protobuf/proto"
)

var (
	DefaultMarshalOptions = proto.MarshalOptions{
		AllowPartial: false,
	}

	DefaultUnmarshalOptions = proto.UnmarshalOptions{
		DiscardUnknown: false,
		AllowPartial:   false,
	}
)

type protoCodec struct {
	opts codec.Options
}

type protoCodecV2 struct {
	opts codec.Options
}

type protoCodecGRPC struct {
	opts codec.Options
	bp   mem.BufferPool
}

var (
	_ codec.Codec      = (*protoCodec)(nil)
	_ codec.CodecV2    = (*protoCodecV2)(nil)
	_ encoding.CodecV2 = (*protoCodecGRPC)(nil)
)

func (c *protoCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	case proto.Message:
		marshalOptions := DefaultMarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(marshalOptionsKey{}).(proto.MarshalOptions); ok {
				marshalOptions = f
			}
		}
		return marshalOptions.Marshal(m)
	case codec.RawMessage:
		return []byte(m), nil
	case *codec.RawMessage:
		return []byte(*m), nil
	default:
		return nil, codec.ErrInvalidMessage
	}
}

func (c *protoCodec) Unmarshal(d []byte, v interface{}, opts ...codec.Option) error {
	if v == nil || len(d) == 0 {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = d
		return nil
	case *pb.Frame:
		m.Data = d
		return nil
	case proto.Message:
		unmarshalOptions := DefaultUnmarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(marshalOptionsKey{}).(proto.UnmarshalOptions); ok {
				unmarshalOptions = f
			}
		}
		return unmarshalOptions.Unmarshal(d, m)
	case *codec.RawMessage:
		*m = append((*m)[0:0], d...)
		return nil
	case codec.RawMessage:
		copy(m, d)
		return nil
	default:
		return codec.ErrInvalidMessage
	}
}

func (c *protoCodec) String() string {
	return "proto"
}

func NewCodec(opts ...codec.Option) codec.Codec {
	return &protoCodec{opts: codec.NewOptions(opts...)}
}

func (c *protoCodecV2) Marshal(d []byte, v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	case proto.Message:
		marshalOptions := DefaultMarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(marshalOptionsKey{}).(proto.MarshalOptions); ok {
				marshalOptions = f
			}
		}
		return marshalOptions.MarshalAppend(d[:0], m)
	default:
		return nil, codec.ErrInvalidMessage
	}
}

func (c *protoCodecV2) Unmarshal(d []byte, v interface{}, opts ...codec.Option) error {
	if v == nil || len(d) == 0 {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = d
		return nil
	case *pb.Frame:
		m.Data = d
		return nil
	case proto.Message:
		unmarshalOptions := DefaultUnmarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(marshalOptionsKey{}).(proto.UnmarshalOptions); ok {
				unmarshalOptions = f
			}
		}
		return unmarshalOptions.Unmarshal(d, m)
	default:
		return codec.ErrInvalidMessage
	}
}

func (c *protoCodecV2) String() string {
	return "proto"
}

func NewCodecV2(opts ...codec.Option) codec.CodecV2 {
	options := codec.NewOptions(opts...)

	bp, ok := options.Context.Value(memBufferPoolKey{}).(mem.BufferPool)
	if !ok || bp != nil {
		bp = mem.NewTieredBufferPool(1024, 2048, 4096, 8092, 1*1024*1024*1024, 4*1024*1024*1024)
	}

	return &protoCodecV2{opts: codec.NewOptions(opts...)}
}

func NewCodecGRPC(opts ...codec.Option) *protoCodecGRPC {
	options := codec.NewOptions(opts...)

	bp, ok := options.Context.Value(memBufferPoolKey{}).(mem.BufferPool)
	if !ok || bp != nil {
		bp = mem.NewTieredBufferPool(1024, 2048, 4096, 8092, 1*1024*1024*1024, 4*1024*1024*1024)
	}

	return &protoCodecGRPC{opts: codec.NewOptions(opts...)}
}

func (c *protoCodecGRPC) Marshal(v any) (out mem.BufferSlice, err error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		buf := c.bp.Get(len(m.Data))
		*buf = append((*buf)[:0], m.Data...)
		out = append(out, mem.NewBuffer(buf, c.bp))
		return out, nil
	case *pb.Frame:
		buf := c.bp.Get(len(m.Data))
		*buf = append((*buf)[:0], m.Data...)
		out = append(out, mem.NewBuffer(buf, c.bp))
		return out, nil
	case proto.Message:
		marshalOptions := DefaultMarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(marshalOptionsKey{}).(proto.MarshalOptions); ok {
				marshalOptions = f
			}
		}
		size := proto.Size(m)

		buf := c.bp.Get(size)

		if _, err := marshalOptions.MarshalAppend((*buf)[:0], m); err != nil {
			c.bp.Put(buf)
			return nil, err
		}
		out = append(out, mem.NewBuffer(buf, c.bp))
		return out, nil
	default:
		return nil, codec.ErrInvalidMessage
	}
}

func (c *protoCodecGRPC) Unmarshal(data mem.BufferSlice, v any) (err error) {
	if v == nil || len(data) == 0 {
		return nil
	}

	options := c.opts

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	buf := data.MaterializeToBuffer(c.bp)
	defer buf.Free()

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = buf.ReadOnlyData()
		return nil
	case *pb.Frame:
		m.Data = buf.ReadOnlyData()
		return nil
	case proto.Message:
		unmarshalOptions := DefaultUnmarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(marshalOptionsKey{}).(proto.UnmarshalOptions); ok {
				unmarshalOptions = f
			}
		}
		return unmarshalOptions.Unmarshal(buf.ReadOnlyData(), m)
	default:
		return codec.ErrInvalidMessage
	}
}

func (c *protoCodecGRPC) Name() string {
	return "proto"
}
