# Go Load Balancer

This is my own implementation of an API Load Balancer in Go using different algorithms, concurrent programming and synchronization techniques. Coding Challenges series by John Crickett https://codingchallenges.fyi/challenges/challenge-load-balancer/

-   Active and Passive Healthchecks
-   Healthcheck url on servers with checks dependent resources (Redis, DB, etc)
-   Round-Robin requests to all active servers

TODO:

-   Different algorithms: Round-Robin, Least-Connections (TODO)
-   Retry & Timeout mechanism (TODO)
-   Self register new servers to the LB pool

## Requirements

-   Install [GO v1.21](https://go.dev/dl/)
-   Install Make [Windows](https://gnuwin32.sourceforge.net/packages/make.htm) | [Mac] brew install make

## Setup

```
make build
```

```
make install
```

## Run

### Servers

Start two or more servers using different ports

```
go-lb --type node --addr :8081
```

```
go-lb --type node --addr :8082
```

### Load Balancer

Start the Load Balancer and the urls of the existing servers

```
go-lb --type lb --addr :8080 --urls http://localhost:8081,http://localhost:8082
```

## Tests

```
curl --parallel --parallel-immediate --parallel-max 4 --config urls.txt
```
