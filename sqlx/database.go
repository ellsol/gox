package sqlx

import (
	"strings"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ellsol/gox/typesx"
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
	DB sql.DB
}

type SQLTable interface {
	ColumnNames() []string
	Name() string
	KeyTag() string
	CreateStatement() string
}

func OpenSqlDB(params string) (*SQLDB, error) {
	fmt.Println("Trying to open connection to postgres with: ", info.dbinfo())
	db, err := sql.Open("postgres", params)

	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(20)

	return &SQLDB{
		DB: *db,
	}, nil
}


func (db *SQLDB) InitializeDatabase(databaseName string, schema string, tables map[string]SQLTable, forceRecreate bool) error {
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

func (pg *SQLDB) MaybeCreateDatabase(database string) error {
	statement := fmt.Sprintf(CreateDatabaseStatement, database)

	stmt, err := pg.DB.Prepare(statement)
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

func (pg *SQLDB) DropDatabaseIfExist(database string) (error) {
	statement := fmt.Sprintf(DropDatabaseStatement, database)
	stmt, err := pg.DB.Prepare(statement)
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (pg *SQLDB) MaybeCreateScheme(scheme string) error {
	statement := fmt.Sprintf(CreateSchemaStatement, scheme)

	stmt, err := pg.DB.Prepare(statement)
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

func (pg *SQLDB) DropSchemaIfExist(schema string) (error) {
	statement := fmt.Sprintf(DropSchemaStatement, schema)
	stmt, err := pg.DB.Prepare(statement)
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (pg *SQLDB) MaybeCreateTable(table SQLTable) (error) {
	stmt, err := pg.DB.Prepare(table.CreateStatement())
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

func (pg *SQLDB) DropTableIfExist(table SQLTable) (error) {
	statement := fmt.Sprintf(DropTableStatement, table.Name())
	stmt, err := pg.DB.Prepare(statement)
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
// Statements nobody needs abstracted but sometimes helpful
//
/////////////////////////////////////////////////////////////////

func (pg *SQLDB) Insert(table SQLTable, values []interface{}) (int, error) {
	statement := GetPostgresInsertStatementNoIncrement(table)
	o, err := pg.DB.Query(statement, values...)
	if err != nil {
		return -1, err
	}
	var lastInsertId int
	o.Scan(&lastInsertId)
	o.Close()

	return lastInsertId, nil
}

func (pg *SQLDB) Update(table SQLTable, values []interface{}) error {
	statement := CreateUpdateStatement(table)
	return pg.UpdateWithStatement(statement, table, values)
}


func (pg *SQLDB) UpdateWithStatement(statement string, table SQLTable, values []interface{}) error {
	updated, err := pg.DB.Exec(statement, values...)

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
func (pg *SQLDB) Delete(key interface{}, table SQLTable) error {
	sqlStatement := fmt.Sprintf(DeleteStatement, table.Name(), table.KeyTag())
	_, err := pg.DB.Exec(sqlStatement, key)
	if err != nil {
		return err
	}

	return nil
}

// Number Of Rows
func (pg *SQLDB) Count(table SQLTable) (int, error) {
	sqlStatement := fmt.Sprintf(NumberOfRowsStatement, table.Name())
	rows, err := pg.DB.Query(sqlStatement)
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

// Number Of Rows
func (pg *SQLDB) Max(table SQLTable, column string) (int64, error) {
	sqlStatement := fmt.Sprintf(MaxStatement, column, table.Name())
	rows, err := pg.DB.Query(sqlStatement)
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
	paramsJoin := typesx.CommaSeparatedString(t.ColumnNames())
	paramsPlaceholder := typesx.CommaSeparatedString(typesx.MapStringListWithPos(t.ColumnNames(), func(key int, value string) string {
		return fmt.Sprintf("$%v", key+1)
	}))

	return fmt.Sprintf(InsertStatement, t.Name(), paramsJoin, paramsPlaceholder)
}

/*
	 Maps SQLTable to update statement
 */
func CreateUpdateStatement(table SQLTable) string {
	set := typesx.MapStringListWithPos(table.ColumnNames()[1:], func(pos int, tag string) string {
		return fmt.Sprintf("%v = $%v", tag, pos+2)
	})

	return fmt.Sprintf("UPDATE %v SET %v WHERE %v = $1;", table.Name(), typesx.CommaSeparatedString(set), table.KeyTag())
}

func CreateGetStatement(table SQLTable) string {
	return fmt.Sprintf("SELECT * FROM %v WHERE %v= $1", table.Name(), table.KeyTag())
}