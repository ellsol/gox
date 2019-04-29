package sqlx

import (
	"bytes"
)


type ColumnDefinition struct {
	TableName string
	Columns   []TableColumn
}

func NewColumnDefinition(tableName string) *ColumnDefinition {
	return &ColumnDefinition{
		TableName: tableName,
		Columns:   make([]TableColumn, 0),
	}
}

/*
	Encodes column definition as SQL string
 */
func (it *ColumnDefinition) SqlString() string {
	var buffer bytes.Buffer

	buffer.WriteString("CREATE TABLE ")
	buffer.WriteString(it.TableName)
	buffer.WriteString("(")

	for k, v := range it.Columns {
		lastElement := k == len(it.Columns)-1
		buffer.WriteString(v.Statement(!lastElement))
	}

	buffer.WriteString(");")
	return buffer.String()
}


func (it *ColumnDefinition) ColumnNames() []string {
	result := make([]string, 0)

	for _, v := range it.Columns {
		result = append(result, v.Name)
	}

	return result
}


func (it *ColumnDefinition) WithColumnDefinition(name string, valueType string, notNull bool) *ColumnDefinition {
	col := TableColumn{
		Name:      name,
		NotNull:   notNull,
		Type:      valueType,
	}

	it.Columns = append(it.Columns, col)
	return it
}

/*
	Gets the last column element and marks it as primary
 */
func (it *ColumnDefinition) AsPrimary() *ColumnDefinition {
	lastElementPos := len(it.Columns) - 1

	// ignore if no element set yet
	if lastElementPos < 0 {
		return it
	}

	lastElement := it.Columns[lastElementPos]
	lastElement.IsPrimary = true

	// need to put it back
	it.Columns = append(it.Columns[:lastElementPos], lastElement)

	return it
}

func (builder *ColumnDefinition) WithSerialColumn(name string) *ColumnDefinition {
	return builder.WithColumnDefinition(name, "SERIAL", false)
}

func (builder *ColumnDefinition) WithTextColumn(name string, notNull bool) *ColumnDefinition {
	return builder.WithColumnDefinition(name, "TEXT", notNull)
}

func (builder *ColumnDefinition) WithBooleanColumn(name string, notNull bool) *ColumnDefinition {
	return builder.WithColumnDefinition(name, "BOOLEAN", NotNull)
}

func (builder *ColumnDefinition) WithBigIntColumn(name string, notNull bool) *ColumnDefinition {
	return builder.WithColumnDefinition(name, "BIGINT", notNull)
}

func (builder *ColumnDefinition) WithIntColumn(name string, notNull bool) *ColumnDefinition {
	return builder.WithColumnDefinition(name, "INT", notNull)
}

func (builder *ColumnDefinition) WithByteAColumn(name string,  notNull bool) *ColumnDefinition {
	return builder.WithColumnDefinition(name, "BYTEA", notNull)
}




//////////////////////////////////////
//
// TableColumn
//
/////////////////////////////////////

type TableColumn struct {
	Name      string
	Type      string
	IsPrimary bool
	NotNull   bool
}

func (column *TableColumn) Statement(withComma bool) string {
	var buffer bytes.Buffer

	buffer.WriteString(column.Name)
	buffer.WriteString(" ")
	buffer.WriteString(column.Type)

	if column.IsPrimary {
		buffer.WriteString(" PRIMARY KEY")
	}

	if column.NotNull {
		buffer.WriteString(" NOT NULL")
	}

	if withComma {
		buffer.WriteString(",")
	}

	return buffer.String()
}