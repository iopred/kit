package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
)

func main() {
	var kit Node
	if kit, err := loadKit(); err == nil {
		if err := generateQR("qr.kit.png", "https://kit.iop.red"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.k.png", "kit.iop.red/k"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.i.png", "kit.iop.red/i"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.t.png", "kit.iop.red/t"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.21.png", "kit.iop.red/21"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.81.21.png", "kit.iop.red/81.21"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.r.png", "kit.iop.red/r"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.g.png", "https://qr.kit.iop.red./"+kit.now()); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.t.png", "https://qr.kit.iop.red./"+kit.next()); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.kit.iop.red.png", "https://kit.iop.red:3242/qr.kit.iop.red"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.description.kit.png", "there is a yellow smiley face with a big smile on it"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.ufo.naa.mba.png", "https://ufo.naa.mba"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.naa.mba.png", "https://naa.mba"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}

		if err := generateQR("qr.the.keeper.png", "https://www.keeperproject.com.au/"); err != nil {
			fmt.Errorf("%w", errors.WithStack(err))
		}
	} else {
		panic(fmt.Errorf("%w", errors.WithStack(err)))
	}

	port := 3242
	host := "naa.mba"
	filename := "ufo.png"

	type templateData struct {
		Host     string
		Port     int
		Filename string
	}

	generateQR("localhost.png", fmt.Sprintf("http://localhost:%d/", port))

	if err := generateQR("localhost.qr.kit.png", fmt.Sprintf("http://localhost:%d/qr.kit", port)); err != nil {
		fmt.Errorf("%w", errors.WithStack(err))
	}

	kitHandler := func(w http.ResponseWriter, r *http.Request) {
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
				"r g b t"
				"a k k k"
				"d k k k"
				"d k k k";
			place-self: center;
			// background-image: url('kit.png');
			// background-repeat: no-repeat;
		}
		div > * {
			mix-blend-mode: multiply;
		}
		span,
		iframe {
			border: 0px;
		}
		#k {
			grid-area: k;
		}
	</style>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
<div> <!- tl, br ->
	<img id="r" src="http://{{.Host}}:{{.Port}}/r.png">
	<span>// three.js</span>
	<img id="k" src='http://{{.Host}}:{{.Port}}/qr.kit.iop.red.png'/>
</div>
</body>
</html>`

		t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "kit"}}%s{{end}}`, html))
		if err != nil {
			panic("undefined")
		}

		err = t.ExecuteTemplate(w, "kit", templateData{Host: host, Port: port, Filename: filename})
		if err != nil {
			panic("undefined")
		}

	}

	htmlHandler := func(w http.ResponseWriter, r *http.Request, filename string) {
		fmt.Println("html", filename)

		if filename == "qr.kit.iop.red" {
			kitHandler(w, r)
			return
		}

		type templateData struct {
			Host     string
			Port     int
			Filename string
		}

		htmlt := `
<html>
<!- kit ->
<head>
	<style>
		* {
			margin: 0;
			padding: 0;
		}
		div {
			display: inline-grid;
			grid-template-areas:
				"r g b t"
				"a k k k"
				"d k k k"
				"d k k k";
			place-self: center;
			// background-image: url('kit.png');
			// background-repeat: no-repeat;
		}
		iframe {
			position: absolute;
			width: 100%;
			height: 100%;
			border: none;
		}

		img {
			position: absolute;
			right: 10px;
			bottom: 10px;
		}
	</style>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
<div> <!-- tl, br -->
	<iframe src='{{.Filename}}'/>
	<img src='/qr.{{.Filename}}.png'/>
</div>
</body>
</html>`

		t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "kit"}}%s{{end}}`, htmlt))
		if err != nil {
			panic("undefined")
		}
		err = t.ExecuteTemplate(w, "kit", templateData{Host: host, Port: port, Filename: filename})
		if err != nil {
			panic("undefined")
		}

	}

	pngHandler := func(w http.ResponseWriter, r *http.Request, filename string) {
		filename = strings.trim(filename, ".png")
		fmt.Println("png", filename)
		qrCode, _ := qrcode.New(filename, qrcode.Low)
		if err := qrCode.Write(11, w); err != nil {
			fmt.Errorf("error creating qr code: %w\n", errors.WithStack(err))
		}
	}

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			paths := strings.Split(path, "/")
			path = paths[len(paths)-1]

			segments := strings.Split(path, ".")

			filetype := segments[len(segments)-1]
			filename := ""
			if len(segments) == 1 {
				filename = filetype
				filetype = ""
			} else {
				filename = strings.Join(segments[:len(segments)-1], ".")
			}

			if filename == "" && filetype == "" {
				filename = "kit.iop.red"
			}

			switch filetype {
			case "ico":
				fallthrough
			case "png":
				if strings.HasPrefix(filename, "qr.") {
					pngHandler(w, r, "https://naa.mba/"+strings.Trim(filename, "qr.")+".png")
				} else {
					http.ServeFile(w, r, filename+".png")
				}
			case "html":
			default:
				htmlHandler(w, r, filename)
			}
		},
	)

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
