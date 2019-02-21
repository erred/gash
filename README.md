# gash

Google Analytics, Self Hosted

[![License](https://img.shields.io/github/license/seankhliao/gash.svg?style=for-the-badge&maxAge=31536000)](LICENSE)
[![Build](https://badger.seankhliao.com/i/github_seankhliao_gash)](https://badger.seankhliao.com/l/github_seankhliao_gash)

## About

Proxy Google Analytics through your own domain / server

1. Get the gtag script from: `https://https://www.googletagmanager.com/gtag/js`
   a. Replace `https://www.google-analytics.com/` with `https://gash.seankhliao.com`
   b. Cache the script
2. Get the analytics.js script from `https://www.google-analytics.com/analytics.js`
   a. Replace `https://www.google-analytics.com/` with `https://gash.seankhliao.com`
   b. Cache the script
3. Collect reports at `https://gash.seankhliao.com/collect/`
   a. add `uip=<true client ip address`
   b. forward to `https://www.google-analytics.com/collect`
4. Periodically update the scripts

## Usage

#### Prerequisites

- go

or

- docker

#### Install

go:

```sh
go get github.com/seankhliao/gash
```

#### Run

```sh
gash [-p 8080] [-t 48h]
  -p port to serve on
  -t update interval
```

docker:

```sh
docker run --rm \
  -p 8080:8080 \
  seankhliao/gash
```

#### Build

go:

```sh
go build
```

docker:

```sh
docker build \
  --network host \
  ,
```

## TODO

- [ ] verify client ip is valid
