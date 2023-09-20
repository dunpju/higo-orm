package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
	"time"
)

func main() {
	_, err := him.DbConfig(him.DefaultConnect).
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
		SetSlowThreshold(1).
		SetColorful(true).
		Init()
	if err != nil {
		panic(err)
	}

	connect, err := him.DBConnect(him.DefaultConnect)
	if err != nil {
		panic(err)
	}

	sql, args, err := connect.Insert("users").
		Columns("user_name", "day").
		Values("ghgh", time.Now().Format(time.DateOnly)).
		ToSql()
	fmt.Println("insert: ", sql, args, err)

	db19, id := connect.Insert("users").
		Columns("user_name", "day", "is_delete", "create_time").
		Values("ghgh19", time.Now().Format(time.DateOnly), 1, time.Now().Format(time.DateTime)).
		LastInsertId()
	fmt.Println("db19: ", id, db19.Error)

	// 事务 https://learnku.com/docs/gorm/v2/transactions/9745
	// https://learnku.com/docs/gorm/v2/create/9732
	//users20 := &Users{UserName: "h20", Day: time.Now(), IsDelete: 1, CreateTime: time.Now()}
	//tx.Select("user_name", "day", "is_delete", "create_time").Create(&users20)
	//fmt.Println("db20: ", users20, tx.Error)

	db21, id := connect.Begin().Insert("users").
		Columns("user_name", "day", "create_time").
		Values("ghgh21", time.Now().Format(time.DateOnly), time.Now().Format(time.DateTime)).
		LastInsertId()
	fmt.Println("db21: ", id, db21.Error)

	db22, affected := connect.TX(db21).
		Update().
		Table("users").
		Set("user_name", "user_name_98").
		Where("user_id", "=", 2).
		Exec()
	fmt.Println("db22: ", affected, fmt.Sprintf("%p", db22), db22.Error)
	db22.Rollback()

	insert23, id := connect.Insert("users").
		Set("user_name", "insert23").
		Set("day", time.Now().Format(time.DateOnly)).
		Set("create_time", time.Now().Format(time.DateTime)).
		LastInsertId()
	fmt.Println("insert23: ", id, insert23.Error)

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
