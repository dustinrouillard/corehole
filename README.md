# Corehole

### DNS Adblocker to rid your network of those pesky malicious domains

Adblocking for the CoreDNS users :)

### CoreDNS Config Options

Set the `corehole` config option with `redis` and `resolve_host` to configure the module.

```
. {
  corehole {
    redis redis://valkey:6379/1
    resolve_host 10.90.1.9
  }
}
```

### Implementing with a blocked page

Originally I created this due to other DNS adblocking solutions not having a way to point to a custom A record so that a service can handle displaying a blocked page.

You'll notice this in the config section above under `resolve_host`, if you don't want this behavior for handling the blocked request to show a block page, you can set this to `127.0.0.1` *(the default)* like the other blockers do by default.

Be aware that some hosts the browser will refuse to load unless certain criteria is met, such as HSTS which would prevent the blocked page from loading on that domain.

### Example blocked records to store in redis:

```sh
valkey.kush:6379> KEYS dns_block:*
1) "dns_block:googletagservices.com"
2) "dns_block:analytics.twitter.com"
3) "dns_block:doubleclick.net"
valkey.kush:6379> GET dns_block:googletagservices.com
"1"
valkey.kush:6379> 
```

The record just has to exist it doesn't matter the value, however I just use "1", and I've run a few scripts to insert around 700k records from various blocklists found [here](https://github.com/blocklistproject/Lists)

### Using this with CoreDNS

In order to use this you will have to switch to a coredns build compiled with this plugin, for simplicity sake you can use the docker image that is built in github actions.

`ghcr.io/dustinrouillard/corehole`

If you want to build it yourself, you can just clone the [coredns](https://github.com/coredns/coredns) repo and add the following to the `plugin.cfg` file in the root, directly under `etcd:etcd`

```
etcd:etcd # Put it below this line, this is the order it executes the plugin on your query.
corehole:github.com/dustinrouillard/corehole
```

Get the package `go get github.com/dustinrouillard/corehole`

Then run `go generate`, followed by `go build` with whatever settings you need to build for your target environment.

and then to build the image as seen in the [workflow](https://github.com/dustinrouillard/corehole/blob/main/.github/workflows/production.yaml#L57) we replace the base image to debian and remove the nonroot user line.
