package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/olahol/melody"
)

//go:embed templates/*
var embededTemplates embed.FS

//go:embed images/*
var images embed.FS

var templates = template.Must(template.ParseFS(embededTemplates, "templates/*.html"))

var currentImage string

type ClientMessage struct {
	MessageType string
	Payload     string
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
		SelectedImage: currentImage,
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

func getAvailableImages(files []string) []string {
	result := make([]string, 0)
	validFileTypes := []string{".jpg", ".jpeg", ".png", ".gif"}

	for _, file := range files {
		for _, fileType := range validFileTypes {
			if strings.HasSuffix(file, fileType) {
				result = append(result, "/"+file)
			}
		}
	}

	log.Println(result)

	return result
}

func sendJsonMessage(s *melody.Session, messageType string, payload string) {
	var (
		newline = []byte{'\n'}
		space   = []byte{' '}
	)

	jsonMessage := []byte(`{"messageType": "` + messageType + `", "payload": "` + payload + `"}`)
	messageString := bytes.TrimSpace(bytes.Replace(jsonMessage, newline, space, -1))

	s.Write(messageString)
}

func main() {
	flag.Parse()

	m := melody.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	files, _ := getAllFilenames(images)
	availableImages := getAvailableImages(files)

	currentImage = availableImages[0]

	log.SetFlags(0)
	http.HandleFunc("/display", display)
	http.Handle("/images/", http.FileServer(http.FS(images)))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Recieved request for ws")
		m.HandleRequest(w, r)
	})
	http.HandleFunc("/", home)

	m.HandleConnect(func(s *melody.Session) {
		sendJsonMessage(s, "setAvailableImages", strings.Join(availableImages[:], ","))
		sendJsonMessage(s, "setImage", currentImage)
	})

	m.HandleMessage(func(s *melody.Session, message []byte) {
		var clientMessage ClientMessage

		err := json.Unmarshal(message, &clientMessage)

		if err == nil && clientMessage.MessageType == "setImage" {
			currentImage = clientMessage.Payload
		}

		m.Broadcast(message)
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
