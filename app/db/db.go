package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/michgur/puncher/app/model"
	_ "modernc.org/sqlite"
)

var sqlsPath = "./sqls"
var sqls = map[string]string{}
var db *sql.DB

func ReadSQLs() error {
	// iterate over ./sqls, store the contents of each file in a map
	// key: filename, value: contents

	fnames, err := os.ReadDir(sqlsPath)
	if err != nil {
		return err
	}

	for _, fname := range fnames {
		f, err := os.ReadFile(sqlsPath + "/" + fname.Name())
		if err != nil {
			return err
		}
		sqls[fname.Name()] = string(f)
		println(fname.Name())
		println(fname.Name())
		println(fname.Name())
	}

	return nil
}

func init() {
	// fetch queries
	var err = ReadSQLs()
	// create the database
	db, err = sql.Open("sqlite", "puncher.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	// Create CardInstances table
	_, err = db.Exec(sqls["card-instances-create.sql"])
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}
}

func GetAllCardInstances() (cardInstances []model.CardInstance, err error) {
	rows, err := db.Query(sqls["card-instances-select-all.sql"])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ci model.CardInstance
		rows.Scan(&ci.ID, &ci.CardID, &ci.Slots)
		cardInstances = append(cardInstances, ci)
	}

	return cardInstances, nil
}

func InsertCardInstance(cardInstance model.CardInstance) error {
	_, err := db.Exec(sqls["card-instance-insert.sql"], cardInstance.CardID, cardInstance.Slots)
	return err
}

// wrappers for db functions, using SQL files from ./sqls
func Exec(qName string, args ...interface{}) (sql.Result, error) {
	if query, ok := sqls[qName]; ok {
		return db.Exec(query, args...)
	}
	return nil, fmt.Errorf("query %s not found", qName)
}

func Query(qName string, args ...interface{}) (*sql.Rows, error) {
	if query, ok := sqls[qName]; ok {
		return db.Query(query, args...)
	}
	return nil, fmt.Errorf("query %s not found", qName)
}

func QueryRow(qName string, args ...interface{}) *sql.Row {
	if query, ok := sqls[qName]; ok {
		return db.QueryRow(query, args...)
	}
	return nil
}
