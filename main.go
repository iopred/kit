package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
)

func main() {
	var kit Node
	if kit, err := loadKit(); err == nil {
		if err := generateQR("kit.qr.png", "kit.iop.red"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("r.qr.png", "kit.iop.red/r"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("g.qr.png", "kit.iop.red/"+kit.now()+".qr"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("t.qr.png", "http://kit.iop.red/"+kit.next()+".qr"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("kit.iop.red.qr.png", "http://localhost:3242/kit.iop.red.qr"); err != nil {
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

			//kit := fmt.Sprintf("http://localhost:%d", port)

			html := `
		<html>
		<head>
			<style>
				* {
					margin: 0;
					padding: 0;
				}
				div {
					display: inline-grid;
					grid-template-areas:
						"r g b"
						"t k k"
						"d k k";
					place-self: center;
					background-image: src('kit.png');
					background-repeat: no-repeat;
				}
				div > img {
					mix-blend-mode: multiply;
				}
				div > iframe {
					border: 0px;
					grid-area: k;
					mix-blend-mode: multiply;
				}
			</style>
		</head>
		<body>
		<div> <!- tl, br ->
			<img src='http://localhost:{{.}}/kit.iop.red.png'/>
			<iframe src='http://localhost:{{.}}/kit.iop.red.qr'/>
		</div>
		</body>
		</html>`

			t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "kit"}}%s{{end}}`, html))
			if err != nil {
				panic("undefined")
			}
			err = t.ExecuteTemplate(w, "kit", port)
			if err != nil {
				panic("undefined")
			}

		},
	)

	http.HandleFunc(
		"/kit.iop.red.qr",
		func(w http.ResponseWriter, r *http.Request) {

			//kit := fmt.Sprintf("http://localhost:%d", port)

			html := `
		<html>
		<head>
			<style>
				* {
					margin: 0;
					padding: 0;
				}
				div {
					display: inline-grid;
					grid-template-rows: 33px 33px 33px;
					grid-template-columns: 33px 33px 33px;
					place-self: center;
				}
				div > * {
					mix-blend-mode: multiply;
				}
				#k {
					grid-column: 2;
					grid-row: 2;
					width: 66px;
					height: 66px;
				}
				
			</style>
		</head>
		<body>
		<div> <!- tl, br ->
			<img src='http://localhost:{{.}}/kit.iop.red.qr.png'/>
			<img id="k" src='http://localhost:{{.}}/kit.png'/>
			<img id="k" src='http://localhost:{{.}}/kit.iop.red.png'/>
		</div>
		</body>
		</html>`

			t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "kit"}}%s{{end}}`, html))
			if err != nil {
				panic("undefined")
			}
			err = t.ExecuteTemplate(w, "kit", port)
			if err != nil {
				panic("undefined")
			}

		},
	)

	http.HandleFunc("/qr.png", func(w http.ResponseWriter, r *http.Request) {
		qrCode, _ := qrcode.New("qr", qrcode.Low)
		if err := qrCode.Write(11, w); err != nil {
			fmt.Errorf("error creating qr code: %w\n", errors.WithStack(err))
		}
	})

	http.HandleFunc("/kit.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "kit.png")
	})

	http.HandleFunc("/kit.qr.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "kit.qr.png")
	})

	http.HandleFunc("/kit.iop.red.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "kit.iop.red.png")
	})

	http.HandleFunc("/kit.iop.red.qr.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "kit.iop.red.qr.png")
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
	err := qrCode.WriteFile(11, filename)
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

func (n Node) x() string {
	if n.X {
		return "1"
	}
	return "0"
}

func (n Node) y() string {
	if n.Y {
		return "1"
	}
	return "0"
}

func (n Node) z() string {
	if n.Z {
		return "1"
	}
	return "0"
}

func (n Node) g() string {
	if n.Gravity {
		return "1"
	}
	return "0"
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
