package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
)

func main() {
	dbc := him.DbConfig(him.DefaultConnect).
		SetHost("192.168.8.99").
		SetPort("3306").
		SetDatabase("test").
		SetUsername("root").
		SetPassword("1qaz2wsx").
		SetCharset("utf8mb4").
		SetDriver("mysql").
		SetPrefix("tl_").
		SetMaxIdle(5).
		SetMaxOpen(10).
		SetMaxLifetime(1000).
		SetLogMode("Info").
		SetColorful(true)
	_, err := him.Init(dbc)
	if err != nil {
		panic(err)
	}

	connect, err := him.DBConnect(him.DefaultConnect)
	if err != nil {
		panic(err)
	}
	/*defer func() {
		if r := recover(); r != nil {
			fmt.Println("delete Rollback: ", r)
			tx.Rollback()
		}
	}()*/

	delete1, affected := connect.Begin().Delete().
		From("users").
		Where("user_id", "=", 1).
		Exec()
	fmt.Println("delete1: ", affected, fmt.Sprintf("%p", delete1), delete1.Error)

	delete2, affected := connect.Begin(delete1).
		Update().
		Table("users").
		Set("user_name", "user_name_delete111").
		Where("user_id", "=", 2).
		Exec()
	fmt.Println("delete2: ", affected, fmt.Sprintf("%p", delete2), delete2.Error)
	delete2.Rollback()
}
