package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/outreach-demo/backend/handlers"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://app:"+handlers.DBPassword+"@localhost/outreach")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUser(db, w, r)
	})
	http.HandleFunc("/users/search", func(w http.ResponseWriter, r *http.Request) {
		handlers.SearchUsers(db, w, r)
	})
	http.HandleFunc("/users/role", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUserRole(db, w, r)
	})
	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/audit", handlers.AuditLog)
	http.HandleFunc("/files/download", handlers.DownloadFile)
	http.HandleFunc("/files/upload", handlers.UploadFile)
	http.HandleFunc("/proxy", handlers.ProxyFetch)
	http.HandleFunc("/webhook/relay", handlers.WebhookRelay)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
