package main

import (
	"log"
	"net/http"

	"forum/internal/web/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app, err := server.InitApp()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/chat.html")
	})

	serveur := http.Server{
		Addr:    ":8228",
		Handler: app.Routes(),
	}

	log.Println("\u001b[38;2;255;165;0mListing on http://localhost:8228...\033[0m")
	log.Fatal(serveur.ListenAndServe())
}
