package utils

import (
	"encoding/json"
	"errors"
	"reflect"
)

// MapToStruct converts a map to a struct of the specified type.
func MapToStruct(m map[string]interface{}, s interface{}) error {
	if s == nil {
		return errors.New("destination struct cannot be nil")
	}
	
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, s)
}

// StructToMap converts a struct to a map.
func StructToMap(s interface{}) (map[string]interface{}, error) {
	if s == nil {
		return nil, errors.New("source struct cannot be nil")
	}
	
	result := make(map[string]interface{})
	val := reflect.ValueOf(s).Elem()
	
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		result[field.Name] = val.Field(i).Interface()
	}
	
	return result, nil
}

// QueryBuilder helps to build SQL queries dynamically.
type QueryBuilder struct {
	table      string
	columns    []string
	conditions []string
}

// NewQueryBuilder initializes a new QueryBuilder for the specified table.
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

// Select specifies the columns to select.
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.columns = columns
	return qb
}

// Where adds conditions to the query.
func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// Build constructs the SQL query as a string.
func (qb *QueryBuilder) Build() string {
	query := "SELECT "
	if len(qb.columns) > 0 {
		query += join(qb.columns, ", ")
	} else {
		query += "*"
	}
	query += " FROM " + qb.table
	if len(qb.conditions) > 0 {
		query += " WHERE " + join(qb.conditions, " AND ")
	}
	return query
}

// join is a helper function to join strings with a separator.
func join(elements []string, separator string) string {
	if len(elements) == 0 {
		return ""
	}
	
	result := elements[0]
	for _, element := range elements[1:] {
		result += separator + element
	}
	return result
}