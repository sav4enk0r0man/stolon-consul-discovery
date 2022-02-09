package api

import (
	"fmt"
	// "github.com/sav4enk0r0man/stolon-consul-discovery/logger"
	"strconv"
	"strings"
)

type Node struct {
	uid        string
	name       string
	role       string
	address    string
	healthy    bool
	serialized string
}

func NewNode(uid, name, role, address string, healthy bool) *Node {
	return &Node{
		uid:     uid,
		name:    name,
		role:    role,
		healthy: healthy,
		address: address,
	}
}

func (n *Node) Serialize() {
	n.serialized = fmt.Sprintf("[%s,%s]", n.Name(),
		strings.Join([]string{n.role, n.address,
			strconv.FormatBool(n.healthy), n.uid}, ","))
}

func (n Node) Uid() string {
	return n.uid
}

func (n Node) Name() string {
	return n.name
}

func (n Node) Role() string {
	return n.role
}

func (n Node) Address() string {
	return n.address
}

func (n Node) Healthy() bool {
	return n.healthy
}

func (n Node) Serialized() string {
	return n.serialized
}
