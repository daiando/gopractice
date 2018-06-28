package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const circle = `
{
	"id": "001",
	"type": "circle",
	"shape": {
		"radius": 5
	}
}
`

const rectangle = `
{
	"id": "002",
	"type": "rectangle",
	"shape": {
		"height": 5,
		"width": 2
	}
}
`

type Drawer interface {
	Draw() string
}

type Figure struct {
	Id    string
	Type  string
	Shape Drawer
}

type Circle struct {
	Radius int
}

type Rectangle struct {
	Height int
	Width  int
}

func main() {
	var c, r Figure
	if err := json.Unmarshal([]byte(circle), &c); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("circle %s\n", c.Shape.Draw())

	if err := json.Unmarshal([]byte(rectangle), &r); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rectangle %s\n", r.Shape.Draw())
}

func (c *Circle) Draw() string {
	msg := fmt.Sprintf("draw with radius %d", c.Radius)
	return msg
}

func (r *Rectangle) Draw() string {
	msg := fmt.Sprintf("draw with height %d, width %d", r.Height, r.Width)
	return msg
}

func (f *Figure) UnmarshalJSON(data []byte) error {
	type Alias Figure
	a := &struct {
		Shape json.RawMessage
		*Alias
	}{
		Alias: (*Alias)(f),
	}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	switch f.Type {
	case "circle":
		var c Circle
		if err := json.Unmarshal(a.Shape, &c); err != nil {
			return err
		}
		f.Shape = &c
	case "rectangle":
		var r Rectangle
		if err := json.Unmarshal(a.Shape, &r); err != nil {
			return err
		}
		f.Shape = &r
	default:
		return fmt.Errorf("unknown type: %q", f.Type)
	}
	return nil
}
