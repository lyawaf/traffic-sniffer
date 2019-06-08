package parser

import (
	"regexp"
)

type OwnerType string

const (
	Client = OwnerType("Client")
	Server = OwnerType("Server")
)

type Packet struct {
	Owner OwnerType
	Data  string
}

type TCPSession struct {
	ServerAddr     string
	ClientAddr     string
	ServerPort     uint16
	ClientPort     uint16
	SequenceNumber uint32
	Packets        []Packet
	Labels         []Label
}

// LabelType is marker for applying regexp:
// for IN or for OUT
type LabelType string

const (
	PacketIN  = LabelType("in")
	PacketOUT = LabelType("out")
)

// Label uses for traffic clustering.
type Label struct {
	Name   string
	Type   LabelType
	Color  string
	Regexp *regexp.Regexp
	Count  int
}

var Labels = []Label{
	{
		Name:   "FlagIN",
		Type:   PacketIN,
		Regexp: regexp.MustCompile("[A-Z0-9]{31}="),
	},
}
