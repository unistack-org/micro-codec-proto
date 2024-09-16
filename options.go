package proto

import (
	codec "go.unistack.org/micro/v3/codec"
	"google.golang.org/protobuf/proto"
)

type unmarshalOptionsKey struct{}

func UnmarshalOptions(o proto.UnmarshalOptions) codec.Option {
	return codec.SetOption(unmarshalOptionsKey{}, o)
}

type marshalOptionsKey struct{}

func MarshalOptions(o proto.MarshalOptions) codec.Option {
	return codec.SetOption(marshalOptionsKey{}, o)
}
