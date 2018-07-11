package gedcom

import "bytes"

// Document represents a whole GEDCOM document. It is possible for a
// Document to contain zero Nodes, this means the GEDCOM file was empty. It
// may also (and usually) contain several Nodes.
type Document struct {
	Nodes []Node
}

// String will render the entire GEDCOM document.
func (doc *Document) String() string {
	buf := bytes.NewBufferString("")

	encoder := NewEncoder(buf, doc)
	err := encoder.Encode()
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func (doc *Document) Individuals() []*IndividualNode {
	individuals := []*IndividualNode{}

	for _, node := range doc.Nodes {
		if n, ok := node.(*IndividualNode); ok {
			individuals = append(individuals, n)
		}
	}

	return individuals
}

func (doc *Document) NodeByPointer(ptr string) Node {
	for _, node := range doc.Nodes {
		if node.Pointer() == ptr {
			return node
		}
	}

	return nil
}

func (doc *Document) Families() []*FamilyNode {
	families := []*FamilyNode{}

	for _, node := range doc.Nodes {
		if n, ok := node.(*FamilyNode); ok {
			families = append(families, n)
		}
	}

	return families
}
