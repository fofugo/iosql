[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges/)
<h2>import</h2>

```
go get -u github.com/fofugo/iosql
```

<h2>contribution</h2>
<ul>
  <li>No guide line</li>
  <li>Just contribute</li>
</ul>

<h2>feature</h2>
<ul>
  <li>Bulk insert need safety way, against panic</li>
  <li>so iosql offers the way by file io query in dir</li>
</ul>

<h2>example</h2>

```
+-------+---------+------+-----+---------+-------+
| Field | Type    | Null | Key | Default | Extra |
+-------+---------+------+-----+---------+-------+
| col1  | char(5) | YES  |     | NULL    |       |
| col2  | int     | YES  |     | NULL    |       |
| col3  | int     | YES  |     | NULL    |       |
+-------+---------+------+-----+---------+-------+
```

```go
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
		msg <- []interface{}{"value1", 1, 3}
	}
	time.Sleep(4 * time.Second)
	cancel()
	wg.Wait()
	return
}

```
