package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"strings"
)

var data []byte

type jsonEncoderCtx struct {
	data map[string]any
}

func main() {
	dir, err := os.ReadDir("xml")
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		data, err = os.ReadFile("xml/" + file.Name())
		if err != nil {
			panic(err)
		}

		// First is always the root node
		var node Node
		err = binary.Read(bytes.NewReader(data[36:]), binary.BigEndian, &node)
		if err != nil {
			panic(err)
		}

		// Get first child
		err = binary.Read(bytes.NewReader(data[32+node.ChildNodeOffset:]), binary.BigEndian, &node)
		if err != nil {
			panic(err)
		}

		encoder := jsonEncoderCtx{data: make(map[string]any)}
		for {
			name := node.GetName()
			bounds := node.GetDataBounds()
			_type := node.GetType()

			if _type == 0 {
				// This is a parent node.
				encoder.data[name] = node.ParseChildren()

				// It also includes how many repeating of itself can exist.
				encoder.data[name].(map[string]any)["bounds"] = bounds
			} else if _type == 4 {
				// Boolean value
				encoder.data[name] = "Boolean"
			} else if _type == 8 {
				// Special type. Bounds are not defined in the binary, rather in the executable DOL.
				encoder.data[name] = "No bounds in binary."
			} else {
				encoder.data[name] = bounds
			}

			if node.NextNodeOffset == 0 {
				break
			}

			// Get next node
			node = node.GetNextNode()
		}

		j, err := json.MarshalIndent(encoder.data, "", "\t")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("output/"+strings.Replace(file.Name(), "bin", "json", -1), j, 0644)
		if err != nil {
			panic(err)
		}
	}
}
