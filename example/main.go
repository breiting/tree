package main

import (
	"os"

	"github.com/breiting/tree"
)

func generateSample() {

	p := tree.Node{
		ID:   "foo",
		Name: "p",
		Data: "This is a simple HTML document.",
	}

	h1 := tree.Node{
		ID:   "bar2",
		Name: "h1",
		Data: 1.33,
		Attributes: map[string]string{
			"shape": "doublecircle",
			"color": "blue",
		},
	}

	h2 := tree.Node{
		ID:       "bar",
		Name:     "h2",
		Data:     1.33,
		Children: []*tree.Node{&h1},
	}

	html := tree.Node{
		ID:   "root",
		Name: "head",
		Attributes: map[string]string{
			"shape": "box",
			"color": "red",
		},
		Children: []*tree.Node{&p, &h2},
	}

	f, _ := os.Create("sample1.dot")
	defer f.Close()
	err := tree.WriteToDot(&html, f)
	if err != nil {
		panic(err)
	}
}

func deserialize() {

	var relations []tree.Relation
	relations = append(relations, tree.Relation{ID: "6", ParentID: "2"})
	relations = append(relations, tree.Relation{ID: "5", ParentID: "2"})
	relations = append(relations, tree.Relation{ID: "4", ParentID: "3"})
	relations = append(relations, tree.Relation{ID: "3", ParentID: "1"})
	relations = append(relations, tree.Relation{ID: "2", ParentID: "1"})
	relations = append(relations, tree.Relation{ID: "1"})

	t, err := tree.Deserialize(relations)
	f, _ := os.Create("sample2.dot")
	defer f.Close()
	err = tree.WriteToDot(t, f)
	if err != nil {
		panic(err)
	}
}

func main() {

	generateSample()
	deserialize()
}
