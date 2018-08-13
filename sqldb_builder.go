package gox

import "fmt"

type SQLConfig struct {
	DatabaseName string
	Schema       string
	User         string
	Password     string
	Host         string
	Tables       map[string]SQLTable
}

type SQLDBBuilder struct {
	SQLConfig *SQLConfig
}

func NewDatabaseBuilder(databaseName string) *SQLDBBuilder {
	return &SQLDBBuilder{SQLConfig: &SQLConfig{
		DatabaseName: databaseName,
		Tables:       make(map[string]SQLTable),
	}}
}

func (builder *SQLDBBuilder) WithHost(host string) *SQLDBBuilder {
	builder.SQLConfig.Host = host
	return builder
}

func (builder *SQLDBBuilder) WithSchema(schema string) *SQLDBBuilder {
	builder.SQLConfig.Schema = schema
	return builder
}

func (builder *SQLDBBuilder) WithUser(user string) *SQLDBBuilder {
	builder.SQLConfig.User = user
	return builder
}

func (builder *SQLDBBuilder) WithPassword(password string) *SQLDBBuilder {
	builder.SQLConfig.Password = password
	return builder
}

func (builder *SQLDBBuilder) AddTable(table SQLTable) *SQLDBBuilder {
	tables := builder.SQLConfig.Tables
	tables[table.TableName()] = table
	builder.SQLConfig.Tables = tables
	return builder
}

func (builder *SQLDBBuilder) OpenAndInitializeDB(forceRecreate bool) (*SQLDB, *SQLConfig, error) {
	config := builder.SQLConfig

	if config.Host == "" {
		return nil, nil, fmt.Errorf("no schema provided, use builder with .WithHost(...)")
	}

	if config.DatabaseName == "" {
		return nil, nil, fmt.Errorf("no database provided, use builder with NewSQLBuilder(...)")
	}

	if config.Schema == "" {
		return nil, nil, fmt.Errorf("no schema provided, use builder with .WithSchema(...)")
	}

	if config.User == "" {
		return nil, nil, fmt.Errorf("no user provided, use builder with .WithUser(...)")
	}

	if config.Password == "" {
		return nil, nil, fmt.Errorf("no password provided, use builder with .WithPassword(...)")
	}

	info := &SQLDBInfo{
		DBName:   "postgres",
		Password: config.Password,
		User:     config.User,
		Host:     config.Host,
	}

	db, err := OpenSqlDB(info)
	if err != nil {
		return nil, nil, err
	}

	err = db.MaybeCreateDatabase(config.DatabaseName)
	if err != nil {
		return nil, nil, err
	}

	// reopen connection with correct database
	err = db.DB.Close()
	if err != nil {
		return nil, nil, err
	}

	info.DBName = config.DatabaseName

	db, err = OpenSqlDB(info)
	if err != nil {
		return nil, nil, err
	}

	err = db.InitializeDatabase(config.DatabaseName, config.Schema, config.Tables, forceRecreate)
	if err != nil {
		return nil, nil, err
	}


	err = db.DB.Ping()
	if err != nil {
		return nil, nil, err
	}

	return db, builder.SQLConfig, nil
}
