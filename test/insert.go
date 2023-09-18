package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/orm"
	"github.com/dunpju/higo-orm/orm/Transaction"
	"time"
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

	sql, args, err := orm.Insert("users").
		Columns("user_name", "day").
		Values("ghgh", time.Now().Format(time.DateOnly)).
		ToSql()
	fmt.Println("Insert: ", sql, args, err)

	db19, id := orm.Insert("users").
		Columns("user_name", "day", "is_delete", "create_time").
		Values("ghgh19", time.Now().Format(time.DateOnly), 1, time.Now().Format(time.DateTime)).
		LastInsertId()
	fmt.Println("db19: ", id, db19.Error)

	db20, _ := orm.Gorm()
	// 事务 https://learnku.com/docs/gorm/v2/transactions/9745
	tx := db20.Begin()
	// https://learnku.com/docs/gorm/v2/create/9732
	//users20 := &Users{UserName: "h20", Day: time.Now(), IsDelete: 1, CreateTime: time.Now()}
	//tx.Select("user_name", "day", "is_delete", "create_time").Create(&users20)
	//fmt.Println("db20: ", users20, tx.Error)

	db21, id := Transaction.Begin(tx).Insert("users").
		Columns("user_name", "day", "create_time").
		Values("ghgh21", time.Now().Format(time.DateOnly), time.Now().Format(time.DateTime)).
		LastInsertId()
	fmt.Println("db21: ", id, db21.Error)

	db22, affected := Transaction.Begin(tx).
		Update().
		Table("users").
		Set("user_name", "user_name_98").
		Where("user_id = ?", 2).
		Save()
	fmt.Println("db22: ", affected, fmt.Sprintf("%p", db22), db22.Error)
	db21.Rollback()

}

type Users struct {
	ID         int64
	UserId     int64
	UserName   string
	Day        time.Time
	IsDelete   int8
	CreateTime time.Time
	UpdateTime time.Time
}
