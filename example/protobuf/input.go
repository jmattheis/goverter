package protobuf

import "goverter/example/pb"

// goverter:converter
type Converter interface {
	FromProtobuf(*pb.Event) *OutputEvent

	// goverter:ignore state sizeCache unknownFields
	ToProtobuf(*OutputEvent) *pb.Event
}

type OutputEvent struct {
	Content  string
	Priority int32
}
