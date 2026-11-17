package main // предать в стор который работает с бд запрос ошибки тут обработать и выслать в ебало клиенту

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

//выше можете описать структуры, которые будете принимать на вход в случае POST/PUT запросов

type GopherHandler struct { //Обработчик запросов от сервера // тут быдыут функции CRUD но они будут активировать функции store у которых есть доступ к базе
	store *Store //передаем
}

func NewGopherHandler(store *Store) *GopherHandler {
	return &GopherHandler{
		store: store,
	}
}

// READ
func (h *GopherHandler) GetGophers(w http.ResponseWriter, r *http.Request) { //функция GopherHandler имеет доступ к store прослойке которая реализует общение с базой
	//логика обработки GET запроса
	gophers, err := h.store.selectAllGophers() // на выходе получили слайс из базы данных через прослойку в виде Store
	if err != nil {
		http.Error(w, "Ошибка базы данных", 500)
		return
	}
	//записываем ответ
	w.Header().Set("Content-Type", "application/json") //настраиваем ответ в виде json формата
	w.WriteHeader(http.StatusCreated)
	//gopher - это слайс структур, который мы достали ранее. С помощью Encode() мы преобразуем его в json и записываем в ответ
	if err := json.NewEncoder(w).Encode(gophers); err != nil { //прикольно можно одновременно создать err выполнить функцию и проверить условие
		return
	}
}

// CREATE
func (h *GopherHandler) CreateGopher(w http.ResponseWriter, r *http.Request) {
	//Сначаоа принять json с гофером и оправить его в Store чтоб он создал в базе его
	var g Gopher //создаем переменную типа Gopher

	json.NewDecoder(r.Body).Decode(&g)  //в гофера распаковали данные из Джейсона
	id, err := h.store.insertGopher(&g) // функция возвращает айди и ощтибку принимаем соответственно

	if err != nil {
		return
	}
	g.Id = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created — успешно создано!
	json.NewEncoder(w).Encode(g)      //обновиоли перменные
}

// UPDATE
func (h *GopherHandler) UpdateGopher(w http.ResponseWriter, r *http.Request) {
	//чначала принимаем какого гофера обновлятьб

	idString := chi.URLParam(r, "id") //URL был /gopher/5, то idStr будет "5" РАЗОБРАТЬСЯ НУЖНО ПРЕРЕДАТЬ Сюда уже инт
	id, err := strconv.Atoi(idString)
	if err != nil {
		return
	}

	var g Gopher
	json.NewDecoder(r.Body).Decode(&g)             //в гофера распаковали данные из Джейсона
	if err := h.store.update(id, &g); err != nil { // функция возвращает айди и ощтибку принимаем соответственно
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE

func (h *GopherHandler) DeleteGopher(w http.ResponseWriter, r *http.Request) {
	//чначала принимаем какого гофера обновлятьб

	idString := chi.URLParam(r, "id") //URL был /gopher/5, то idStr будет "5" РАЗОБРАТЬСЯ НУЖНО ПРЕРЕДАТЬ Сюда уже инт
	id, err := strconv.Atoi(idString)
	if err != nil {
		return
	}

	if err := h.store.delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
