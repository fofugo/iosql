package iosql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPrepareQuery(t *testing.T) {
	expectation, err := GetPrepareQuery(QueryState{
		Columns: []string{
			"col1", "col2", "col3",
		},
		Table: "iosql_table",
	})
	assert.NoError(t, err)
	assert.Equal(t, expectation, "INSERT INTO iosql_table (col1,col2,col3) VALUES ")
}

func TestGetValues(t *testing.T) {
	assert.Equal(t, GetValues([]interface{}{"value", 1, 2}), "('value',1,2)")
}
func TestGetValuesWithNull(t *testing.T) {
	assert.Equal(t, GetValues([]interface{}{"value", "NULL", 2}), "('value',NULL,2)")
}

func TestConvertValues(t *testing.T) {
	assert.Equal(t,
		convertValues([]interface{}{"column1", 1, "column3"}),
		[]string{"'column1'", "1", "'column3'"},
	)
}
func TestConvertValuesWithNull(t *testing.T) {
	assert.Equal(t,
		convertValues([]interface{}{"NULL", 1, "column3"}),
		[]string{"NULL", "1", "'column3'"},
	)
}

func TestWrap(t *testing.T) {
	assert.Equal(t, wrap([]string{"v1", "v2", "v3"}), "(v1,v2,v3)")
}
