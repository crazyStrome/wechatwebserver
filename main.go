package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I love you"))
	})
	http.HandleFunc("/wechat", wechat)
	http.ListenAndServe(":8080", nil)
}
func wechat(w http.ResponseWriter, r *http.Request) {
	var content, _ = ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	log.Println(content)
	w.Write([]byte("crazstom"))
}
