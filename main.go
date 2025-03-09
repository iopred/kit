package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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

			metadata.CLDR = strings.ToLower(strings.ReplaceAll(metadata.CLDR, " ", "_"))

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
	// Create a context that listens for interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Received interrupt signal, shutting down...")
		cancel() // Cancel the context to signal shutdown
	}()

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

	// Check for command-line arguments
	if len(os.Args) > 1 {
		// todo merge with the emoji image.
		emoji := string(os.Args[1][0])           // Get the first character of the first argument
		filename := fmt.Sprintf("%s.png", emoji) // Create the filename
		if err := generateQR(filename, emoji); err != nil {
			fmt.Printf("Error generating QR code: %v\n", err)
			os.Exit(1)
		}
		return // Exit after generating the QR code
	}

	port := 3242
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	if portStr := os.Getenv("PORT"); portStr != "" {
		if portParsed, err := strconv.Atoi(portStr); err == nil {
			port = portParsed
		}
	}

	type templateData struct {
		Host     string
		Port     int
		Filename string
		URL      string
	}

	htmlHandler := func(w http.ResponseWriter, r *http.Request, filename string) {
		fmt.Println("html", filename)

		type templateData struct {
			Host     string
			Port     int
			Filename string
			URL      string
		}

		htmlt := `
<!DOCTYPE html>
<html>
<!-- kit -->
<head>
	<style>
		* {
			margin: 0;
			padding: 0;
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
<iframe src='{{.URL}}'></iframe>
<img src='/{{.Filename}}' id="qr" onclick="hideElement(this)" style="cursor: pointer;"/>

<script>
	function hideElement(element) {
		element.style.display = 'none';
	}
</script>
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
			filename = "https://kit.iop.red/qr.🛸.png"
			url = "https://ufo.naa.mba"
		case "naamba":
			fallthrough
		case "naa.mba":
			filename = "naamba.png"
			url = "https://naa.mba"
		case "the.keeper":
			url = "https://www.keeperproject.com.au"
			filename = "https://kit.iop.red/the.keeper.png"
		case "bad.habit":
			url = "https://badhabitrecords.com.au"
			filename = "https://kit.iop.red/qr.bad.habit.png"
		case "the.presynct":
			url = "https://www.thepresynct.com.au"
			filename = "https://kit.iop.red/qr.the.presynct.png"
		case "🧑":
			url = "https://kit.iop.red/emoji/🧑"
			filename = "https://kit.iop.red/qr.🧑.png"
		case "🏡":
			url = "https://kit.iop.red/emoji/🏡"
			filename = "https://kit.iop.red/qr.🏡.png"
		case "🌞":
			url = "https://kit.iop.red/emoji/🌞"
			filename = "https://kit.iop.red/qr.🌞.png"
		case "🌝":
			url = "https://kit.iop.red/emoji/🌝"
			filename = "https://kit.iop.red/qr.🌝.png"
		case "🌚":
			url = "https://kit.iop.red/emoji/🌚"
			filename = "https://kit.iop.red/qr.🌚.png"
		case "🌛": // oops, mirrored (:P), this is good information in a fight :)
			url = "https://kit.iop.red/emoji/🌜"
			filename = "https://kit.iop.red/qr.🌜.png"
		case "🌜":
			url = "https://kit.iop.red/emoji/🌛"
			filename = "https://kit.iop.red/qr.🌛.png"
		case "📼":
			url = "https://kit.iop.red/tape.png"
			filename = "https://kit.iop.red/qr.📼.png"
		case "🚁":
			url = "https://🚁.heliattack.com"
			filename = "https://kit.iop.red/qr.🚁.png"
		case "🚁🧑":
			url = "https://heliattack.com/🚁🧑.png"
			filename = "https://kit.iop.red/qr.🚁🧑.png"
		case "🚁🪖":
			url = "https://heliattack.com/🚁🪖"
			filename = "https://kit.iop.red/qr.🚁🪖.png"
		case "🚁🪖🧑":
			url = "https://heliattack.com/🚁🪖🧑"
			filename = "https://kit.iop.red/qr.🚁🪖🧑.png"
		case "🚁🔫":
			url = "https://heliattack.com/🚁🔫";
			filename = "https://kit.iop.red/qr.🚁🔫.png";
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
				filename = "🌞"
			}

			switch filetype {
			case "ico":
				fallthrough
			case "png":
				if strings.HasPrefix(filename, "qr.") {
					pngHandler(w, r, "https://naa.mba/"+strings.Trim(filename, "qr.")+".png")
				} else {
					fmt.Println("serving", filename+".png")
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

	http.HandleFunc("/kit.png", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "kit.png")
	}))

	http.HandleFunc("/the.keeper", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "the.keeper")
	})

	http.HandleFunc("/bad.habit", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "bad.habit")
	})

	http.HandleFunc("/the.presynct", func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "the.presynct")
	})

	http.HandleFunc("/🌞", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🌞")
	}))

	http.HandleFunc("/🌝", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🌝")
	}))

	http.HandleFunc("/🌚", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🌚")
	}))

	http.HandleFunc("/🌛", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🌛")
	}))

	http.HandleFunc("/🌜", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🌜")
	}))

	http.HandleFunc("/🏡", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🏡")
	}))

	http.HandleFunc("/🧑", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🧑")
	}))

	http.HandleFunc("/🚁", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "🚁")
	}))

	http.HandleFunc("/📼", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "📼")
	}))

	http.HandleFunc("/🖭", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		htmlHandler(w, r, "📼")
	}))

	http.HandleFunc("/emoji/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		emojiHandler(w, r, r.URL.Path[len("/emoji/"):])
	}))

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server closed")
		} else if err != nil {
			fmt.Println("error starting server: %w", errors.WithStack(err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Println("server shutdown failed: %w", err)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if strings.HasSuffix(origin, ".heliattack.com") {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next(w, r)
    }
}

func generateQR(filename, url string) error {
	qrCode, _ := qrcode.New(url, qrcode.Medium)
	err := qrCode.WriteFile(11, filename)
	if err != nil {
		return fmt.Errorf("%w", errors.WithStack(err))
	}
	return nil
}
