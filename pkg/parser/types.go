package parser

import (
	"regexp"
	"sync"

	"github.com/google/gopacket"
	"go.mongodb.org/mongo-driver/mongo"
)

const WAIT_TIMEOUT = 5

var DBClient *mongo.Client

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
	Name      string         `json:"name"`
	Type      LabelType      `json:"type"`
	Color     string         `json:"color"`
	Regexp    *regexp.Regexp `json:"-"`
	RawRegexp string         `json:"regexp"`
}

var Labels = struct {
	sync.Mutex
	L []Label
}{L: []Label{
	{
		Name:      "test label",
		Type:      PacketIN,
		Regexp:    regexp.MustCompile("asdf"),
		RawRegexp: "IkNlbGxzIg==",
		Color:     "#ffffff",
	},
}}
