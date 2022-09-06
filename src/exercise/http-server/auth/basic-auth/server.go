package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// We store bcrypt-ed passwords for each user. The actual passwords are "1234"
// for "joe" and "strongerpassword9902" for "mary", but these should not be
// stored anywhere.
var usersPasswords = map[string][]byte{
	"joe":  []byte("$2a$12$aMfFQpGSiPiYkekov7LOsu63pZFaWzmlfm1T8lvG6JFj2Bh4SZPWS"),
	"mary": []byte("$2a$12$l398tX477zeEBP6Se0mAv.ZLR8.LZZehuDgbtw2yoQeMjIyCNCsRW"),
}

// go run .\generate_cert.go -ecdsa-curve P256 -host localhost
// generate_cert.go 文件位置位于go源码包位置如：C:\Program Files\Go\src\crypto\tls
func main() {
	addr := flag.String("addr", ":4000", "HTTPS network address")
	certFile := flag.String("certfile", "cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "key.pem", "key PEM file")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Proudly served with Go and HTTPS!\n")
	})

	mux.HandleFunc("/secret/", func(w http.ResponseWriter, r *http.Request) {
		// http 标准认证源自：https://tools.ietf.org/html/rfc7617
		// 现已因没有安全性而基本禁止
		user, pass, ok := r.BasicAuth()
		if ok && verifyUserPass(user, pass) {
			fmt.Fprintf(w, "You get to see the secret\n")
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})

	srv := &http.Server{
		Addr:    *addr,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	log.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServeTLS(*certFile, *keyFile)
	log.Fatal(err)
}

func verifyUserPass(username, password string) bool {
	wantPass, hasUser := usersPasswords[username]
	if !hasUser {
		return false
	}
	if cmperr := bcrypt.CompareHashAndPassword(wantPass, []byte(password)); cmperr == nil {
		return true
	}
	return false
}
