// Copyright 2014 The Yxorp Authors. All rights reserved.

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"

	"github.com/hjr265/go-zrsc/zrsc"
)

func main() {
	log.Print("starting yxorp")

	go func() {
		base, err := url.Parse(cfg.Get("core.base").(string))
		catch(err)

		addr, ok := cfg.Get("core.addr").(string)
		if !ok {
			log.Fatal("missing core.addr in config.tml")
		}

		insecure, _ := cfg.Get("core.http.transport.tls.insecure").(bool)

		proxy := httputil.NewSingleHostReverseProxy(base)
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		}

		s := http.Server{
			Addr: addr,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				host, _, _ := net.SplitHostPort(r.RemoteAddr)
				hub.Send([]interface{}{"SRVD", fmt.Sprintf("%s - %s %s", host, r.Method, r.URL.String())})

				proxy.ServeHTTP(w, r)
			}),
		}

		log.Printf("core listening on %s", addr)
		s.ListenAndServe()
	}()

	go func() {
		addr, ok := cfg.Get("mond.addr").(string)
		if !ok {
			log.Fatal("missing mond.addr in config.tml")
		}

		http.HandleFunc("/hub", handleConnect)
		http.Handle("/", http.FileServer(zrsc.HttpDir("mond")))

		log.Printf("mond listening on %s", addr)
		err := http.ListenAndServe(addr, nil)
		catch(err)
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	log.Printf("received %s", <-c)
	log.Print("exiting")
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
