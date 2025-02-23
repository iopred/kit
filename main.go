package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
)

var emojiMetadata map[string]string

func init() {
	emojiMetadata = make(map[string]string)
	err := filepath.Walk("./emoji", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "metadata.json" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var metadata struct {
				Glyph            string   `json:"glyph"`
				CLDR             string   `json:"cldr"`
				UnicodeSkinTones []string `json:"unicodeSkintones"`
			}
			if err := json.Unmarshal(data, &metadata); err != nil {
				return err
			}

			metadata.CLDR = strings.ReplaceAll(metadata.CLDR, " ", "_")

			if len(metadata.UnicodeSkinTones) > 0 {
				emojiMetadata[metadata.Glyph] = filepath.Join(filepath.Dir(path), "Default/3D/", metadata.CLDR+"_3d_default.png")
			} else {
				emojiMetadata[metadata.Glyph] = filepath.Join(filepath.Dir(path), "3D/", metadata.CLDR+"_3d.png")
			}
		}
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("error loading emoji metadata: %w", errors.WithStack(err)))
	}
}

func emojiHandler(w http.ResponseWriter, r *http.Request, emoji string) {
	// Maybe check that the emoji is a rune?

	if path, ok := emojiMetadata[emoji]; ok {
		fmt.Println("serving ", emoji, " :", path)
		http.ServeFile(w, r, path)
	} else {
		http.Error(w, "Emoji not found", http.StatusNotFound)
	}
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %w", errors.WithStack(err)))
	}

	if kit, err := loadKit(); err == nil {
		fmt.Println(kit.now(), "->", kit.next())
	} else {
		panic(fmt.Errorf("error loading kit: %w", errors.WithStack(err)))
	}

	portStr := os.Getenv("QR_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Errorf("invalid port: %w", errors.WithStack(err)))
	}

	host := "naa.mba"
	filename := "ufo.png"

	type templateData struct {
		Host     string
		Port     int
		Filename string
		URL      string
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
			URL      string
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
		case "heliattack2000":
			fallthrough
		case "ha2k":
			fallthrough
		case "ha2000":
			url = "https://ha2k.heliattack.com"
			filename = "ha2000"
		case "ufo":
			fallthrough
		case "ufo.naa.mba":
			url = "https://ufo.naa.mba"
		case "naamba":
			fallthrough
		case "naa.mba":
			url = "https://naa.mba"
		case "the.keeper":
			url = "https://www.keeperproject.com.au"
		case "bad.habit":
			url = "https://badhabitrecords.com.au"
			filename = "qr.bad.habit"
		case "the.presynct":
			url = "https://www.thepresynct.com.au"
			filename = "qr.the.presynct"
		case "🧑":
			url = "https://kit.iop.red/emoji/🧑"
			filename = "https://kit.iop.red/qr.🧑.png"
		case "🚁":
			url = "https://heliattack.com"
			filename = "https://kit.iop.red/qr.🚁.png"
		case "📼":
			url = "https://ha2k.heliattack.com"
			filename = "https://kit.iop.red/qr.📼.png"
		}

		err = t.ExecuteTemplate(w, "kit", templateData{Host: host, Port: port, Filename: filename, URL: url})
		if err != nil {
			panic("undefined")
		}

	}

	pngHandler := func(w http.ResponseWriter, r *http.Request, filename string) {
		fmt.Println("png", filename)

		filename = strings.Trim(filename, ".png")

		if filename == "📼" {
			filename = "tape"
		}

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

	http.HandleFunc("/bad.habit", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "bad.habit")
	})

	http.HandleFunc("/the.presynct", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "the.presynct")
	})

	http.HandleFunc("/🧑", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🧑")
	})

	http.HandleFunc("/🚁", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🚁")
	})

	http.HandleFunc("/📼", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "📼")
	})

	http.HandleFunc("/🖭", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "📼")
	})

	http.HandleFunc("/emoji/", func(w http.ResponseWriter, r *http.Request) {
		emojiHandler(w, r, r.URL.Path[len("/emoji/"):])
	})

	fmt.Println("Listening on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		panic("server closed")
	} else if err != nil {
		fmt.Println("error starting server: %w", errors.WithStack(err))
		os.Exit(1)
	}

	fmt.Println("exiting qr.kit")
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
