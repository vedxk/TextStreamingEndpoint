package main

import (
	"log"
	"myapp/websocketserver"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", websocketserver.WSHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
