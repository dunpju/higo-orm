package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
	"sync"
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
		SetColorful(true).
		Init()
	if err != nil {
		panic(err)
	}

	/*gorm, _ := orm.Gorm()
	tx := gorm.Begin()
	fmt.Printf("%p\n", tx)*/
	connect, err := him.DBConnect(him.DefaultConnect)
	if err != nil {
		panic(err)
	}

	update0, affected := connect.Update().
		Table("users").
		Set("user_name", "user_5").
		Where("user_id", "=", 5).
		Exec()
	fmt.Println("update0: ", affected, fmt.Sprintf("%p", update0), update0.Error)
	update1, affected := connect.Begin().Update().
		Table("users").
		Set("user_name", "user_6").
		Where("user_id", "=", 5).
		Exec()
	fmt.Println("update1: ", affected, fmt.Sprintf("%p", update1), update1.Error)
	update1.Rollback()

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			connect, err := him.DBConnect(him.DefaultConnect)
			if err != nil {
				panic(err)
			}
			update2, affected := connect.Begin().Update().
				Table("users").
				Set("user_name", "user_name_5").
				Where("user_id", "=", 5).
				Exec()
			fmt.Println("update2: ", affected, fmt.Sprintf("%p", update2), update2.Error)
			if update2.Error != nil {
				update2.Rollback()
				fmt.Println("update2 Rollback")
			} else {
				update2.Commit()
			}
		}()
		go func() {
			wg.Add(1)
			defer wg.Done()
			connect, err := him.DBConnect(him.DefaultConnect)
			if err != nil {
				panic(err)
			}
			/*update3, affected := connect.Begin().Update().
				Table("users").
				Set("user_name", "user_name_5").
				Where("user_id", "=", 5).
				Exec()
			fmt.Println("update3: ", affected, fmt.Sprintf("%p", update3), update3.Error)*/

			update4, affected := connect.Begin().
				Update().
				Table("users").
				Set("user_name", "user_name_update_2").
				Where("user_id", "=", 2).
				Exec()
			// UPDATE users SET user_name = 'user_name_update_2' WHERE (user_id = 2)
			fmt.Println("update4: ", affected, fmt.Sprintf("%p", update4), update4.Error)
			if update4.Error != nil {
				update4.Rollback()
				fmt.Println("update4 Rollback")
			} else {
				update4.Commit()
			}
		}()
	}

	connect1, err := him.DBConnect(him.DefaultConnect)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()
			update5 := connect1.Begin().GormDB()
			update := connect1.TX(update5).Update().Table("users")
			update = update.Where("user_id", "=", 2)
			update = update.Set("user_name", fmt.Sprintf("update5_%d", i))

			update6, affected := connect1.Begin(update5).Update().
				Table("users").
				Set("user_name", fmt.Sprintf("update6_%d", i)).
				Where("user_id", "=", 7).
				Exec()

			// sql交叉测试
			update5, affected5 := update.Exec()
			// UPDATE users SET user_name = 'update5_9' WHERE (user_id = 2)
			fmt.Println("update5: ", affected5, fmt.Sprintf("%p", update5), update5.Error)
			if i%2 == 0 {
				update5.Error = fmt.Errorf("测试异常回滚%d", i)
			}

			// UPDATE users SET user_name = 'update6_9' WHERE (user_id = 2)
			fmt.Println("update6: ", affected, fmt.Sprintf("%p", update6), update6.Error)
			if update6.Error != nil {
				update6.Rollback()
				fmt.Println("update6 Rollback", update6.Error)
			} else {
				update6.Commit()
				fmt.Println("update6 Commit")
			}
		}(i)
	}
	wg.Wait()

}
