package lib

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/wheresalice/pronounsapi/lib/database"
	"log"
	"strings"
)

func parseHostname(fqdn string) string {
	return strings.Split(fqdn, ".")[0]
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeTXT:
			log.Printf("Query for %s\n", q.Name)
			username := parseHostname(q.Name)
			txt := database.GetUserData(username).Pronouns
			rr, err := dns.NewRR(fmt.Sprintf("%s TXT %s", q.Name, txt))
			if err != nil {
				log.Print(err)
			} else {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func HandleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}
