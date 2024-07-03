package main

import (
	"bytes"
	"encoding/binary"
)

type Node struct {
	Type             uint32
	NodeNameOffset   uint32
	ChildNodeOffset  uint32
	NextNodeOffset   uint32
	DataBoundsOffset uint32
}

type DataBounds struct {
	Min uint32 `json:"min"`
	Max uint32 `json:"max"`
}

func (n *Node) GetName() string {
	var name []byte
	for i := 0; data[32+int(n.NodeNameOffset)+i] != 0; i++ {
		name = append(name, data[32+int(n.NodeNameOffset)+i])
	}

	return string(name)
}

func (n *Node) GetType() uint32 {
	return n.Type
}

func (n *Node) GetChildNode() Node {
	var node Node
	err := binary.Read(bytes.NewReader(data[32+n.ChildNodeOffset:]), binary.BigEndian, &node)
	if err != nil {
		panic(err)
	}

	return node
}

func (n *Node) GetNextNode() Node {
	var node Node
	err := binary.Read(bytes.NewReader(data[32+n.NextNodeOffset:]), binary.BigEndian, &node)
	if err != nil {
		panic(err)
	}

	return node
}

func (n *Node) GetDataBounds() DataBounds {
	var bounds DataBounds
	err := binary.Read(bytes.NewReader(data[32+n.DataBoundsOffset:]), binary.BigEndian, &bounds)
	if err != nil {
		panic(err)
	}

	return bounds
}

func (n *Node) ParseChildren() map[string]any {
	node := n.GetChildNode()
	children := make(map[string]any)
	for {
		name := node.GetName()
		bounds := node.GetDataBounds()
		_type := node.GetType()

		if _type == 0 {
			children[name] = node.ParseChildren()
		} else if _type == 4 {
			// Boolean value
			children[name] = "Boolean"
		} else if _type == 8 {
			// Special type. Bounds are not defined in the binary, rather in the executable DOL.
			children[name] = "No bounds in binary."
		} else {
			children[name] = bounds
		}

		if node.NextNodeOffset == 0 {
			break
		}

		// Get next node
		node = node.GetNextNode()
	}

	return children
}
