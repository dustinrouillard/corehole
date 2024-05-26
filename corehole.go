package corehole

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/upstream"

	redisCon "github.com/gomodule/redigo/redis"
)

type CoreHole struct {
	Next        plugin.Handler
	Pool        *redisCon.Pool
	Redis       string
	resolveHost string
	ttl         int
	dnsPrefix   string
	upstream    *upstream.Upstream
}

func (corehole *CoreHole) A(name string) (answers, extras []dns.RR) {
	r := new(dns.A)
	r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeA,
		Class: dns.ClassINET, Ttl: uint32(corehole.ttl)}

	r.A = net.ParseIP(corehole.resolveHost)
	answers = append(answers, r)
	return
}

func (corehole *CoreHole) get(key string) int {
	conn := corehole.Pool.Get()
	if conn == nil {
		fmt.Println("error connecting to redis")
		return 0
	}
	defer conn.Close()

	reply, err := conn.Do("GET", corehole.dnsPrefix+key)
	if err != nil {
		return 0
	}
	val, err := redisCon.Int(reply, nil)
	if err != nil {
		return 0
	}
	return val
}

func (corehole *CoreHole) Connect() {
	corehole.Pool = &redisCon.Pool{
		MaxIdle:     3,
		MaxActive:   5,
		IdleTimeout: 0,
		Dial: func() (redisCon.Conn, error) {
			opts := []redisCon.DialOption{}
			opts = append(opts, redisCon.DialConnectTimeout(time.Duration(10)*time.Second))
			opts = append(opts, redisCon.DialReadTimeout(time.Duration(10)*time.Second))

			return redisCon.DialURL(corehole.Redis, opts...)
		},
	}
}
