package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
	"sync"
	"time"
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
		SetMaxIdle(3).
		SetMaxOpen(5).
		SetMaxLifetime(60).
		SetLogMode("Info").
		SetColorful(true)
	db, err := him.Init(dbc)
	if err != nil {
		panic(err)
	}

	connect, err := him.DBConnect(him.DefaultConnect)
	if err != nil {
		panic(err)
	}

	users66 := make([]map[string]interface{}, 0)
	select66 := connect.Query().Select("user_id", "user_name").
		From("users")
	select66.Where("user_id", "=", 8)
	select66.First(&users66)
	// SELECT user_id, user_name FROM users WHERE (user_id = 8) LIMIT 1
	fmt.Println("users66-1:", users66)
	select66 = select66.Query().Select("user_id").From("users")
	select66.Where("user_id", "=", 9)
	select66.First(&users66)
	fmt.Println("users66-2:", users66)
	return

	userNames := make([]string, 0)
	userNames = append(userNames, "ggg")
	userNames = append(userNames, "ttttt")
	sql, args, err := connect.Query().Select("*").
		From("users").
		WhereIn("user_name", userNames).
		OrWhere("is_delete", "=", 1).
		WhereNull("update_time").
		ToSql()
	// SELECT * FROM users WHERE user_name IN (?,?) OR is_delete = ? AND update_time IS NULL [ggg ttttt 1] <nil>
	fmt.Println("users1:", sql, args, err)
	if err != nil {
		panic(err)
	}

	users1 := make([]map[string]interface{}, 0)
	// SELECT * FROM users WHERE user_name IN ('ggg','ttttt') OR is_delete = 1 AND update_time IS NULL
	db.Raw(sql, args...).Scan(&users1)
	fmt.Println(users1)

	users2 := make([]map[string]interface{}, 0)
	sql, args, err = connect.Query().Select("*").
		From("users").
		WhereBetween("day", "2023-06-11", "2023-06-12").
		ToSql()
	fmt.Println("users2:", sql, args, err)
	db.Raw(sql, args...).Scan(&users2)
	fmt.Println(users2)

	users3 := make([]map[string]interface{}, 0)
	sql, args, err = connect.Query().Select("*").
		From("users").
		WhereRaw(func(builder him.WhereRawBuilder) him.WhereRawBuilder {
			return builder.Where("user_id", "=", 3).OrWhere("user_id", "=", 5)
		}).
		ToSql()
	// SELECT * FROM users WHERE ((user_id = ?) OR (user_id = ?)) [3 5] <nil>
	fmt.Println("users3:", sql, args, err)
	// SELECT * FROM users WHERE ((user_id = 3) OR (user_id = 5))
	db.Raw(sql, args...).Scan(&users3)
	fmt.Println(users3)

	users4 := make([]map[string]interface{}, 0)
	sql, args, err = connect.Query().Select("*").
		From("users").
		Where("user_id", "=", 4).
		OrWhereRaw(func(builder him.WhereRawBuilder) him.WhereRawBuilder {
			// return builder.Where("user_id", "=", 3).Where("user_id", "=", 5)
			userIds := make([]int64, 0)
			userIds = append(userIds, 2)
			userIds = append(userIds, 3)
			b := builder.WhereIn("user_id", userIds)
			//b = b.Where("user_id", "=", 3)
			//b = b.Where("user_id", "=", 5)
			b = b.OrWhere("user_id", "=", 1)
			return b
		}).
		ToSql()
	// SELECT * FROM users WHERE (user_id = ?) OR ((user_id = ?) AND (user_id = ?)) [4 3 5] <nil>
	// SELECT * FROM users WHERE (user_id = ?) OR ((user_id IN (?,?)) AND (user_id = ?) AND (user_id = ?)) [4 2 3 3 5] <nil>
	// SELECT * FROM users WHERE (user_id = ?) OR ((user_id IN (?,?)) OR (user_id = ?)) [4 2 3 1] <nil>
	fmt.Println("users4:", sql, args, err)
	// SELECT * FROM users WHERE (user_id = 4) OR ((user_id = 3) AND (user_id = 5))
	db.Raw(sql, args...).Scan(&users4)
	fmt.Println(users4)

	users5 := make([]map[string]interface{}, 0)
	sql, args, err = connect.Query().Select("user_id", "user_name", "day").
		From("users").
		Where("user_id", "=", 4).
		ToSql()
	// SELECT * FROM users WHERE (user_id = ?) [4] <nil>
	fmt.Println("users5:", sql, args, err)
	// SELECT * FROM users WHERE (user_id = 4)
	db.Raw(sql, args...).Scan(&users5)
	fmt.Println(users5)

	users6 := make([]map[string]interface{}, 0)
	select6 := connect.Query().Select("user_id", "user_name").
		From("users")
	select6.Where("user_id", "=", 8)
	select6.First(&users6)
	// SELECT user_id, user_name FROM users WHERE (user_id = 8) LIMIT 1
	fmt.Println("users6-1:", users6)
	select6.Query().Select("user_id").From("users")
	select6.Where("user_id", "=", 9)
	select6.First(&users6)
	fmt.Println("users6-2:", users6)

	users7 := make([]map[string]interface{}, 0)
	db7 := connect.Query().Select("user_id", "user_name").
		From("users1").
		Where("user_id", "=", 8).
		First(&users7)
	// SELECT user_id, user_name FROM users1 WHERE (user_id = 8) LIMIT 1
	fmt.Println("users7:", users7)
	// Error 1146 (42S02): Table 'test.users1' doesn't exist
	fmt.Println(db7.Error)

	users8 := make([]map[string]interface{}, 0)
	db8 := connect.Query().Select("user_id", "user_name").
		From("users").
		Where("user_id", "=", 7).
		First(&users8)
	// SELECT user_id, user_name FROM users WHERE (user_id = 7) LIMIT 1
	fmt.Println(users8)
	fmt.Println(db8.Error) // <nil>

	users9 := make([]map[string]interface{}, 0)
	db9, paginate := connect.Query().Select("user_id", "user_name").
		From("users").
		Where("user_name", "=", "kkk").
		Paginate(2, 2, &users9)
	// SELECT user_id, user_name FROM users LIMIT 2 OFFSET 0    {8 2 1 0 0xc0002aabe8}
	// SELECT user_id, user_name FROM users WHERE (user_name = 'kkk') LIMIT 2 OFFSET 0    {4 2 1 0 0xc0002aabe8}
	// SELECT user_id, user_name FROM users WHERE (user_name = 'kkk') LIMIT 2 OFFSET 2    {4 2 1 0 0xc0002aabe8}
	fmt.Println(users9, paginate)
	fmt.Println(db9.Error) // <nil>

	select10, count := connect.Query().
		Select("count(distinct(user_name))").
		From("users").
		// Where("user_name", "=", "kkk").
		GroupBy("user_name").
		Count()
	// SELECT count(*) FROM `users` WHERE (user_name = 'kkk')
	// SELECT count(distinct(user_name)) FROM `users`
	// SELECT count(distinct(user_name)) FROM `users` GROUP BY `user_name`
	fmt.Println("select10: ", select10, count)

	users11 := make([]map[string]interface{}, 0)
	db11 := connect.Query().
		//Select("count(distinct(user_name)) count", "user_name").
		Select("count(user_name) count", "user_name").
		From("users").
		// Where("user_name", "=", "kkk").
		GroupBy("user_name").
		OrderBy("count desc").
		//Having("count > ?", 2).
		//Having("count > ? AND count <= 4", 2).
		Having("count > ?", 2).
		Having("count <= ?", 4).
		Get(&users11)
	// SELECT count(distinct(user_name)) count, user_name FROM users GROUP BY user_name
	// SELECT count(user_name) count, user_name FROM users GROUP BY user_name
	// SELECT count(user_name) count, user_name FROM users GROUP BY user_name ORDER BY count desc
	// SELECT count(user_name) count, user_name FROM users GROUP BY user_name HAVING count >= (2) ORDER BY count desc
	// SELECT count(user_name) count, user_name FROM users GROUP BY user_name HAVING count > (2) AND count <= 4 ORDER BY count desc
	// SELECT count(user_name) count, user_name FROM users GROUP BY user_name HAVING count > (2) AND count <= (4) ORDER BY count desc
	fmt.Println("db11: ", users11)
	fmt.Println(db11.Error) // <nil>

	select12, sum := connect.Query().
		Select("*").
		From("users").
		Where("user_name", "=", "jjj").
		Sum("is_delete")
	// SELECT SUM(is_delete) count_ FROM users LIMIT 1
	// SELECT SUM(is_delete) count_ FROM users WHERE (user_name = 'jjj') LIMIT 1
	fmt.Println("select12: ", select12, sum)

	users13 := make([]map[string]interface{}, 0)
	db13 := connect.Query().Raw("SELECT * FROM users").
		Get(&users13)
	// SELECT * FROM users
	// SELECT SUM(is_delete) count_ FROM users WHERE (user_name = 'jjj') LIMIT 1
	fmt.Println("users13: ", users13)
	fmt.Println(db13.Error) // <nil>

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			connect14, err := him.DBConnect(him.DefaultConnect)
			if err != nil {
				panic(err)
			}
			users14 := make([]map[string]interface{}, 0)
			db14 := connect14.Query().Select("*").
				From("users").
				WhereRaw(func(builder him.WhereRawBuilder) him.WhereRawBuilder {
					return builder.Raw("user_id = ?", 1)
				}).
				Get(&users14)
			// SELECT * FROM users WHERE (user_id = 1)
			fmt.Println("users14: ", users14)
			fmt.Println(db14.Error) // <nil>
		}()
		go func() {
			wg.Add(1)
			defer wg.Done()
			connect15, err := him.DBConnect(him.DefaultConnect)
			if err != nil {
				panic(err)
			}
			users15 := make([]map[string]interface{}, 0)
			db15 := connect15.Query().Select("*").
				From("users").
				OrderBy("user_id desc").
				First(&users15)
			// SELECT * FROM users ORDER BY user_id desc LIMIT 1
			fmt.Println("users15: ", users15)
			fmt.Println(db15.Error) // <nil>
		}()
		go func(i int) {
			wg.Add(1)
			defer wg.Done()
			connect16, err := him.DBConnect(him.DefaultConnect)
			if err != nil {
				panic(err)
			}

			users16 := make([]map[string]interface{}, 0)
			db16 := connect16.Query().
				Distinct().
				Select("user_name").
				From("users").
				OrderBy("user_id desc").
				Get(&users16)
			// SELECT DISTINCT user_name FROM users ORDER BY user_id desc
			fmt.Println("users16: ", users16)
			fmt.Println(db16.Error) // <nil>
			if i%2 == 0 {
				time.Sleep(time.Duration(i) * time.Second)
			}

			users17 := make([]map[string]interface{}, 0)
			db17 := connect16.Query().
				Select("*").
				From("users AS A").
				Join("ts_user AS B", "B.uname", "=", "A.user_name").
				OrderBy("A.user_id desc").
				Get(&users17)
			// SELECT * FROM users AS A JOIN ts_user AS B ON B.uname = A.user_name ORDER BY A.user_id desc
			fmt.Println("users17: ", users17)
			fmt.Println(db17.Error) // <nil>

			users18 := make([]map[string]interface{}, 0)
			db18 := connect16.Query().
				Select("*").
				From("users AS A").
				InnerJoin("ts_user AS B", "B.uname", "=", "A.user_name").
				OrderBy("A.user_id desc").
				Get(&users18)
			// SELECT * FROM users AS A INNER JOIN ts_user AS B ON B.uname = A.user_name ORDER BY A.user_id desc
			fmt.Println("users18: ", users18)
			fmt.Println(db18.Error) // <nil>

			users19 := make([]map[string]interface{}, 0)
			db19 := connect16.Query().
				Select("*").
				From("ts_user AS A").
				LeftJoin("users AS B", "B.user_name", "=", "A.uname").
				OrderBy("B.user_id desc").
				Get(&users19)
			// SELECT * FROM ts_user AS A LEFT JOIN users AS B ON B.user_name = A.uname ORDER BY B.user_id desc
			fmt.Println("users19: ", users19)
			fmt.Println(db19.Error) // <nil>
		}(i)
	}
	wg.Wait()
	select1 := make([]map[string]interface{}, 0)
	select1DB := connect.Query().Raw("SELECT * FROM users LIMIT 1").
		Get(&select1)
	fmt.Println("select1:", select1DB, select1)
	select2 := make([]map[string]interface{}, 0)
	select2DB := connect.Query().Select().From("users").
		Get(&select2)
	fmt.Println("select2:", select2DB, select2)
}
