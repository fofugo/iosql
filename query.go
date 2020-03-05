package iosql

import (
	"database/sql"
	"fmt"
	"strings"
)

type QueryState struct {
	Path    string
	Db      *sql.DB
	Table   string
	Columns []string
	Size    int
	Msg     <-chan []interface{}
}

func GetPrepareQuery(queryState QueryState) (query string, err error) {
	if len(queryState.Columns) == 0 {
		err = fmt.Errorf("error wrong Columns: %v", queryState.Columns)
		return
	}
	if queryState.Table == "" {
		err = fmt.Errorf("error empty Table: %v", queryState.Table)
		return
	}
	wColumns := wrap(queryState.Columns)
	query = fmt.Sprintf("INSERT INTO %s %s VALUES ", queryState.Table, wColumns)
	return
}
func GetValues(values []interface{}) (wValues string) {
	cValues := convertValues(values)
	wValues = wrap(cValues)
	return
}
func convertValues(values []interface{}) (results []string) {
	for _, value := range values {
		switch v := value.(type) {
		case string:
			format := ""
			if v == "NULL" {
				format = "%s"
			} else {
				format = "'%s'"
			}
			results = append(results, fmt.Sprintf(format, v))
		default:
			results = append(results, fmt.Sprintf("%v", v))
		}
	}
	return
}

func wrap(contents []string) (result string) {
	result += "("
	for _, content := range contents {
		result += content + ","
	}
	result = strings.Trim(result, ",")
	result += ")"
	return
}
