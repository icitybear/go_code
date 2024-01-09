package my_map

import (
	"fmt"
	"sync"
	"testing"
)

type Student struct {
	Age  int
	Name string
}

func TestMap(t *testing.T) {
	var m sync.Map
	s := Student{
		Age:  1,
		Name: "a",
	}
	//m.Store(s, s.Name)
	val, ok := m.Load(s)
	fmt.Println(val, ok)
	s1 := Student{
		Age:  2,
		Name: "b",
	}
	m.Store(s1, s1.Name)
	val, ok = m.Load(s1)
	fmt.Println(val, ok)
	// delete
}
