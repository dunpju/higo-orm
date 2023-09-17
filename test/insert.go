package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/orm"
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
		Columns("user_name", "day", "create_time").
		Values("ghgh", time.Now().Format(time.DateOnly), time.Now().Format(time.DateTime)).
		LastInsertId()
	fmt.Println("db19: ", id, db19.Error)

	users20 := &Users{UserName: "h20", Day: time.Now(), IsDelete: 1, CreateTime: time.Now()}
	db20, _ := orm.Gorm()
	db20.Select("user_name", "day", "is_delete", "create_time").Create(&users20)
	fmt.Println("db20: ", users20, db20.Error)

	/*users21 := map[string]interface{}{
		"user_name":   "h21",
		"day":         time.Now(),
		"is_delete":   1,
		"create_time": time.Now(),
	}
	db21, _ := orm.Gorm()
	db20.Select("user_name", "day", "is_delete", "create_time").Create(&users21)
	fmt.Println("db21: ", users21, db21.Error)*/
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
