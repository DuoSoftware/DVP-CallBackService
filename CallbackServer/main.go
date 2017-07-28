// CallbackServer project main.go
package main

import (
	"fmt"
	"github.com/DuoSoftware/gorest"
	"github.com/rs/cors"
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
	InitiateRedis()
	go InitiateService()
	for {
		go ExecuteCallback()
		time.Sleep(externalCallbackRequestFrequency)
	}
}

func InitiateService() {
	gorest.RegisterService(new(CallbackServerSelfHost))
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"accept", "authorization"},
	})
	handler := c.Handler(gorest.Handle())
	addr := fmt.Sprintf(":%s", port)
	s := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//s.SetKeepAlivesEnabled(false)
	s.ListenAndServe()
}
