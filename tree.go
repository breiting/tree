// Package tree is a very simple package to work with trees. You can easily deserialize a parent
// to child relationship and also dump the result into a .dot file (graphviz).
package tree

import (
	"fmt"
	"io"
)

const (
	defaultColor  = "lightblue"
	defaultShape  = "ellipse"
	maxIterations = 10000
)

// Node provides a hierarchical tree based data structure with optional payload per node
type Node struct {
	ID       string
	Name     string
	Children []*Node

	// Data is the attached data for every node, can be anything, depending on the use case
	Data interface{}

	// Attributes can carry any additional information such as shape and color
	Attributes map[string]string
}

// Relation is a simple relationship between a node and a parent, which can be used for a seralized
// representation of the tree itself
type Relation struct {
	ID       string
	ParentID string
}

// NewNode creates a new root node of a tree (no parent)
func NewNode(id string) *Node {
	return &Node{
		ID:         id,
		Attributes: make(map[string]string),
	}
}

// Deserialize takes the relations and build a tree structure, make sure that you provide one and
// only one root node with no parent ID
func Deserialize(relations []Relation) (*Node, error) {

	var root *Node
	assigned := make(map[string]bool) // denotes if the node is already in the tree

	// find root, only allow one root node
	for _, v := range relations {
		assigned[v.ID] = false

		if v.ParentID == "" {
			// only allow one root
			if root != nil {
				return nil, fmt.Errorf("Multiple roots have been found")
			}
			root = NewNode(v.ID)
			assigned[v.ID] = true
		}
	}

	if root == nil {
		return nil, fmt.Errorf("No root was found")
	}

	iteration := 0
	for {
		iteration++
		if checkAllAssigned(assigned) || iteration > maxIterations {
			break // we are done
		}
		for i := 0; i < len(relations); i++ {
			if assigned[relations[i].ID] {
				continue
			}
			if parent := FindByID(root, relations[i].ParentID); parent != nil {
				// we have found the parent in the tree
				parent.Children = append(parent.Children, NewNode(relations[i].ID))
				assigned[relations[i].ID] = true
			}
		}
	}

	if iteration > maxIterations {
		return root, fmt.Errorf("Max iteration reached, some nodes could not be assigned")
	}

	return root, nil
}

// FindByID searches in BFS manner
func FindByID(root *Node, id string) *Node {

	queue := make([]*Node, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.ID == id {
			return nextUp
		}
		if len(nextUp.Children) > 0 {
			for _, child := range nextUp.Children {
				queue = append(queue, child)
			}
		}
	}
	return nil
}

// FindByIDDFS searchs in depth first search
func FindByIDDFS(node *Node, id string) *Node {
	if node.ID == id {
		return node
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			FindByIDDFS(child, id)
		}
	}
	return nil
}

// WriteToDot writes the whole tree into a graphviz-dot structure
func WriteToDot(node *Node, w io.Writer) error {

	fmt.Fprintf(w, "digraph {\n")

	writeNodeShape(node, w)
	writeNodeRelationship(node, w)

	fmt.Fprintln(w, "}")
	return nil
}

// String dumps a tree
func (t Node) String() string {

	if &t == nil {
		return "()"
	}
	s := ""
	for _, v := range t.Children {
		s += v.String() + " "
	}
	s += fmt.Sprint(t.Name)
	return "(" + s + ")"
}

func writeNodeRelationship(n *Node, w io.Writer) {

	if n == nil {
		return
	}
	for _, v := range n.Children {
		fmt.Fprintf(w, "  \"%s\" -> \"%s\";\n", n.ID, v.ID)
		writeNodeRelationship(v, w)
	}
}

func writeNodeShape(n *Node, w io.Writer) {

	if n == nil {
		return
	}
	color := defaultColor
	if n.Attributes["color"] != "" {
		color = n.Attributes["color"]
	}
	shape := defaultShape
	if n.Attributes["shape"] != "" {
		shape = n.Attributes["shape"]
	}
	fmt.Fprintf(w, "  \"%s\" [shape=%s,style=filled,color=%s]\n", n.ID, shape, color)
	for _, v := range n.Children {
		writeNodeShape(v, w)
	}
}

func checkAllAssigned(assigned map[string]bool) bool {

	for _, v := range assigned {
		if v == false {
			return false
		}
	}
	return true
}
