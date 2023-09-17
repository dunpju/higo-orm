package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/orm"
	"github.com/dunpju/higo-orm/orm/Transaction"
)

func main() {
	orm.DbConfig().
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
	_, err := orm.Init()
	if err != nil {
		panic(err)
	}

	tx, _ := orm.Gorm()
	fmt.Println(tx)
	db22, affected := Transaction.Begin(tx).Update().
		Table("users").
		Set("user_name", "user_name_95").
		Where("user_id = ?", 95).
		Save()
	fmt.Println("db22: ", affected, db22, db22.Error)

	db23, affected := Transaction.Begin(tx).
		Update().
		Table("users").
		Set("user_name", "user_name_98").
		Where("user_id = ?", 98).
		Save()
	fmt.Println("db23: ", affected, db23, db23.Error)
	db23.Rollback()
}
