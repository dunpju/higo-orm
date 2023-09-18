package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
)

func main() {
	him.DbConfig().
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
	db, err := him.Init()
	if err != nil {
		panic(err)
	}

	userNames := make([]string, 0)
	userNames = append(userNames, "ggg")
	userNames = append(userNames, "ttttt")
	sql, args, err := him.Query().Select("*").
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
	sql, args, err = him.Query().Select("*").
		From("users").
		WhereBetween("day", "2023-06-11", "2023-06-12").
		ToSql()
	fmt.Println("users2:", sql, args, err)
	db.Raw(sql, args...).Scan(&users2)
	fmt.Println(users2)

	users3 := make([]map[string]interface{}, 0)
	sql, args, err = him.Query().Select("*").
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
	sql, args, err = him.Query().Select("*").
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
	sql, args, err = him.Query().Select("user_id", "user_name", "day").
		From("users").
		Where("user_id", "=", 4).
		ToSql()
	// SELECT * FROM users WHERE (user_id = ?) [4] <nil>
	fmt.Println("users5:", sql, args, err)
	// SELECT * FROM users WHERE (user_id = 4)
	db.Raw(sql, args...).Scan(&users5)
	fmt.Println(users5)

	users6 := make([]map[string]interface{}, 0)
	him.Query().Select("user_id", "user_name").
		From("users").
		Where("user_id", "=", 8).
		First(&users6)
	// SELECT user_id, user_name FROM users WHERE (user_id = 8) LIMIT 1
	fmt.Println(users6)

	users7 := make([]map[string]interface{}, 0)
	db7 := him.Query().Select("user_id", "user_name").
		From("users1").
		Where("user_id", "=", 8).
		First(&users7)
	// SELECT user_id, user_name FROM users1 WHERE (user_id = 8) LIMIT 1
	fmt.Println(users7)
	// Error 1146 (42S02): Table 'test.users1' doesn't exist
	fmt.Println(db7.Error)

	users8 := make([]map[string]interface{}, 0)
	db8 := him.Query().Select("user_id", "user_name").
		From("users").
		Where("user_id", "=", 7).
		First(&users8)
	// SELECT user_id, user_name FROM users WHERE (user_id = 7) LIMIT 1
	fmt.Println(users8)
	fmt.Println(db8.Error) // <nil>

	users9 := make([]map[string]interface{}, 0)
	db9, paginate := him.Query().Select("user_id", "user_name").
		From("users").
		Where("user_name", "=", "kkk").
		Paginate(2, 2, &users9)
	// SELECT user_id, user_name FROM users LIMIT 2 OFFSET 0    {8 2 1 0 0xc0002aabe8}
	// SELECT user_id, user_name FROM users WHERE (user_name = 'kkk') LIMIT 2 OFFSET 0    {4 2 1 0 0xc0002aabe8}
	// SELECT user_id, user_name FROM users WHERE (user_name = 'kkk') LIMIT 2 OFFSET 2    {4 2 1 0 0xc0002aabe8}
	fmt.Println(users9, paginate)
	fmt.Println(db9.Error) // <nil>

	db10, count := him.Query().
		Select("count(distinct(user_name))").
		From("users").
		// Where("user_name", "=", "kkk").
		GroupBy("user_name").
		Count()
	// SELECT count(*) FROM `users` WHERE (user_name = 'kkk')
	// SELECT count(distinct(user_name)) FROM `users`
	// SELECT count(distinct(user_name)) FROM `users` GROUP BY `user_name`
	fmt.Println("db10: ", count)
	fmt.Println(db10.Error) // <nil>

	users11 := make([]map[string]interface{}, 0)
	db11 := him.Query().
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

	db12, sum := him.Query().
		From("users").
		Where("user_name", "=", "jjj").
		Sum("is_delete")
	// SELECT SUM(is_delete) count_ FROM users LIMIT 1
	// SELECT SUM(is_delete) count_ FROM users WHERE (user_name = 'jjj') LIMIT 1
	fmt.Println("db12: ", sum)
	fmt.Println(db12.Error) // <nil>

	users13 := make([]map[string]interface{}, 0)
	db13 := him.Query().Raw("SELECT * FROM users").
		Get(&users13)
	// SELECT * FROM users
	// SELECT SUM(is_delete) count_ FROM users WHERE (user_name = 'jjj') LIMIT 1
	fmt.Println("users13: ", users13)
	fmt.Println(db13.Error) // <nil>

	users14 := make([]map[string]interface{}, 0)
	db14 := him.Query().Select("*").
		From("users").
		WhereRaw(func(builder him.WhereRawBuilder) him.WhereRawBuilder {
			return builder.Raw("user_id = ?", 1)
		}).
		Get(&users14)
	// SELECT * FROM users WHERE (user_id = 1)
	fmt.Println("users14: ", users14)
	fmt.Println(db14.Error) // <nil>

	users15 := make([]map[string]interface{}, 0)
	db15 := him.Query().Select("*").
		From("users").
		OrderBy("user_id desc").
		First(&users15)
	// SELECT * FROM users ORDER BY user_id desc LIMIT 1
	fmt.Println("users15: ", users15)
	fmt.Println(db15.Error) // <nil>

	users16 := make([]map[string]interface{}, 0)
	db16 := him.Query().
		Distinct().
		Select("user_name").
		From("users").
		OrderBy("user_id desc").
		Get(&users16)
	// SELECT DISTINCT user_name FROM users ORDER BY user_id desc
	fmt.Println("users16: ", users16)
	fmt.Println(db16.Error) // <nil>

	users17 := make([]map[string]interface{}, 0)
	db17 := him.Query().
		Select("*").
		From("users AS A").
		Join("ts_user AS B", "B.uname", "=", "A.user_name").
		OrderBy("A.user_id desc").
		Get(&users17)
	// SELECT * FROM users AS A JOIN ts_user AS B ON B.uname = A.user_name ORDER BY A.user_id desc
	fmt.Println("users17: ", users17)
	fmt.Println(db17.Error) // <nil>

	users18 := make([]map[string]interface{}, 0)
	db18 := him.Query().
		Select("*").
		From("users AS A").
		InnerJoin("ts_user AS B", "B.uname", "=", "A.user_name").
		OrderBy("A.user_id desc").
		Get(&users18)
	// SELECT * FROM users AS A INNER JOIN ts_user AS B ON B.uname = A.user_name ORDER BY A.user_id desc
	fmt.Println("users18: ", users18)
	fmt.Println(db18.Error) // <nil>

	users19 := make([]map[string]interface{}, 0)
	db19 := him.Query().
		Select("*").
		From("ts_user AS A").
		LeftJoin("users AS B", "B.user_name", "=", "A.uname").
		OrderBy("B.user_id desc").
		Get(&users19)
	// SELECT * FROM ts_user AS A LEFT JOIN users AS B ON B.user_name = A.uname ORDER BY B.user_id desc
	fmt.Println("users19: ", users19)
	fmt.Println(db19.Error) // <nil>
}
