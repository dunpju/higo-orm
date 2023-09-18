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
	update1, affected := Transaction.Begin().Update().
		Table("users").
		Set("user_name", "user_name_1").
		Where("user_id", "=", 1).
		Exec()
	fmt.Println("update1: ", affected, fmt.Sprintf("%p", update1), update1.Error)

	update2, affected := Transaction.Begin(update1).
		Update().
		Table("users").
		Set("user_name", "user_name_update_2").
		Where("user_id", "=", 2).
		Exec()
	fmt.Println("update2: ", affected, fmt.Sprintf("%p", update2), update2.Error)
	update2.Rollback()
}
