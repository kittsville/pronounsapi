# Pronouns API

An HTTP API for retrieving people's pronouns, backed by CockroachDB.  Pronouns are also available via DNS query.

## Setup

You'll need an insecure CockroachDB server running locally

```
docker run -d --name=roach1 -p 26257:26257 -p 8080:8080 -v "${PWD}/data/cockroach:/cockroach/cockroach-data" cockroachdb/cockroach:v2.0.3 start --insecure
docker exec -ti roach1 ./cockroach user set pronounsapi --insecure
docker exec -ti roach1 ./cockroach sql -e 'GRANT ALL ON pronounsapi to pronounsapi' --insecure
```

You'll also need some Go dependencies installing, for which we're using Dep - `dep ensure`

And then run `go run init.go` to initialise the tables and add an example record.

You can see the status of Cockroach in a web browser http://localhost:8080

## Running

`go run http.go` will run a webserver on port 3000 which you can query with `curl localhost:3000/u/alice`

`go run dns.go` will run a dns server on port 8053 which you can query with `dig -t txt alice @127.0.0.1 -p 8053`

## TODO

* Secure the database properly
* Run this somewhere
* Add a web interface for adding records
* Come up with a cool domain name
* Design backup strategy
* Add monitoring