package iosql

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func BulkInsert(ctx context.Context, queryState QueryState) (err error) {
	matchs, err := filepath.Glob(queryState.Path + "iosql_*")
	if err != nil {
		return
	}
	for _, v := range matchs {
		if err = insertRows(queryState.Path+v, queryState.Db); err != nil {
			return
		}
		if err = os.Remove(queryState.Path + v); err != nil {
			return
		}
	}
	count := 0
	fileName := queryState.Path + "/iosql_" + getGoId()
	f, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)

	defer func() {
		if err = insertRows(fileName, queryState.Db); err != nil {
			return
		}
		f.Close()
		_ = os.Remove(fileName)
	}()
	if err != nil {
		return
	}
	prefixQuery, err := GetPrepareQuery(queryState)
	fmt.Fprintf(f, "%s ", prefixQuery)
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-queryState.Msg:
			if !ok {
				return
			}
			fmt.Fprintf(f, " %s,", GetValues(v))
			count++
			if count > queryState.Size {
				if err = insertRows(fileName, queryState.Db); err != nil {
					return
				}
				count = 0
				f.Close()
				f, _ = os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
				fmt.Fprintf(f, "%s ", prefixQuery)
			}
		}
	}
	return
}

func insertRows(fileName string, db *sql.DB) (err error) {
	result, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	query := string(result)
	query = strings.Trim(query, ",")
	_, err = db.Exec(query)
	return
}

func getGoId() (id string) {
	var buf [32]byte
	n := runtime.Stack(buf[:], false)
	id = strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	return
}
