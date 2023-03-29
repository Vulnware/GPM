package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/packr/v2"
)

func httpServer() {
	// Start WebSocket server
	templates := packr.New("templates", "./templates")
	statics := packr.New("static", "./static") // static files like css, js, images, etc.
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/", http.FileServer(templates)) // serve index.html
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(statics)))
	// print all files in the box
	for _, name := range templates.List() {
		fmt.Println(name)
	}
	for _, name := range statics.List() {
		fmt.Println(name)
	}

	log.Printf("Starting server on http://%s:%s", CONFIGS["host"], CONFIGS["port"])
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", CONFIGS["host"], CONFIGS["port"]), nil))
}
