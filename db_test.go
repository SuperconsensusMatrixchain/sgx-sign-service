package main

import (
	"fmt"
	"testing"
)

// 测试

func Test_DB(t *testing.T) {

	dbsoucename := "test_db"

	db, err := InitDB(dbsoucename)
	if err != nil {
		return
	}
	defer db.Close()
	db.Add("acc1", "m1")
	db.Add("acc2", "m2")

	result, err := db.Query("acc1")
	if err != nil {
		fmt.Println(result)
	}
	fmt.Println(result)

	result, err = db.Query("acc2")
	if err != nil {
		fmt.Println(result)
	}
	fmt.Println(result)
}
