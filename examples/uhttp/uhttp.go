package uhttp

import (
	"fmt"
	"net/http"
)

type SimpleHanle struct{}

func (*SimpleHanle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello goloader!"))
}

func Enter() {
	go func() {
		panic(http.ListenAndServe(":2300", http.FileServer(http.Dir("."))))
	}()
	sh := &SimpleHanle{}
	mux := http.NewServeMux()
	mux.Handle("/", sh)
	fmt.Println("start listen:9090")
	panic(http.ListenAndServe(":9090", mux))
}
