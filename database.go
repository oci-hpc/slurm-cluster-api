package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() {
	DbConn, err := sql.Open("sqlite3", "./slurm-cluster.db")
	if err != nil {
		log.Fatal("ERROR: InitDatabase: " + err.Error())
	}
	defer DbConn.Close()

	sql, err := ioutil.ReadFile("./deploy.sql")
	if err != nil {
		log.Fatal("ERROR: InitDatabase: " + err.Error())
	}

	sqlString := string(sql)

	_, err = DbConn.Exec(sqlString)
	if err != nil {
		log.Printf("SQL ERROR: InitDatabase: %q, Statement: %s\n", err, sqlString)
		return
	}
}

func GetDbConnection() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./slurm-cluster.db")
	if err != nil {
		log.Fatal("ERROR: InitDatabase: " + err.Error())
	}
	return
}
