package sqlx

import "fmt"

type DatabaseCreator struct {
	DatabaseName string
	Schema       string
	User         string
	Password     string
	Host         string
	Tables       map[string]SQLTable
}

func (it *DatabaseCreator) dbinfo() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", it.Host, it.User, it.Password, it.DatabaseName)
}

func (it *DatabaseCreator) dbDefaultInfo() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", it.Host, it.User, it.Password)
}

func NewDatabaseCreator(databaseName string) *DatabaseCreator {
	return &DatabaseCreator{
		DatabaseName: databaseName,
		Tables:       make(map[string]SQLTable),
	}
}

func (builder *DatabaseCreator) WithHost(host string) *DatabaseCreator {
	builder.Host = host
	return builder
}

func (it *DatabaseCreator) WithSchema(schema string) *DatabaseCreator {
	it.Schema = schema
	return it
}

func (it *DatabaseCreator) WithUser(user string) *DatabaseCreator {
	it.User = user
	return it
}

func (it *DatabaseCreator) WithPassword(password string) *DatabaseCreator {
	it.Password = password
	return it
}

func (it *DatabaseCreator) AddTable(table SQLTable) *DatabaseCreator {
	tables := it.Tables
	tables[table.Name()] = table
	it.Tables = tables
	return it
}

func (it *DatabaseCreator) OpenAndInitializeDB(forceRecreate bool) (*SQLDB, error) {
	config := it

	if config.Host == "" {
		return nil, fmt.Errorf("no schema provided, use creator with .WithHost(...)")
	}

	if config.DatabaseName == "" {
		return nil, fmt.Errorf("no database provided, use creator with NewDatabaseCreator(...)")
	}

	if config.Schema == "" {
		return nil, fmt.Errorf("no schema provided, use creator with .WithSchema(...)")
	}

	if config.User == "" {
		return nil, fmt.Errorf("no user provided, use creator with .WithUser(...)")
	}

	if config.Password == "" {
		return nil, fmt.Errorf("no password provided, use creator with .WithPassword(...)")
	}

	db, err := OpenSqlDB(it.dbDefaultInfo())
	if err != nil {
		return nil, err
	}

	err = db.MaybeCreateDatabase(config.DatabaseName)
	if err != nil {
		return nil, err
	}

	// reopen connection with correct database
	err = db.DB.Close()
	if err != nil {
		return nil, err
	}

	db, err = OpenSqlDB(it.dbinfo())
	if err != nil {
		return nil, err
	}

	err = db.InitializeDatabase(config.DatabaseName, config.Schema, config.Tables, forceRecreate)
	if err != nil {
		return nil, err
	}

	err = db.DB.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
