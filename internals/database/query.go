package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTable() *sql.DB {
	_, errNofile := os.Stat("./internals/database/database.db")

	db, err := sql.Open("sqlite3", "./internals/database/database.db")
	if err != nil {
		log.Println(err.Error())
	}
	if errNofile != nil {
		sqlcode, err := os.ReadFile("./internals/database/table.sql")
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(string(sqlcode))

		if err != nil {
			log.Fatal(err)
		}
	}
	return db
}

func GeneratePrepare(text string) string {
	nb := len(strings.Split(text, ","))
	a := strings.Repeat("?,", nb)
	return "(" + a[:len(a)-1] + ")"
}

func Insert(db *sql.DB, table, values string, data ...interface{}) {
	prepare := GeneratePrepare(values)
	Query := fmt.Sprintf("INSERT INTO %v %v values %v", table, values, prepare)

	fmt.Println(Query)
	insert, err := db.Prepare(Query)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = insert.Exec(data...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

type Table interface {
	ScanRows(rows *sql.Rows) error
}

func (u Users) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&u.Id, &u.Username, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Registration)
}

type Users struct {
	Id           int
	Username     string
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Registration string
}

func Scan(db *sql.DB, table string, data Table) ([]Table, error) {
	Query := fmt.Sprintf("SELECT * FROM %v", table)
	stmt, err := db.Prepare(Query)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf(err.Error())
	}
	row, err := stmt.Query()
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf(err.Error())
	}

	tableData := []Table{}

	for row.Next() {

		dynamicType := reflect.New(reflect.TypeOf(data).Elem()).Interface().(Table)

		if err := dynamicType.ScanRows(row); err != nil {
			return nil, err
		}

		tableData = append(tableData, data)
	}
	return tableData, nil
}
