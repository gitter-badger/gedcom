package gedcom

import (
	"io"
	"bufio"
	"strconv"
	"regexp"
)

// Decoder represents a GEDCOM decoder.
type Decoder struct {
	r *bufio.Reader
}

// Create a new decoder to parse a reader that contain GEDCOM data.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: bufio.NewReader(r),
	}
}

// Decode will parse the entire GEDCOM stream (until EOF is reached) and return
// a DocumentNode. If the GEDCOM stream is not valid then the document node will
// be nil and the error is returned.
//
// A blank GEDCOM or a GEDCOM that only contains empty lines is valid and a
// DocumentNode will be returned with zero nodes.
func (dec *Decoder) Decode() (*DocumentNode, error) {
	documentNode := &DocumentNode{
		Nodes: []*SimpleNode{},
	}
	indents := []*SimpleNode{}

	finished := false
	for !finished {
		line, err := dec.r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				finished = true
			} else {
				return nil, err
			}
		}

		node := parseLine(line)

		// Skip blank lines.
		if node.Tag == "" {
			continue
		}

		// Add a root node to the document.
		if node.Indent == 0 {
			documentNode.Nodes = append(documentNode.Nodes, node)
			indents = append(indents, node)
			continue
		}

		i := indents[node.Indent-1]

		// Move indent pointer if we are changing depth.
		switch {
		case node.Indent >= len(indents):
			indents = append(indents, node)

		case node.Indent < len(indents)-1:
			indents = indents[:len(indents)-1]
		}

		i.Children = append(i.Children, node)
	}

	return documentNode, nil
}

func parseLine(line string) *SimpleNode {
	parts := regexp.
		MustCompile(`^(\d) (@\w+@ )?(\w+)( .+)?\n?$`).
		FindStringSubmatch(line)

	indent := 0
	if len(parts) > 1 {
		indent, _ = strconv.Atoi(parts[1])
	}

	pointer := ""
	if len(parts) > 2 && len(parts[2]) > 4 {
		pointer = parts[2][1 : len(parts[2])-2]
	}

	tag := ""
	if len(parts) > 3 {
		tag = parts[3]
	}

	value := ""
	if len(parts) > 4 && len(parts[4]) > 0 {
		value = parts[4][1:]
	}

	return &SimpleNode{
		Indent:   indent,
		Tag:      tag,
		Value:    value,
		Pointer:  pointer,
		Children: []*SimpleNode{},
	}
}

// SimpleNode is used as the default node type when there is no more appropriate
// or specific type to use.
type SimpleNode struct {
	Indent   int
	Tag      string
	Value    string
	Pointer  string
	Children []*SimpleNode
}

// DocumentNode represents a whole GEDCOM document. It is possible for a
// DocumentNode to contain zero Nodes, this means the GEDCOM file was empty. It
// may also (and usually) contain several Nodes.
type DocumentNode struct {
	Nodes []*SimpleNode
}
