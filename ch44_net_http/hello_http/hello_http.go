package main

import (
	"fmt"
	"net/http"

	"time"
)

func main() {
	//内置http服务
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	http.HandleFunc("/time/", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		timeStr := fmt.Sprintf("{\"time\": \"%s\"}", t)
		w.Write([]byte(timeStr))
	})
	//http://172.17.10.178:10110/  因为使用的是wsl2
	http.ListenAndServe(":10110", nil) //监听的端口 Handler
}
