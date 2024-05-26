package corehole

import (
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

func (corehole *CoreHole) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	qname := state.Name()
	qtype := state.Type()

	location := corehole.get(strings.TrimSuffix(qname, "."))
	if location == 0 {
		return plugin.NextOrFailure(qname, corehole.Next, ctx, w, r)
	}

	answers := make([]dns.RR, 0, 10)

	if qtype == "A" {
		answers, _ = corehole.A(qname)
	} else {
		return plugin.NextOrFailure(qname, corehole.Next, ctx, w, r)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	m.Answer = append(m.Answer, answers...)

	state.SizeAndDo(m)
	m = state.Scrub(m)
	w.WriteMsg(m)

	return dns.RcodeSuccess, nil
}

func (corehole *CoreHole) Name() string { return "corehole" }

func (corehole *CoreHole) errorResponse(state request.Request, zone string, rcode int, err error) (int, error) {
	m := new(dns.Msg)
	m.SetRcode(state.Req, rcode)
	m.Authoritative = true

	state.SizeAndDo(m)
	state.W.WriteMsg(m)
	return dns.RcodeSuccess, err
}
