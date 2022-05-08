package main

import (
	"fmt"
	"kit/kit"
)

func unsafespectrum() kit.System {
	k := kit.New()
	k.AddNode(kit.Node{
		Dt:     1,
		Value:  0,
		Dvalue: 1,
		Name:   "spectrum",
	})
	return k
}

func onoff() kit.System {
	k := kit.New()
	k.AddNode(kit.Node{
		X:      0,
		Dx:     1,
		T:      0,
		Dt:     1,
		Value:  0,
		Dvalue: 1,
		Name:   "input",
	})
	k.AddNode(kit.Node{
		X:      1,
		Dx:     -1,
		T:      1,
		Dt:     -1,
		Value:  1,
		Dvalue: 1,
		Name:   "output",
	})
	return k
}

func main() {
	k := unsafespectrum()
	fmt.Println("spectrum")
	fmt.Println(k.Next(0).Node("spectrum"))
	fmt.Println(k.Next(0.5).Node("spectrum"))
	fmt.Println(k.Next(1).Node("spectrum"))

	of := onoff()
	fmt.Println("On off")
	fmt.Println(of.Next(0).Node("output"))
	fmt.Println(of.Next(0.5).Node("output"))
	fmt.Println(of.Next(1).Node("output"))

	// fmt.Println(kit(system()))
}
