package main

import (
	"github.com/gorilla/mux"
	"log"
	"github.com/wheresalice/pronounsapi/lib"
	"github.com/miekg/dns"
	"github.com/coreos/bbolt"
	"github.com/wheresalice/pronounsapi/lib/database"
	"net/http"
	"os"
	"fmt"
)

func main() {
	database.DBCon, database.Err = bolt.Open("data/bolt.db", 0600, nil)
	if database.Err != nil {
		log.Fatalf("could not open db, %v", database.Err)
	}

	var router = mux.NewRouter()
	router.HandleFunc("/healthcheck", lib.HealthCheck).Methods("GET")

	router.HandleFunc("/u/{userName}", lib.HandleUser).Methods("GET")
	router.HandleFunc("/.well-known/webfinger", lib.HandleWebfinger).Methods("GET")
	router.HandleFunc("/a/{userName}", lib.HandleActor).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	go func() {
		httpPort := os.Getenv("PRONOUNS_HTTP_PORT")
		if len(httpPort) == 0 {
			httpPort = "3000"
		}

		log.Printf("Starting HTTP server at :%s\n", httpPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), router))
	}()

	go func() {
		// attach request handler func
		dns.HandleFunc(".", lib.HandleDnsRequest)

		dnsPort := os.Getenv("PRONOUNS_DNS_PORT")
		if len(dnsPort) == 0 {
			dnsPort = "5053"
		}

		// start dns server
		server := &dns.Server{Addr: ":" + dnsPort, Net: "udp"}
		log.Printf("Starting DNS server at %s\n", dnsPort)
		err := server.ListenAndServe()
		defer server.Shutdown()
		if err != nil {
			log.Fatalf("Failed to start server: %s\n ", err.Error())
		}
	}()
	select {}
}
