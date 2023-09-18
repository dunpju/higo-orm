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

	/*gorm, _ := orm.Gorm()
	tx := gorm.Begin()
	fmt.Printf("%p\n", tx)*/
	db23, affected := Transaction.Begin().Update().
		Table("users").
		Set("user_name", "user_name_95").
		Where("user_id = ?", 1).
		Save()
	fmt.Println("db23: ", affected, fmt.Sprintf("%p", db23), db23.Error)

	db24, affected := Transaction.Begin(db23).
		Update().
		Table("users").
		Set("user_name", "user_name111").
		Where("user_id = ?", 2).
		Save()
	fmt.Println("db24: ", affected, fmt.Sprintf("%p", db24), db24.Error)
	db24.Rollback()
}
