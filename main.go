package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
)

func main() {
	if err, kit := loadKit(); err != nil {
		panic(fmt.Errorf("%w", errors.WithStack(err)))
	}

	if err := generateQR("qr.png", "http://localhost:3249/kit.qr"); err != nil {
		panic(fmt.Errorf("%w", errors.WithStack(err)))
	}

}

func generateQR(filename, url string) error {
	qrCode, _ := qrcode.New(url, qrcode.Medium)
	err := qrCode.WriteFile(256, filename)
	if err != nil {
		return fmt.Errorf("%w", errors.WithStack(err))
	}
	return nil
}

type Node struct {
	X       bool    `@"?"?`
	Y       bool    `@"?"?`
	Z       bool    `@"?"?`
	Gravity bool    `| @("true" | "false")`
	Nodes   []*Node `@@*`
}

func loadKit() (Node, error) {

	return {false, false, false, true}
}
