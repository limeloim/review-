package main //сдесь будет логика работы приложения с базой данных:

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Store struct { //Новая структура db *sql.DB подтаскиваем в неё 1 раз и потом Передаем ссылку на эту структуру где уже есть db *sql.DB Методы которые есть у Store автоматически имеют доступ к db *sql.DB ТУТ реализуем CRUD
	db *sql.DB
}

func initDB(dbPath string) (*Store, error) { //Инициализация базы данных что бы это ни значело Принимает патч в праметры Возвращает ссылку на
	db, err := sql.Open("sqlite", dbPath)
	if err != nil { //НЕ ВАЖНО
		fmt.Println(err)
		return nil, err
	}

	err = createGopherTable(db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Store{db: db}, nil //Возвращает ссылку на БАЗУ ДАННЫХ? и нулевую ошибку
}

func createGopherTable(db *sql.DB) error { //SQL таблицу тут создадим поидее
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS gophers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		age INTEGER
	);`

	_, err := db.Exec(createTableSQL) //отправляем строку в базу данных
	if err != nil {
		err = fmt.Errorf("create gophers table: %w", err) //Передача ошибки наверх %w заполниится err
		return err
	}
	log.Printf("Table 'gophers' created.")
	return nil
}

func (s *Store) insertGopher(g *Gopher) (int64, error) { //CREATE Принимает гофера из запроса к серверу
	name := g.Name
	email := g.Email
	age := g.Age

	insertSQL := `INSERT INTO gophers (name, email, age) VALUES (?, ?, ?);`
	res, err := s.db.Exec(insertSQL, name, email, age)
	if err != nil {
		log.Fatalf("failed to insert: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func (s *Store) selectAllGophers() ([]*Gopher, error) { //На вход принимается База и читается .нужно сложить в слайс гоферов из базы READ метод
	const selectAllSQL = `SELECT id, name, email, age 
						   FROM gophers` //достаем сусликов из таблицы
	rows, err := s.db.Query(selectAllSQL) //он из s достает db и метод Query делает

	if err != nil {
		log.Fatalf("failed to insert: %v", err)
	}
	defer rows.Close()

	var gophers []*Gopher

	for rows.Next() { // Для каждой строки резултата
		var name, email string
		var age, id int
		err := rows.Scan(&id, &name, &email, &age)
		if err != nil {
			log.Fatal(err)
		}
		//Тут нужно создать суслика со сканированными данными поместить в []*Gopher,
		gophers = append(gophers, &Gopher{
			Id:    id,
			Name:  name,
			Email: email,
			Age:   age,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate gophers: %w", err)
	}

	return gophers, nil
}

func (s *Store) update(id int, g *Gopher) error { //принимаем новые данные гофера по адресу
	name := g.Name
	email := g.Email
	age := g.Age

	updateSQL := `UPDATE gophers SET name = ?, email = ?, age = ? WHERE id = ?;`
	_, err := s.db.Exec(updateSQL, name, email, age, id)
	if err != nil {

		log.Fatalf("failed to update: %v", err)
	}
	return nil
}

func (s *Store) delete(id int) error {
	deleteSQL := `DELETE FROM gophers WHERE id = ?;`
	_, err := s.db.Exec(deleteSQL, id)
	if err != nil {
		log.Fatalf("failed to update: %v", err)
	}
	return nil
}

func (s *Store) closeDBConnection() error {
	return s.db.Close()
}
