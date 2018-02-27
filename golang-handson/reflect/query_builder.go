package reflect

import (
	"bytes"
	"fmt"
	"reflect"
	"time"
)

type QueryBuilder struct {
	Type reflect.Type
}

type Employee struct {
	ID        uint32 `orm:"id_pk"`
	FirstName string `orm:"first_name"`
	LastName  string `orm:"last_name"`
	Birthday  time.Time
}

func (qb *QueryBuilder) CreateSelectQuery() string {
	buffer := bytes.NewBufferString("")

	for i := 0; i < qb.Type.NumField(); i++ {
		field := qb.Type.Field(i)

		if i == 0 {
			buffer.WriteString("SELECT ")
		} else {
			buffer.WriteString(", ")
		}
		column := field.Name
		if tag := field.Tag.Get("orm"); tag != "" {
			column = tag
		}
		buffer.WriteString(column)
	}

	if buffer.Len() > 0 {
		fmt.Fprintf(buffer, " FROM %s", qb.Type.Name())
	}

	return buffer.String()
}
