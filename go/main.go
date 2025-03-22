package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
)

func main() {
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<h1>Welcome to the Image Server</h1>
			<img src="/images/image1.jpg" alt="Image 1" width="300">
			<img src="/images/image2.jpg" alt="Image 2" width="300">
		`)
	})

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
