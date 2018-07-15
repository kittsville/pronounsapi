# Pronouns API

An HTTP API for retrieving people's pronouns, backed by BoltDB.  Pronouns are also available via DNS query.

## Setup

You'll need some Go dependencies installing, for which we're using Dep - `dep ensure`

And then run `go run init.go` to initialise the tables and add an example record.

## Running

`go run main.go`

You should now have both an http server and a dns server running.

The default ports are 3000 for http and 5053 for dns, but you can configure these with the environmental variables `PRONOUNS_HTTP_PORT` and `PRONOUNS_DNS_PORT`.

## Usage

```
dig -t txt alice @localhost -p 5053 # get Alice's pronouns by dns
curl localhost:3000/u/alice # get Alice's pronouns by http

curl localhost:3000/.well-known/webfinger?resource=alice@pronouns.tech # web finger!
curl localhost:3000/a/alice # activity feeds!
```

## TODO

* Add a web interface for adding records
* Design backup strategy
* Add monitoring