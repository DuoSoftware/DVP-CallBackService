// CallbackServer project main.go
package main

import (
	"fmt"
	"github.com/DuoSoftware/gorest"
	"net/http"
	"time"
)

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	fmt.Println("Hello World!")
	LoadConfiguration()
	go InitiateService()
	for {
		go ExecuteCallback()
		time.Sleep(externalCallbackRequestFrequency * time.Second)
	}
}

func InitiateService() {
	gorest.RegisterService(new(CallbackServerSelfHost))
	http.Handle("/", gorest.Handle())
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
