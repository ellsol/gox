package sqlx

import (
	"bytes"
)

var IsPrimary = true
var NotNull = true

type SQLTableColumn struct {
	Name      string
	Type      string
	IsPrimary bool
	NotNULL   bool
}

type SQLTableDefinition struct {
	TableName string
	Columns   []SQLTableColumn
}

func (definition *SQLTableDefinition) CreateStatement() string {
	var buffer bytes.Buffer

	buffer.WriteString("CREATE TABLE ")
	buffer.WriteString(definition.TableName)
	buffer.WriteString("(")

	for k, v := range definition.Columns {
		lastElement := k == len(definition.Columns)-1
		buffer.WriteString(v.Statement(!lastElement))
	}

	buffer.WriteString(");")
	return buffer.String()
}

func (definition *SQLTableDefinition) Tags() []string {
	tags := make([]string, 0)

	for _, v := range definition.Columns {
		if v.Type != "SERIAL" {
			tags = append(tags, v.Name)
		}
	}

	return tags
}

type SQLTableBuilder struct {
	Definition *SQLTableDefinition
}

func NewSQLTableBuilder(tableName string) *SQLTableBuilder {
	return &SQLTableBuilder{
		Definition: &SQLTableDefinition{
			TableName: tableName,
			Columns:   make([]SQLTableColumn, 0),
		},
	}
}

func (builder *SQLTableBuilder) Build() *SQLTableDefinition {
	return builder.Definition
}

func (builder *SQLTableBuilder) WithColumn(column *SQLTableColumn) *SQLTableBuilder {
	builder.Definition.Columns = append(builder.Definition.Columns, *column)
	return builder
}

func (builder *SQLTableBuilder) WithColumnDefinition(name string, valueType string, params ...bool) *SQLTableBuilder {
	notNull := false
	isPrimary := false

	if len(params) > 1 {
		isPrimary = true
	}

	if len(params) > 0 {
		notNull = true
	}

	col := &SQLTableColumn{
		Name:      name,
		IsPrimary: isPrimary,
		NotNULL:   notNull,
		Type:      valueType,
	}

	return builder.WithColumn(col)
}

func (builder *SQLTableBuilder) WithSerialColumn(name string, params ...bool) *SQLTableBuilder {
	return builder.WithColumnDefinition(name, "SERIAL", params...)
}

func (builder *SQLTableBuilder) WithTextColumn(name string, params ...bool) *SQLTableBuilder {
	return builder.WithColumnDefinition(name, "TEXT", params...)
}

func (builder *SQLTableBuilder) WithBooleanColumn(name string, params ...bool) *SQLTableBuilder {
	return builder.WithColumnDefinition(name, "BOOLEAN", params...)
}

func (builder *SQLTableBuilder) WithBigIntColumn(name string, params ...bool) *SQLTableBuilder {
	return builder.WithColumnDefinition(name, "BIGINT", params...)
}

func (builder *SQLTableBuilder) WithIntColumn(name string, params ...bool) *SQLTableBuilder {
	return builder.WithColumnDefinition(name, "INT", params...)
}

func (builder *SQLTableBuilder) WithByteAColumn(name string, params ...bool) *SQLTableBuilder {
	return builder.WithColumnDefinition(name, "BYTEA", params...)
}

func (column *SQLTableColumn) Statement(withComma bool) string {
	var buffer bytes.Buffer

	buffer.WriteString(column.Name)
	buffer.WriteString(" ")
	buffer.WriteString(column.Type)

	if column.IsPrimary {
		buffer.WriteString(" PRIMARY KEY")
	}

	if column.NotNULL {
		buffer.WriteString(" NOT NULL")
	}

	if withComma {
		buffer.WriteString(",")
	}

	return buffer.String()
}
