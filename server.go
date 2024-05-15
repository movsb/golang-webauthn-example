package main

import (
	_ "embed"
	"log"
	"net/http"
)

func main() {
	store := NewStore()
	wa := NewWebAuthn(store,
		`localhost`,
		`Golang WebAuthn Example`,
		[]string{`https://localhost:2345`},
	)
	const prefix = `/admin/`
	admin := NewAdmin(store, wa, prefix)
	http.Handle(prefix, admin.Handler())
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	log.Fatalln(http.ListenAndServeTLS(":2345", "localhost.crt", "localhost.key", nil))
}
