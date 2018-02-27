package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/k0kubun/pp"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if dump, err := httputil.DumpRequest(r, false); err == nil {
		pp.Println(string(dump))
	}
	fmt.Fprint(w, "<h1> Hello !</h1>")
}

func main() {
	http.HandleFunc("/", HelloHandler)
	host := ":80"
	fmt.Println("listning", host)
	http.ListenAndServe(host, nil)
}
