package sqlx

import (
	"database/sql"
	"fmt"
	"github.com/ellsol/gox/typex"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

const (
	CreateDatabaseStatement   = "CREATE DATABASE %v;"
	DropDatabaseStatement     = "DROP DATABASE IF EXISTS %v;"
	CreateSchemaStatement     = "CREATE SCHEMA %v;"
	DropSchemaStatement       = "DROP SCHEMA IF EXISTS %v CASCADE;"
	CreateTableStatement      = "CREATE TABLE %v %v;"
	DropTableStatement        = "DROP TABLE IF EXISTS %v;"
	DeleteStatement           = "DELETE FROM %v WHERE %v = $1;"
	InsertStatementWithReturn = "INSERT INTO %v(%v) VALUES(%v) returning %v;"
	InsertStatement           = "INSERT INTO %v(%v) VALUES(%v);"
	NumberOfRowsStatement     = "SELECT count(*) FROM %v;"
	MaxStatement              = "SELECT max(%v) FROM %v;"
)

type SQLDB struct {
	Connection *sql.DB
}

type SqlDBInfo struct {
	Host     string
	User     string
	Password string
	DBName   string
}

func (pi *SqlDBInfo) dbinfo() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", pi.Host, pi.User, pi.Password, pi.DBName)
}

type SQLTable interface {
	ColumnNames() []string
	Name() string
	CreateStatement() string
}

func OpenSqlDB(params string) (*SQLDB, error) {
	fmt.Println("Trying to open connection to postgres with: ", params)
	connection, err := sql.Open("postgres", params)

	if err != nil {
		return nil, err
	}

	connection.SetMaxIdleConns(20)

	return &SQLDB{
		Connection: connection,
	}, nil
}

func (db *SQLDB) InitializeDatabase(databaseName string, schema string, tables map[string]SQLTable, forceRecreate bool) error {
	logMsg(fmt.Sprintf("initializing db %v with scheme %v and forceRecreate: %v", databaseName, schema, forceRecreate))
	if forceRecreate {
		err := db.DropSchemaIfExist(schema)
		if err != nil {
			return err
		}
	}

	err := db.MaybeCreateScheme(schema)
	if err != nil {
		return err
	}

	err = db.MaybeInitializeTables(tables)
	if err != nil {
		return err
	}

	return nil
}

func (it *SQLDB) MaybeCreateDatabase(database string) error {
	statement := fmt.Sprintf(CreateDatabaseStatement, database)

	stmt, err := it.Connection.Prepare(statement)
	if err != nil {
		return nil
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {

		if strings.Contains(err.Error(), "already exists") {
			return nil
		}

		return err
	}
	return nil
}

func (it *SQLDB) DropDatabaseIfExist(database string) (error) {
	statement := fmt.Sprintf(DropDatabaseStatement, database)
	stmt, err := it.Connection.Prepare(statement)
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (it *SQLDB) MaybeCreateScheme(scheme string) error {
	logMsg(fmt.Sprintf("Maybe create schema %v", scheme))
	statement := fmt.Sprintf(CreateSchemaStatement, scheme)
	logMsg(fmt.Sprintf("Maybe create schema statement: %v", statement))
	stmt, err := it.Connection.Prepare(statement)
	if err != nil {
		return nil
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		logMsg(err.Error())
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}

		return err
	}
	return nil
}

func (it *SQLDB) DropSchemaIfExist(schema string) (error) {
	logMsg(fmt.Sprintf("Dropping schema %v", schema))
	statement := fmt.Sprintf(DropSchemaStatement, schema)
	logMsg(fmt.Sprintf("Dropping schema statement: %v", statement))
	stmt, err := it.Connection.Prepare(statement)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec()

	return err
}

func (it *SQLDB) MaybeCreateTable(table SQLTable) (error) {
	logMsg(table.CreateStatement())
	stmt, err := it.Connection.Prepare(table.CreateStatement())
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}

		return err
	}
	return nil
}

