package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterGophers struct {
	handler *GopherHandler
}

func NewRouterGophers(handler *GopherHandler) *RouterGophers { //это из мейна мы будем использовать чтоб насторить роутер
	return &RouterGophers{
		handler: handler,
	}
}

// addr в данном случае это по сути порт, который будет слушать приложение
func (n *RouterGophers) Route(addr string) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger) //это для логирование
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/gopher", func(r chi.Router) {
		r.Get("/", n.handler.GetGophers)       //READ
		r.Post("/", n.handler.CreateGopher)    //CREATE
		r.Put("/{id}", n.handler.UpdateGopher) //UPDATEbyID
		r.Delete("/{id}", n.handler.DeleteGopher)
	})

	fmt.Println("Listening on " + addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}
	return nil
}
