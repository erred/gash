package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	gtagURL = "https://www.googletagmanager.com/gtag/js"
	gaURL   = "https://www.google-analytics.com/analytics.js"
	cURL    = "https://www.google-analytics.com/collect"
	rsURL   = []byte("https://www.google-analytics.com/")
	rtURL   = []byte("https://gash.seankhliao.com/")
)

func main() {
	t := flag.Duration("t", 48*time.Hour, "update frequency of js, also used for cache age")
	p := flag.String("p", "8080", "port to serve on")

	gtag := NewScript(gtagURL)
	go gtag.update(*t)

	analytics := NewScript(gaURL)
	go analytics.update(*t)

	http.Handle("/js", gtag)
	http.Handle("/analytics.js", analytics)
	http.Handle("/collect/analytics.js", analytics)
	http.HandleFunc("/collect", collector)

	http.ListenAndServe(":"+*p, nil)
}

type Script struct {
	u string
	b []byte
}

func NewScript(u string) *Script {
	return &Script{
		u: u,
	}
}

func (s *Script) update(d time.Duration) {
	s.get()
	for range time.NewTicker(d).C {
		s.get()
	}
}

func (s *Script) get() {
	res, err := http.Get(s.u)
	if err != nil {
		log.Println("get: ", err)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		log.Println("get status: ", res.StatusCode, res.Status)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("read get: ", err)
	}
	s.b = bytes.Replace(b, rsURL, rtURL, -1)
}

func (s *Script) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(s.b)
}

func collector(w http.ResponseWriter, r *http.Request) {
	tip := strings.Split(r.RemoteAddr, ":")[0]
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		tip = strings.Split(xff, ", ")[0]
	}

	vals := r.URL.Query()
	vals.Add("uip", tip)
	u, err := url.Parse(cURL)
	if err != nil {
		log.Println("parse url: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	u.RawQuery = vals.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		log.Println("get: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(res.StatusCode)
}
