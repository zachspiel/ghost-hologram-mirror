package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
    "embed"
    "io/fs"
	"os"
)

//go:embed templates/*
var embededTemplates embed.FS

//go:embed images/*
var images embed.FS

var templates = template.Must(template.ParseFS(embededTemplates, "templates/*.html"))

var upgrader = websocket.Upgrader{}
var selectedImage = "/images/ghost.gif"

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func updateDisplay(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer connection.Close()
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("Got message: %s", message)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		WebSocketUrl string
	}{
		WebSocketUrl: "ws://" + r.Host + "/updateDisplay",
	}

	err := templates.ExecuteTemplate(w, "index.html", data)

	if err != nil {
		fmt.Println(err)
	}
}

func display(w http.ResponseWriter, r *http.Request) {
	data := struct {
		WebSocketUrl  string
		SelectedImage string
	}{
		WebSocketUrl:  "ws://" + r.Host + "/updateDisplay",
		SelectedImage: selectedImage,
	}

	err := templates.ExecuteTemplate(w, "display.html", data)

	if err != nil {
		fmt.Println(err)
	}
}

func getAllFilenames(efs embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
 
		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

	files, _ := getAllFilenames(images)
	
	log.SetFlags(0)
	http.HandleFunc("/updateDisplay", updateDisplay)
	http.HandleFunc("/display", display)
	http.Handle("/images/", http.FileServer(http.FS(images)))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    	log.Println("Recieved request for ws")
		serveWs(hub, w, r, files)
	})
	http.HandleFunc("/", home)

    log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
