package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}
var selectedImage = "images/ghost.gif"

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
	lp := filepath.Join("templates", "index.html")

	homeTemplate, parseError := template.ParseFiles(lp)

	if parseError != nil {
		fmt.Println(parseError)
	}

	data := struct {
		WebSocketUrl string
	}{
		WebSocketUrl: "ws://" + r.Host + "/updateDisplay",
	}

	err := homeTemplate.ExecuteTemplate(w, "index.html", data)

	if err != nil {
		fmt.Println(err)
	}
}

func display(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "display.html")
	homeTemplate, parseError := template.ParseFiles(lp)

	if parseError != nil {
		fmt.Println(parseError)
	}

	data := struct {
		WebSocketUrl  string
		SelectedImage string
	}{
		WebSocketUrl:  "ws://" + r.Host + "/updateDisplay",
		SelectedImage: selectedImage,
	}

	err := homeTemplate.ExecuteTemplate(w, "display.html", data)

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	log.SetFlags(0)
	http.HandleFunc("/updateDisplay", updateDisplay)
	http.HandleFunc("/display", display)
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./templates/images"))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
