package tree

import (
	"os"
	"testing"
)

func TestCreation(t *testing.T) {

	p := Node{
		ID:   "foo",
		Name: "p",
		Data: "This is a simple HTML document.",
	}

	h1 := Node{
		ID:   "bar",
		Name: "h1",
		Data: 1.33,
	}

	h2 := Node{
		ID:       "bar",
		Name:     "h2",
		Data:     1.33,
		Children: []*Node{&h1},
	}

	html := Node{
		ID:       "root",
		Name:     "head",
		Children: []*Node{&p, &h2},
	}

	err := WriteToDot(&html, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeserialize(t *testing.T) {

	var relations []Relation
	relations = append(relations, Relation{ID: "6", ParentID: "2"})
	relations = append(relations, Relation{ID: "5", ParentID: "2"})
	relations = append(relations, Relation{ID: "4", ParentID: "3"})
	relations = append(relations, Relation{ID: "3", ParentID: "1"})
	relations = append(relations, Relation{ID: "2", ParentID: "1"})
	relations = append(relations, Relation{ID: "1"})

	tree, err := Deserialize(relations)
	if err != nil {
		t.Fatal(err)
	}

	if tree.Children[0].ID != "3" {
		t.Fatalf("Name does not match expected %s got %s", "3", tree.Children[0].ID)
	}
}

func ExampleDeserialize() {

	var relations []Relation
	relations = append(relations, Relation{ID: "6", ParentID: "2"})
	relations = append(relations, Relation{ID: "5", ParentID: "2"})
	relations = append(relations, Relation{ID: "4", ParentID: "3"})
	relations = append(relations, Relation{ID: "3", ParentID: "1"})
	relations = append(relations, Relation{ID: "2", ParentID: "1"})
	relations = append(relations, Relation{ID: "1"})

	t, err := Deserialize(relations)
	f, _ := os.Create("sample.dot")
	defer f.Close()
	err = WriteToDot(t, f)
	if err != nil {
		panic(err)
	}
}
