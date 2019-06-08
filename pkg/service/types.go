package service

import (
	"github.com/lyawaf/traffic-sniffer/pkg/parser"
)

type RawLabel struct {
	Name   string           `json:"name"`
	Color  string           `json:"color"`
	Type   parser.LabelType `json:"type"`
	Regexp string           `json:"regexp"`
}
