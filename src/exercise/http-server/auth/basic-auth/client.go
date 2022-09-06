package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", "localhost:4000", "HTTPS server address")
	certFile := flag.String("certfile", "cert.pem", "trusted CA certificate")
	user := flag.String("user", "", "username")
	pass := flag.String("pass", "", "password")
	flag.Parse()

	cert, err := os.ReadFile(*certFile)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", *certFile)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	// Set up HTTPS request with basic authorization.
	req, err := http.NewRequest(http.MethodGet, "https://"+*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(*user, *pass)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HTTP Status:", resp.Status)
	fmt.Println("Response body:", string(html))
}
