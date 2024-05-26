package proxy

import (
	"fmt"
	"io"
	"net/http"
)

type MagicHost struct {
	Host string
	Port int
}

func Start(c chan MagicHost) {
	table := make(map[string]int)

	go func() {
		for {
			select {
			case p := <-c:
				table[p.Host] = p.Port
			}
		}
	}()

	go func() {
		client := new(http.Client)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := table[r.Host]
			r.URL.Scheme = "http"
			r.URL.Host = fmt.Sprintf("localhost:%d", p)
			r.RequestURI = ""
			fmt.Println("req", r.URL)
			res, err := client.Do(r)
			if err != nil {
				fmt.Println("bad gateway", err)
				w.WriteHeader(http.StatusBadGateway)
				return
			}

			w.WriteHeader(res.StatusCode)
			for k, v := range res.Header {
				for _, vv := range v {
					w.Header().Add(k, vv)
				}
			}

			defer res.Body.Close()
			io.Copy(w, res.Body)
		})
		err := http.ListenAndServe(":8080", handler)
		if err != nil {
			panic(err)
		}
	}()
}
