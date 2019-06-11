package parser

import (
	"github.com/google/gopacket"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
	"sync"
	"time"
)

const WAIT_TIMEOUT = 5

type Parser struct {
	Source *gopacket.PacketSource
	DBClient *mongo.Client
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
	LastUpdate     time.Time
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

var Labels = []Label{
	{
		Name:      "test label",
		Type:      PacketOUT,
		Regexp:    regexp.MustCompile("Cells"),
		RawRegexp: "IkNlbGxzIg==",
		Color:     "#ffffff",
	},
}