package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	gtagURL = "https://www.googletagmanager.com/gtag/js"
	gaURL   = "https://www.google-analytics.com/"
	prURL   = "https://gash.seankhliao.com/collect/"
)

func main() {
	pURL := flag.String("purl", prURL, "proxy url to collect reports at")
	t := flag.Duration("t", 48*time.Hour, "update frequency of js, also used for cache age")
	p := flag.String("p", "8080", "port to serve on")

	s := NewScript(*t, *pURL)
	go s.tick()

	shrp, err := NewReverseProxy(gaURL)
	if err != nil {
		log.Fatal("create ReverseProxy: ", err)
	}

	http.Handle("/collect", http.StripPrefix("/collect", shrp))
	http.HandleFunc("/js", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/js" {
			http.Redirect(w, r, "/js", http.StatusFound)
			return
		}

		w.Header().Add("Cache-Control", "max-age="+strconv.Itoa(int(s.d.Seconds())))
		w.Write(s.b)
	})

	http.ListenAndServe(":"+*p, nil)
}

func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("parse url: %v", err)
	}
	director := func(req *http.Request) {
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.URL.Path = path.Join(u.Path, req.URL.Path)

		vals := req.URL.Query()
		tip := strings.Split(req.RemoteAddr, ":")[0]
		if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
			tip = strings.Split(xff, ", ")[0]
		}
		vals.Add("uip", tip)
		for k, v := range u.Query() {
			for _, vv := range v {
				vals.Add(k, vv)
			}
		}
		req.URL.RawQuery = vals.Encode()

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}, nil
}

type Script struct {
	d time.Duration // update
	p []byte        // proxy host
	b []byte        // script bytes
}

func NewScript(d time.Duration, p string) *Script {
	return &Script{
		d: d,
		p: []byte(p),
	}
}

func (s *Script) tick() {
	err := s.getGtag()
	if err != nil {
		log.Fatal("get script: ", err)
	}
	for range time.NewTicker(s.d).C {
		err := s.getGtag()
		if err != nil {
			log.Println("get script: ", err)
		}
	}
}

func (s *Script) getGtag() error {
	res, err := http.Get(gtagURL)
	if err != nil {
		return fmt.Errorf("http get: %v", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response: %v", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("status: %v", err)
	}
	s.b = bytes.Replace(b, []byte(gaURL), s.p, -1)
	return nil
}
