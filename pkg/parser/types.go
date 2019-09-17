package parser

import (
	"regexp"
	"sync"

	"github.com/google/gopacket"
)

type Parser struct {
	Source *gopacket.PacketSource
	sync.Mutex
	sessions []TCPSession
}

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
	LastUpdate     int64
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
	Name      string
	Type      LabelType
	Color     string
	Regexp    *regexp.Regexp
	RawRegexp string
}

var Labels = struct {
	sync.Mutex
	L []Label
}{L: []Label{
	{
		Name:      "ASDF label",
		Type:      PacketIN,
		Regexp:    regexp.MustCompile("asdf"),
		RawRegexp: "YXNkZg==",
		Color:     "#ffffff",
	},
	{
		Name:      "SQL quotes",
		Type:      PacketIN,
		Regexp:    regexp.MustCompile(`('(''|[^'])*')`),
		RawRegexp: "KCcoJyd8W14nXSkqJykK",
		Color:     "#ffffff",
	},
	{
		Name:      "SQL commands",
		Type:      PacketIN,
		Regexp:    regexp.MustCompile(`(\b(ALTER|CREATE|DELETE|DROP|EXEC(UTE){0,1}|INSERT( +INTO){0,1}|MERGE|SELECT|UPDATE|UNION( +ALL){0,1})\b)`),
		RawRegexp: "KCcoJyd8W14nXSkqJyl8KFxiKEFMVEVSfENSRUFURXxERUxFVEV8RFJPUHxFWEVDKFVURSl7MCwxfXxJTlNFUlQoICtJTlRPKXswLDF9fE1FUkdFfFNFTEVDVHxVUERBVEV8VU5JT04oICtBTEwpezAsMX0pXGIp",
		Color:     "#ffffff",
	},
}}
