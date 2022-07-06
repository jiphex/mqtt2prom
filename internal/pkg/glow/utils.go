package glow

import (
	"regexp"
)

type PacketType string

const (
	PT_STATE  = "STATE"
	PT_SENSOR = "SENSOR"
)

var TopicPattern = regexp.MustCompile(`(?m)^glow/(?P<topic>[[:xdigit:]]{12})/(?P<type>SENSOR|STATE)(?:/(?P<subtype>[^\/]+))?$`)
