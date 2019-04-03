package main

import (
	"algorithm/set"
	"fmt"
)

func main() {
	var s = set.New()

	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	s.Add(5)
	s.Add(6)

	fmt.Println(s.Size())
	s.Add(1)
	s.Remove(6)

	fmt.Println(s.Size())
	fmt.Println(s.Contains(3))
	fmt.Println(s.Contains(6))
	fmt.Println(s.Contains(7))
	fmt.Println(s.Values())

}
