package main

import (
	"errors"
	"fmt"
	"sql-client/sqlclient"
)

type Language struct {
	Id   int
	Name string
}

var dbClient sqlclient.SqlClient

const (
	queryGetLangById = "SELECT id, name FROM language WHERE id=?;"
)

func init() {

	sqlclient.StartMockServer()

	var err error
	dbClient, err = sqlclient.Open("mysql", "root:root123@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}
}

func GetLangById(id int) (*Language, error) {

	sqlclient.AddMock(sqlclient.Mock{
		Query: "SELECT id, name FROM language WHERE id=?;",
		Args:  []interface{}{1},
		// Error:   errors.New("error creating query"),
		Columns: []string{"id", "name"},
		Rows: [][]interface{}{
			{1, "English"},
			{2, "French"},
		},
	})

	rows, err := dbClient.Query(queryGetLangById, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var language Language
	for rows.HasNext() {
		if err := rows.Scan(&language.Id, &language.Name); err != nil {
			fmt.Println("we got an error")
			return nil, err
		}
		return &language, nil
	}

	return nil, errors.New("language not found")
}

func main() {
	lang, err := GetLangById(2)
	if err != nil {
		panic(err)
	}

	fmt.Println(lang.Id)
	fmt.Println(lang.Name)
}
