package corehole

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	caddy.RegisterPlugin("corehole", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	r, err := configParse(c)
	if err != nil {
		return plugin.Error("corehole", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		r.Next = next
		return r
	})

	return nil
}

func configParse(c *caddy.Controller) (*CoreHole, error) {
	corehole := CoreHole{
		ttl:         300,
		resolveHost: "127.0.0.1",
		dnsPrefix:   "dns_block:",
	}

	for c.Next() {
		if c.NextBlock() {
			for {
				switch c.Val() {
				case "redis":
					if !c.NextArg() {
						return &CoreHole{}, c.ArgErr()
					}
					corehole.Redis = c.Val()
				case "ttl":
					if !c.NextArg() {
						return &CoreHole{}, c.ArgErr()
					}
					i, _ := strconv.Atoi(c.Val())
					corehole.ttl = i
				case "resolve_host":
					if !c.NextArg() {
						return &CoreHole{}, c.ArgErr()
					}
					corehole.resolveHost = c.Val()
				default:
					if c.Val() != "}" {
						return &CoreHole{}, c.Errf("unknown property '%s'", c.Val())
					}
				}

				if !c.Next() {
					break
				}
			}
		}

		corehole.Connect()

		return &corehole, nil
	}
	return &CoreHole{}, nil
}
