package proxy

import (
	"io"
	"net"
	"net/http"
)

func HTTP(target string) http.Handler {
	return TCP(target)
}

func WS(target string) http.Handler {
	return TCP(target)
}

func TCP(target string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, _ := net.Dial("tcp", target)
		defer d.Close()

		h, _, _ := w.(http.Hijacker).Hijack()
		defer h.Close()

		r.Write(d)

		c := make(chan error, 2)
		p := func(w io.Writer, r io.Reader) {
			io.Copy(w, r)
		}
		go p(d, h)
		go p(h, d)
		<-c
	})
}
