package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/miekg/dns"
	"github.com/wheresalice/pronounsapi/lib"
	"strings"
)

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeTXT:
			log.Printf("Query for %s\n", q.Name)
			txt := lib.GetUserPronouns(strings.TrimSuffix(q.Name, "."))
			txt = strings.TrimPrefix(txt, "{")
			txt = strings.TrimSuffix(txt, "}")
			if txt != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s TXT %s", q.Name, txt))
				if err != nil {
					log.Print(err)
				} else {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	port := 8053
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
