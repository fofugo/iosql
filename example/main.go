package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/fofugo/iosql"
	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	Dialect  string
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Charset  string
}

type DB struct {
	Db *sql.DB
}

func (DB *DB) Initialize(config DbConfig) (err error) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.Charset)
	if DB.Db, err = sql.Open(config.Dialect, dbURI); err != nil {
		return
	}
	if err = DB.Db.Ping(); err != nil {
		return
	}
	return
}

func main() {
	dbConfig := DbConfig{
		Dialect:  "mysql",
		Host:     "127.0.0.1",
		Port:     "",
		Username: "",
		Password: "",
		Name:     "",
		Charset:  "",
	}
	db := DB{}
	if err := db.Initialize(dbConfig); err != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	msg := make(chan []interface{})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := iosql.BulkInsert(ctx, iosql.QueryState{
			Path:  "./",
			Db:    db.Db,
			Table: "iosql_table",
			Columns: []string{
				"col1", "col2", "col3",
			},
			Size: 3,
			Msg:  msg})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()
	for i := 0; i < 50; i++ {
		msg <- []interface{}{"co1", "1", "1"}
	}
	time.Sleep(4 * time.Second)
	cancel()
	wg.Wait()
	return
}
