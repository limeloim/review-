package main

import (
	"fmt"
)

type Gopher struct {
	Id    int
	Name  string
	Age   int
	Email string
}

func main() {
	//создаем все струтктуры
	store, err := initDB("gopher.db") //Создали базу данных млм подключились к ней
	if err != nil {
		fmt.Println(err)
		return
	}
	//закроет соединение с дб после завершения работы приложения
	defer store.closeDBConnection()

	handler := NewGopherHandler(store)  //создали хандлер
	router := NewRouterGophers(handler) //создвли роутер
	err = router.Route(":8080")         //запускаем приложение
	if err != nil {
		fmt.Println(err)
	}
}
