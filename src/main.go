package main

import (
	"fmt"
	"kit/kit"
)

func spectrum() kit.System {
	k := kit.New()
	n := kit.Node{
		T:      2022,
		Dvalue: 7950000000,
		Name:   "people",
	}
	k.AddNode(n)

	start := k.Next(0).Node("people")

	if start.Value < 0 {
		panic("no people")
	}

	k.AddNode(kit.Node{
		Dt:     1,
		Dvalue: 1,
		Name:   "spectrum",
		Origin: start,
	})

	return k
}

func tunnel(s kit.System) (kit.System, error) {
	start := s.Next(0).Node("people")
	end := s.Next(1).Node("people")

	if start.Value < 0 || end.Value < 0 {
		panic("no people")
	}

	k := s.Next(1)
	k.AddNode(kit.Node{
		Dt:     1,
		Name:   "input",
		Origin: start,
	})
	k.AddNode(kit.Node{
		T:      1,
		Dt:     -1,
		Name:   "output",
		Origin: end,
	})
	return k, nil
}

func main() {
	k := spectrum()
	fmt.Println("spectrum")
	fmt.Println(k.Next(0).Node("spectrum"))
	fmt.Println(k.Next(0.25).Node("spectrum"))
	fmt.Println(k.Next(0.5).Node("spectrum"))
	fmt.Println(k.Next(0.75).Node("spectrum"))
	fmt.Println(k.Next(1).Node("spectrum"))

	t, err := tunnel(k)
	if err != nil {
		return
	}
	fmt.Println("tunnel")
	fmt.Println(t.Next(0).Node("output"))
	fmt.Println(t.Next(0.5).Node("output"))
	fmt.Println(t.Next(1).Node("output"))

	// fmt.Println(kit(system()))
}
