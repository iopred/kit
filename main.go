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
		kit.now();	
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
		URL	     string
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
<div> <!-- r, g, b -->
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

		err = t.ExecuteTemplate(w, "kit", templateData{Host: host, Port: port, Filename: filename, URL: filename})
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
			URL	     string
		}

		htmlt := `
<html>
<!-- kit -->
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
			right: 0px;
			bottom: 0px;
		}
	</style>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
<div> <!-- tl, br -->
	<iframe src='{{.URL}}'></iframe>
	<a href="{{.URL}}" target="_blank"><img src='/{{.Filename}}.png' id="qr"/></a>
</div>
</body>
</html>`

		t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "kit"}}%s{{end}}`, htmlt))
		if err != nil {
			panic("undefined")
		}

		url := filename
		switch filename {
		case "heliattack":
			url = "https://heliattack.com"
			filename = "ha2000"
		case "ha2000":
			url = "https://heliattack.com/game"
		case "ufo.naa.mba":
			url = "https://ufo.naa.mba"
		case "naa.mba":
			url = "https://naa.mba"
		case "ufo":
			url = "https://ufo.naa.mba"
		case "the.keeper":
			url = "https://www.keeperproject.com.au"
		case "bad.habit":
			url = "https://badhabitrecords.com.au"
		case "the.presynct":
			url = "https://www.thepresynct.com.au"
		case "naamba":
			url = "https://naa.mba"
		}

		err = t.ExecuteTemplate(w, "kit", templateData{Host: host, Port: port, Filename: filename, URL: url})
		if err != nil {
			panic("undefined")
		}

	}

	pngHandler := func(w http.ResponseWriter, r *http.Request, filename string) {
		fmt.Println("png", filename)

		filename = strings.Trim(filename, ".png")

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

	http.HandleFunc("/ufo.naa.mba", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "ufo.naa.mba")
	})

	http.HandleFunc("/naa.mba", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "naa.mba")
	})

	http.HandleFunc("/ufo", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "ufo")
	})

	http.HandleFunc("/kit.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "kit.png")
	})

	http.HandleFunc("/the.keeper", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "the.keeper")
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
