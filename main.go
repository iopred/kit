package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
)

func main() {
	var kit Node
	if kit, err := loadKit(); err == nil {
		if err := generateQR("kit.png", "http://kit.iop.red/kit.qr"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("0001.png", "http://kit.iop.red/0001.qr"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("0001.png", "http://kit.iop.red/3249.qr"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("now.png", "http://kit.iop.red/"+kit.now()+".qr"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("g.png", "http://kit.iop.red/g.qr"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}
	} else {
		panic(fmt.Errorf("%w", errors.WithStack(err)))
	}

	port := 3242

	generateQR("localhost.png", fmt.Sprintf("http://localhost:%d/", port))

	if err := generateQR("localhost.kit.png", fmt.Sprintf("http://localhost:%d/kit.qr", port)); err != nil {
		fmt.Errorf("%w", errors.WithStack(err))
	}

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `
		<html>
		<head>
			<style>
				* {
					margin: 0;
					padding: 0;
				}
				div {
					display: grid;
					height: 100%;
				}
				.c {
					max-width: 100%;
					max-height: 100vh;
					margin: auto;
				}
			</style>
		</head>
		<body>
		<div>
			<img class='c' src='http://localhost:3242/qr.png'>
		</div>
		</body>
		</html>`, port)
		},
	)

	http.HandleFunc("/qr.png", func(w http.ResponseWriter, r *http.Request) {
		qrCode, _ := qrcode.New(kit.now(), qrcode.Medium)
		if err := qrCode.Write(11, w); err != nil {
			fmt.Errorf("error creating qr code: %w\n", errors.WithStack(err))
		}
	})

	fmt.Println(kit.now())

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Errorf("server closed\n")
	} else if err != nil {
		fmt.Errorf("error starting server: %w\n", errors.WithStack(err))
		os.Exit(1)
	}

	fmt.Println("exiting")
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
	X       bool
	Y       bool
	Z       bool
	Gravity bool
	Nodes   byte
}

func (n Node) now() string {
	return "0009"
}

func (n Node) next() string {
	return "🌞"
}

func loadKit() (Node, error) {
	return Node{false, false, false, true, 0001}, nil
}
