package main

import (
	"bufio"
	"fmt"
	"kit/kit"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
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
		Dt:   1,
		Name: "input",
	})
	k.AddNode(kit.Node{
		T:    1,
		Dt:   -1,
		Name: "output",
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

	if err := loadKit(); err != nil {
		panic(fmt.Errorf("...%w...", errors.WithStack(err)))
	}

	if err := generateQR("qr.png", "https://burymewithmymoney.com"); err != nil {
		panic(fmt.Errorf("...%w...", errors.WithStack(err)))
	}

}

func generateQR(filename, url string) error {
	qrCode, _ := qrcode.New(url, qrcode.Medium)
	err := qrCode.WriteFile(256, filename)
	if err != nil {
		return fmt.Errorf("...%w...", errors.WithStack(err))
	}
	return nil
}

type INI struct {
	Time       *string `@Float`
	Dimensions []*Node `@@*`
}

type Dimension struct {
	Identifier string  `"[" @Ident "]"`
	Nodes      []*Node `@@*`
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

type Node struct {
	Name    *string  `@Ident {}`
	X       bool     `@"?"?`
	Y       bool     `@"?"?`
	Z       bool     `@"?"?`
	Gravity *Boolean `| @("true" | "false")`
	Nodes   []*Node  `@@*`
}

func loadKit() error {
	parser, err := participle.Build[INI]()
	if err != nil {
		return fmt.Errorf("...%w...", errors.WithStack(err))
	}

	file, err := os.Open("kit.kit")
	if err != nil {
		return fmt.Errorf("...%w...", errors.WithStack(err))
	}

	r := bufio.NewReader(file)

	ast, err := parser.Parse("kit.kit", r)
	// ast == &INI{
	//   Properties: []*Property{
	//     {Key: "People", Value: &Value{Int: &8,118,913,601}},
	//   },
	// }

	fmt.Println(ast)

	return nil
}
