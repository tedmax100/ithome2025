package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "🎉 Hello from DevContainer! 🚀\n")
		fmt.Fprintf(w, "⏰ Current time: %s\n", now)
		fmt.Fprintf(w, "🐳 Running in container!\n")
	})

	fmt.Println("🚀 Server starting on :8080")
	fmt.Println("🌐 Visit: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
