package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var (
	GDB *DB
)

// 直接使用sqllite3 存放
type DB struct {
	dbSourceName string
	db           *sql.DB
}

// 初始化db连接
func InitDB(dbSourceName string) (*DB, error) {
	dbDriverName := "sqlite3"
	db, err := sql.Open(dbDriverName, dbSourceName)
	if err != nil {
		fmt.Println("Failed to open the database:", err)
		return nil, err
	}

	// ping测试
	if err = db.Ping(); err != nil {
		fmt.Println("Failed to establish a connection to the database:", err)
	}

	// 检查 XuperChainAccount 表是否存在，没有则创建表
	sqlStatement := `create table if not exists XuperChainAccount (Address varchar(255), Mnemonic varchar(255));`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		fmt.Println("Failed to create table:", sqlStatement, err)
		return nil, err
	}

	return &DB{
		dbSourceName: dbSourceName,
		db:           db,
	}, nil
}

// 关闭连接
func (d *DB) Close() {
	d.db.Close()
}

////////////////////
// curd
///////////////////
func (d *DB) Add(address, mnemonic string) {
	// 开启事务
	tx, err := d.db.Begin()
	if err != nil {
		fmt.Println("Failed to start a database transaction:", err)
	}
	statement, err := tx.Prepare("insert into XuperChainAccount(Address, Mnemonic) values(?, ?)")
	if err != nil {
		fmt.Println("Failed to prepare SQL statements:", err)
	}
	defer statement.Close()
	_, err = statement.Exec(address, mnemonic)
	if err != nil {
		fmt.Println("Failed to prepare SQL statements:", err)
	}
	_ = tx.Commit()
}

func (d *DB) Query(address string) (string, error) {
	sqlStatement := `select Address, Mnemonic from XuperChainAccount where Address=$1`
	row := d.db.QueryRow(sqlStatement, address)

	var Address string
	var Mnemonic string
	err := row.Scan(&Address, &Mnemonic)
	if err != nil {
		return "", err
	}
	return Mnemonic, nil
}

func (d *DB) IsExist(address string) bool {
	sqlStatement := `select Address from XuperChainAccount where Address=$1`
	row := d.db.QueryRow(sqlStatement, address)
	var resultAddress string
	err := row.Scan(&resultAddress)
	if err != nil {
		return false
	}

	if address == resultAddress {
		return true
	}
	return false
}
