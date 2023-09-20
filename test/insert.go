package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
	"sync"
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

	sql, args, err := connect.Insert().
		Into("users").
		Columns("user_name", "day").
		Values("ghgh", time.Now().Format(time.DateOnly)).
		ToSql()
	fmt.Println("insert: ", sql, args, err)

	db19, id := connect.Insert().
		Into("users").
		Columns("user_name", "day", "is_delete", "create_time").
		Values("ghgh19", time.Now().Format(time.DateOnly), 1, time.Now().Format(time.DateTime)).
		LastInsertId()
	fmt.Println("db19: ", id, db19.Error)

	// 事务 https://learnku.com/docs/gorm/v2/transactions/9745
	// https://learnku.com/docs/gorm/v2/create/9732
	//users20 := &Users{UserName: "h20", Day: time.Now(), IsDelete: 1, CreateTime: time.Now()}
	//tx.Select("user_name", "day", "is_delete", "create_time").Create(&users20)
	//fmt.Println("db20: ", users20, tx.Error)

	db21, id := connect.Begin().Insert().
		Into("users").
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

	insert23, id := connect.Insert().
		Into("users").
		Set("user_name", "insert23").
		Set("day", time.Now().Format(time.DateOnly)).
		Set("create_time", time.Now().Format(time.DateTime)).
		LastInsertId()
	// INSERT INTO users (user_name,day,create_time) VALUES ('insert23','2023-09-20','2023-09-20 21:24:48')
	fmt.Println("insert23: ", id, insert23.Error)

	insert24, id := connect.Insert().
		Into("users").
		Columns("user_name", "day", "create_time").
		Values("ghgh24_1", time.Now().Format(time.DateOnly), time.Now().Format(time.DateTime)).
		Values("ghgh24_2", time.Now().Format(time.DateOnly), time.Now().Format(time.DateTime)).
		Save()
	// INSERT INTO users (user_name,day,create_time) VALUES ('ghgh24_1','2023-09-20','2023-09-20 21:24:48'),('ghgh24_2','2023-09-20','2023-09-20 21:24:48')
	// insert24:  2 <nil>
	fmt.Println("insert24: ", id, insert24.Error)

	insert25 := connect.Begin().Insert().
		Into("users").
		Set("user_name", "ghgh25_1").
		Set("day", time.Now().Format(time.DateOnly)).
		Set("create_time", time.Now().Format(time.DateTime))
	insert25DB, id := insert25.LastInsertId()
	// INSERT INTO users (user_name,day,create_time) VALUES ('ghgh25_1','2023-09-20','2023-09-20 21:36:52')
	// insert25:  141 <nil>
	fmt.Println("insert25: ", id, insert25DB.Error)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()
			insert26 := connect.Begin(insert25DB).Insert().
				Into("users").
				Set("user_name", fmt.Sprintf("ghgh26_%d", i)).
				Set("day", time.Now().Format(time.DateOnly)).
				Set("create_time", time.Now().Format(time.DateTime))
			if i%2 == 0 {
				time.Sleep(time.Duration(i) * time.Second)
				insert25DB.Error = fmt.Errorf("测试插入异常%d", i)
			}
			insert26DB, id := insert26.LastInsertId()
			fmt.Println(fmt.Sprintf("ghgh26_%d: ", i), id, insert26DB.Error)
		}(i)
	}
	wg.Wait()
	if insert25DB != nil {
		insert25DB.Rollback()
		fmt.Println("跨协程事务 Rollback", insert25DB.Error)
	} else {
		insert25DB.Commit()
		fmt.Println("跨协程事务 Commit", insert25DB.Error)
	}

	for true {

	}
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
