package sqlx

import (
	"bytes"
)


type CreateStatementBuilder struct {
	TableName string
	Columns   []TableColumn
}

func NewCreateStatementBuilder(tableName string) *CreateStatementBuilder {
	return &CreateStatementBuilder{
		TableName: tableName,
		Columns:   make([]TableColumn, 0),
	}
}

/*
	Encodes column definition as SQL string
 */
func (it *CreateStatementBuilder) SqlString() string {
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

func (it *CreateStatementBuilder) WithColumnDefinition(name string, valueType string, notNull bool) *CreateStatementBuilder {
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
func (it *CreateStatementBuilder) AsPrimary() *CreateStatementBuilder {
	lastElementPos := len(it.Columns) - 1

	// ignore if no element set yet
	if lastElementPos < 0 {
		return it
	}

	lastElement := it.Columns[lastElementPos]
	lastElement.IsPrimary = true

	// need to put it back
	it.Columns = append(it.Columns[:lastElementPos-1], lastElement)

	return it
}

func (builder *CreateStatementBuilder) WithSerialColumn(name string, notNull bool) *CreateStatementBuilder {
	return builder.WithColumnDefinition(name, "SERIAL", notNull)
}

func (builder *CreateStatementBuilder) WithTextColumn(name string, notNull bool) *CreateStatementBuilder {
	return builder.WithColumnDefinition(name, "TEXT", notNull)
}

func (builder *CreateStatementBuilder) WithBooleanColumn(name string, notNull bool) *CreateStatementBuilder {
	return builder.WithColumnDefinition(name, "BOOLEAN", NotNull)
}

func (builder *CreateStatementBuilder) WithBigIntColumn(name string, notNull bool) *CreateStatementBuilder {
	return builder.WithColumnDefinition(name, "BIGINT", notNull)
}

func (builder *CreateStatementBuilder) WithIntColumn(name string, notNull bool) *CreateStatementBuilder {
	return builder.WithColumnDefinition(name, "INT", notNull)
}

func (builder *CreateStatementBuilder) WithByteAColumn(name string,  notNull bool) *CreateStatementBuilder {
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