func (it *SQLDB) DropTableIfExist(table SQLTable) (error) {
	statement := fmt.Sprintf(DropTableStatement, table.Name())
	stmt, err := it.Connection.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (db *SQLDB) MaybeInitializeTables(tables map[string]SQLTable) error {
	for _, v := range tables {
		err := db.MaybeCreateTable(v)

		if err != nil {
			return err
		}
	}

	return nil
}

/////////////////////////////////////////////////////////////////
//
// Statements nobody needs, abstracted but sometimes helpful
//
/////////////////////////////////////////////////////////////////

func (pg *SQLDB) Insert(table SQLTable, values []interface{}) (int, error) {
	statement := GetPostgresInsertStatementNoIncrement(table)
	o, err := pg.Connection.Query(statement, values...)
	if err != nil {
		return -1, err
	}
	var lastInsertId int
	o.Scan(&lastInsertId)
	o.Close()

	return lastInsertId, nil
}

func (pg *SQLDB) InsertOmitPrimary(table SQLTable, values []interface{}) (int, error) {
	statement := GetPostgresInsertStatementNoIncrementOmitPrimary(table) 
	o, err := pg.Connection.Query(statement, values...)
	if err != nil {
		return -1, err
	}
	var lastInsertId int
	o.Scan(&lastInsertId)
	o.Close()

	return lastInsertId, nil
}

func (pg *SQLDB) Update(table SQLTable, keyLabel string, values []interface{}) error {
	statement := CreateUpdateStatement(table, keyLabel)
	return pg.UpdateWithStatement(statement, table, values)
}

func (pg *SQLDB) UpdateWithStatement(statement string, table SQLTable, values []interface{}) error {
	updated, err := pg.Connection.Exec(statement, values...)

	if err != nil {
		return err
	}

	count, err := updated.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return fmt.Errorf("failed to update %v", table.Name())
	}

	return nil
}

// Delete Row
func (pg *SQLDB) Delete(key interface{}, keyLabel string, table SQLTable) error {
	sqlStatement := fmt.Sprintf(DeleteStatement, table.Name(), keyLabel)
	_, err := pg.Connection.Exec(sqlStatement, key)
	if err != nil {
		return err
	}

	return nil
}

// Number Of Rows
func (pg *SQLDB) Count(table SQLTable) (int, error) {
	sqlStatement := fmt.Sprintf(NumberOfRowsStatement, table.Name())
	rows, err := pg.Connection.Query(sqlStatement)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	if rows.Next() {
		var count int
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
		return count, nil
	}
	return -1, nil
}

func (it *SQLDB) CountByStatement(table SQLTable, statement string, params ... interface{}) (int, error) {

	var count int
	row := it.Connection.QueryRow(statement, params...)
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

// Number Of Rows
func (pg *SQLDB) Max(table SQLTable, column string) (int64, error) {
	sqlStatement := fmt.Sprintf(MaxStatement, column, table.Name())
	rows, err := pg.Connection.Query(sqlStatement)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	if rows.Next() {
		var max int64

		err := rows.Scan(&max)

		if err != nil {
			return 0, nil
		}

		return max, nil
	}
	return -1, nil
}

// Statements

/*
	Transforms a table a model key list [tag1, tag2,...] into and a returningStatement
	INSERT INTO table(tag1, tag2,...) VALUES(1,2,...) returning returningStatement
 */
func GetPostgresInsertStatementNoIncrement(t SQLTable) string {
	paramsJoin := typex.CommaSeparatedString(t.ColumnNames())
	paramsPlaceholder := typex.CommaSeparatedString(typex.MapStringListWithPos(t.ColumnNames(), func(key int, value string) string {
		return fmt.Sprintf("$%v", key+1)
	}))

	return fmt.Sprintf(InsertStatement, t.Name(), paramsJoin, paramsPlaceholder)
}

func GetPostgresInsertStatementNoIncrementOmitPrimary(t SQLTable) string {
	paramsJoin := typex.CommaSeparatedString(t.ColumnNames()[1:])
	paramsPlaceholder := typex.CommaSeparatedString(typex.MapStringListWithPos(t.ColumnNames()[1:], func(key int, value string) string {
		return fmt.Sprintf("$%v", key+1)
	}))

	return fmt.Sprintf(InsertStatement, t.Name(), paramsJoin, paramsPlaceholder)
}

/*
	 Maps SQLTable to update statement
 */
func CreateUpdateStatement(table SQLTable, keyLabel string) string {
	set := typex.MapStringListWithPos(table.ColumnNames()[1:], func(pos int, tag string) string {
		return fmt.Sprintf("%v = $%v", tag, pos+2)
	})

	return fmt.Sprintf("UPDATE %v SET %v WHERE %v = $1;", table.Name(), typex.CommaSeparatedString(set), keyLabel)
}

var LogDatabase bool = true

func logMsg(msg string) {
	if(LogDatabase) {
		log.Println(msg)
	}
}