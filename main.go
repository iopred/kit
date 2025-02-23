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
	<img src='/{{.Filename}}.png' id="qr" onclick="hideElement(this)" style="cursor: pointer;"/>

	<script>
        function hideElement(element) {
            element.style.display = 'none';
        }
    </script>
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
	Nodes   rune
}

func (n Node) x() string {
	if n.X {
		return fmt.Sprintf("Node{X:%v, Y:%v, Z:%v, Gravity:%v, Nodes:%d} intersects with spacetime vector (x, 0, 0, 0)", n.X, n.Y, n.Z, n.Gravity, n.Nodes)
	}
	return "0"
}

func (n Node) y() string {
	if n.Y {
		return fmt.Sprintf("Node{X:%v, Y:%v, Z:%v, Gravity:%v, Nodes:%d} intersects with spacetime vector (0, y, 0, 0)", n.X, n.Y, n.Z, n.Gravity, n.Nodes)
	}
	return "0"
}

func (n Node) z() string {
	if n.Z {
		return fmt.Sprintf("Node{X:%v, Y:%v, Z:%v, Gravity:%v, Nodes:%d} intersects with spacetime vector (0, 0, z, 0)", n.X, n.Y, n.Z, n.Gravity, n.Nodes)
	}
	return "0"
}

func (n Node) g() string {
	if n.Gravity {
		return fmt.Sprintf("Node{X:%v, Y:%v, Z:%v, Gravity:%v, Nodes:%d} intersects with spacetime vector (0, 0, 0, g)", n.X, n.Y, n.Z, n.Gravity, n.Nodes)
	}
	return "0"
}

func (n Node) kat() Node {
	if n.Gravity {
		return Node{n.X, true, n.Z, n.Gravity, n.Nodes}
	}
	return Node{n.X, n.Y, n.Z, n.Gravity, n.Nodes}
}

func (n Node) now() string {
	return fmt.Sprintf("%03d9", n.Nodes)
}

func (n Node) next() string {
	switch n.now() {
	case "0000":
		return "0001"
	case "0001":
		return "0002"
	case "0002":
		return "0003"
	case "0003":
		return "0004"
	case "0004":
		return "0005"
	case "0005":
		return "0006"
	case "0006":
		return "0007"
	case "0007":
		return "0008"
	case "0008":
		return "0009"
	case "0009":
		return "🌞"
	}

	return string(n.Nodes)
}

func (n Node) renderFrame() string {
	// Create a recursive tree structure using emojis
	// Root emoji is the current time/state
	root := n.now()

	// Build branches using different emojis based on node properties
	var branches []string

	// Add branches based on node properties
	if n.X {
		branches = append(branches, "→🔴") // Red for X
	}
	if n.Y {
		branches = append(branches, "↑🟡") // Yellow for Y
	}
	if n.Z {
		branches = append(branches, "↗️🔵") // Blue for Z
	}
	if n.Gravity {
		branches = append(branches, "↓⚫") // Black for Gravity
	}

	// Add sub-nodes based on Nodes byte value

	// todo: append sub nodes to string following the path.
	// for i := byte(0); i < n.Nodes; i++ {
	//     branches = append(branches, "⤵️"+string(rune('0'+i)))
	// }

	// Combine root with branches
	result := root
	if len(branches) > 0 {
		result += " " + strings.Join(branches, " ")
	}

	// Add next state
	// result += " → " + n.kat().next()

	return result
}

// Frame: 0009 →🔴 ↗️🔵 ↓⚫ ⤵️0 ⤵️1 → 🌞

/*
func (n Node) renderFrame() THREE.Group {
	// Create a group to hold the nodes
	group := THREE.NewGroup()

	// Create a box for each node
	box := THREE.NewBoxGeometry(1, 1, 1)

	// Create a material for the nodes
	material := THREE.NewMeshBasicMaterial(THREE.Color(0x00ff00))

	// Create a mesh for each node
	mesh := THREE.NewMesh(box, material)

	mesh.Position().Set(0, 0, 0)
	group.Add(mesh)

	// Define a recursive function to build the node tree.
	// This function adds a child node (as a red box) at a given position and scale,
	// then recurses to add child nodes in all six cardinal directions.
	var addChildNodes func(parent THREE.Group, depth int, x, y, z, scale float64)
	addChildNodes = func(parent THREE.Group, depth int, x, y, z, scale float64) {
		if depth <= 0 {
			return
		}
		// Create a child node mesh with a scaled box and red material.
		childBox := THREE.NewBoxGeometry(scale, scale, scale)
		childMaterial := THREE.NewMeshBasicMaterial(THREE.Color(0xff0000))
		childMesh := THREE.NewMesh(childBox, childMaterial)
		childMesh.Position().Set(x, y, z)
		parent.Add(childMesh)

		// Set offset and next scale for deeper recursion.
		offset := scale * 2.0
		nextScale := scale * 0.5

		// Recursively add child nodes in six directions: +X, -X, +Y, -Y, +Z, -Z.
		addChildNodes(parent, depth-1, x+offset, y, z, nextScale)
		addChildNodes(parent, depth-1, x-offset, y, z, nextScale)
		addChildNodes(parent, depth-1, x, y+offset, z, nextScale)
		addChildNodes(parent, depth-1, x, y-offset, z, nextScale)
		addChildNodes(parent, depth-1, x, y, z+offset, nextScale)
		addChildNodes(parent, depth-1, x, y, z-offset, nextScale)
	}

	// Start the recursive generation with a chosen depth and initial parameters.
	// Here, we use depth 3 and begin from position (2.0, 2.0, 2.0) with an initial scale of 0.5.
	addChildNodes(group, 3, 2.0, 2.0, 2.0, 0.5)

	// Return the complete group representing the recursive threejs scene.
	return group


}
*/

func loadKit() (Node, error) {
	return Node{false, false, false, true, 0001}, nil
}
