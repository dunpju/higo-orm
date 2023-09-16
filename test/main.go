package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dunpju/higo-orm/orm"
	"time"
)

func main() {

	sql, args, err := orm.Query().Select("*").
		From("tl_privilege_admin").
		Where("admin_name", "=", "Q88888888").
		ToSql()
	fmt.Println(sql, args, err) // SELECT * FROM tl_privilege_admin WHERE admin_name = ? [Q88888888] <nil>

	sql, args, err = orm.Query().Select("*").
		From("tl_privilege_admin").
		Where("admin_name", "=", "Q88888888").
		WhereNull("delete_time").
		ToSql()
	fmt.Println(sql, args, err) // SELECT * FROM tl_privilege_admin WHERE admin_name = ? AND delete_time IS NULL [Q88888888] <nil>

	privilegeActionSns := make([]string, 0)
	privilegeActionSns = append(privilegeActionSns, "pas28289907155466605743739230141")
	privilegeActionSns = append(privilegeActionSns, "pas35231745699179860431602784285")
	sql, args, err = orm.Query().Select("*").
		From("tl_privilege_action").
		WhereIn("privilege_action_sn", privilegeActionSns).
		WhereNull("delete_time").
		ToSql()
	fmt.Println(sql, args, err) // SELECT * FROM tl_privilege_action WHERE privilege_action_sn IN (?,?) AND delete_time IS NULL [pas28289907155466605743739230141 pas35231745699179860431602784285] <nil>

	sql, args, err = orm.Query().Select("*").
		From("tl_privilege_action").
		WhereIn("privilege_action_sn", privilegeActionSns).
		OrWhere("privilege_project_id", "=", 1).
		WhereNull("delete_time").
		ToSql()
	// SELECT * FROM tl_privilege_action WHERE privilege_action_sn IN (?,?) OR privilege_project_id = ? AND delete_time IS NULL [pas28289907155466605743739230141 pas35231745699179860431602784285 1] <nil>
	fmt.Println(sql, args, err)

	// SELECT * FROM tt WHERE a = ? OR b = ? [a b ] <nil>
	fmt.Println(squirrel.Select("*").From("tt").Where("a = ? OR b = ?", "a", "b ").ToSql())

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
	db, err := orm.Init()
	if err != nil {
		panic(err)
	}
	userNames := make([]string, 0)
	userNames = append(userNames, "ggg")
	userNames = append(userNames, "ttttt")
	sql, args, err = orm.Query().Select("*").
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
	sql, args, err = orm.Query().Select("*").
		From("users").
		WhereBetween("day", "2023-06-11", "2023-06-12").
		ToSql()
	fmt.Println("users2:", sql, args, err)
	db.Raw(sql, args...).Scan(&users2)
	fmt.Println(users2)

	users3 := make([]map[string]interface{}, 0)
	sql, args, err = orm.Query().Select("*").
		From("users").
		WhereRaw(func(builder orm.WhereRawBuilder) orm.WhereRawBuilder {
			return builder.Where("user_id", "=", 3).OrWhere("user_id", "=", 5)
		}).
		ToSql()
	// SELECT * FROM users WHERE ((user_id = ?) OR (user_id = ?)) [3 5] <nil>
	fmt.Println("users3:", sql, args, err)
	// SELECT * FROM users WHERE ((user_id = 3) OR (user_id = 5))
	db.Raw(sql, args...).Scan(&users3)
	fmt.Println(users3)

	users4 := make([]map[string]interface{}, 0)
	sql, args, err = orm.Query().Select("*").
		From("users").
		Where("user_id", "=", 4).
		OrWhereRaw(func(builder orm.WhereRawBuilder) orm.WhereRawBuilder {
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
	sql, args, err = orm.Query().Select("user_id", "user_name", "day").
		From("users").
		Where("user_id", "=", 4).
		ToSql()
	// SELECT * FROM users WHERE (user_id = ?) [4] <nil>
	fmt.Println("users5:", sql, args, err)
	// SELECT * FROM users WHERE (user_id = 4)
	db.Raw(sql, args...).Scan(&users5)
	fmt.Println(users5)

	users6 := make([]map[string]interface{}, 0)
	orm.Query().Select("user_id", "user_name").
		From("users").
		Where("user_id", "=", 8).
		First(&users6)
	// SELECT user_id, user_name FROM users WHERE (user_id = 8) LIMIT 1
	fmt.Println(users6)

	users7 := make([]map[string]interface{}, 0)
	db7 := orm.Query().Select("user_id", "user_name").
		From("users1").
		Where("user_id", "=", 8).
		First(&users7)
	// SELECT user_id, user_name FROM users1 WHERE (user_id = 8) LIMIT 1
	fmt.Println(users7)
	// Error 1146 (42S02): Table 'test.users1' doesn't exist
	fmt.Println(db7.Error)

	users8 := make([]map[string]interface{}, 0)
	db8 := orm.Query().Select("user_id", "user_name").
		From("users").
		Where("user_id", "=", 7).
		First(&users8)
	// SELECT user_id, user_name FROM users WHERE (user_id = 7) LIMIT 1
	fmt.Println(users8)
	fmt.Println(db8.Error) // <nil>

	users9 := make([]map[string]interface{}, 0)
	db9, paginate := orm.Query().Select("user_id", "user_name").
		From("users").
		Where("user_name", "=", "kkk").
		Paginate(2, 2, &users9)
	// SELECT user_id, user_name FROM users LIMIT 2 OFFSET 0    {8 2 1 0 0xc0002aabe8}
	// SELECT user_id, user_name FROM users WHERE (user_name = 'kkk') LIMIT 2 OFFSET 0    {4 2 1 0 0xc0002aabe8}
	// SELECT user_id, user_name FROM users WHERE (user_name = 'kkk') LIMIT 2 OFFSET 2    {4 2 1 0 0xc0002aabe8}
	fmt.Println(users9, paginate)
	fmt.Println(db9.Error) // <nil>

	db10, count := orm.Query().
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
	db11 := orm.Query().
		//Select("count(distinct(user_name)) count", "user_name").
		Select("count(user_name) count", "user_name").
		From("users").
		// Where("user_name", "=", "kkk").
		GroupBy("user_name").
		Get(&users11)
	// SELECT count(distinct(user_name)) count, user_name FROM users GROUP BY user_name
	// SELECT count(user_name) count, user_name FROM users GROUP BY user_name
	fmt.Println("db11: ", users11)
	fmt.Println(db11.Error) // <nil>

	db12, sum := orm.Query().
		From("users").
		Where("user_name", "=", "jjj").
		Sum("is_delete")
	// SELECT SUM(is_delete) count_ FROM users LIMIT 1
	// SELECT SUM(is_delete) count_ FROM users WHERE (user_name = 'jjj') LIMIT 1
	fmt.Println("db12: ", sum)
	fmt.Println(db12.Error) // <nil>

	users13 := make([]map[string]interface{}, 0)
	db13 := orm.Query().Raw("SELECT * FROM users").
		Where("user_id", "=", 1).
		Get(&users13)
	// SELECT * FROM users
	// SELECT SUM(is_delete) count_ FROM users WHERE (user_name = 'jjj') LIMIT 1
	fmt.Println("users13: ", users13)
	fmt.Println(db13.Error) // <nil>

	/*
		for i := 0; i < 100; i++ {
			go func() {
				admin := &Admin{}
				db.Raw("SELECT * FROM tl_privilege_admin WHERE admin_name = 'Q88888888' AND isnull(`delete_time`)").Scan(admin)
				fmt.Println(admin)
				privilegeFlag := &PrivilegeFlag{}
				db.Raw("SELECT * FROM `tl_privilege_flag`  WHERE (`privilege_project_id` = 1) AND (isnull(`delete_time`)) ORDER BY `interior_sort` asc").Scan(privilegeFlag)
				fmt.Println(privilegeFlag)
			}()
		}*/
}

type Admin struct {
	PrivilegeAdminId int64     `gorm:"column:privilege_admin_id" json:"privilege_admin_id" comment:"主键"`
	AdminName        string    `gorm:"column:admin_name" json:"admin_name" comment:"用户名"`
	CreateTime       time.Time `gorm:"column:create_time" json:"create_time" comment:"创建时间"`
	UpdateTime       time.Time `gorm:"column:update_time" json:"update_time" comment:"更新时间"`
	DeleteTime       time.Time `gorm:"column:delete_time" json:"delete_time" comment:"删除时间"`
}

type PrivilegeFlag struct {
	PrivilegeFlagId    int64     `gorm:"column:privilege_flag_id" json:"privilege_flag_id" comment:"主键"`
	PrivilegeFlagSn    string    `gorm:"column:privilege_flag_sn" json:"privilege_flag_sn" comment:"权限标签sn"`
	PrivilegeProjectId int64     `gorm:"column:privilege_project_id" json:"privilege_project_id" comment:"项目id"`
	Name               string    `gorm:"column:name" json:"name" comment:"标记名称"`
	Type               int       `gorm:"column:type" json:"type" comment:"类型:1-功能,2-菜单,3-数据"`
	ParentId           int64     `gorm:"column:parent_id" json:"parent_id" comment:"父级id:0-表示顶级"`
	Sort               int       `gorm:"column:sort" json:"sort" comment:"排序"`
	FrontRouteSn       string    `gorm:"column:front_route_sn" json:"front_route_sn" comment:"前端路由sn"`
	FrontButtonTag     string    `gorm:"column:front_button_tag" json:"front_button_tag" comment:"前端按钮标签"`
	State              int       `gorm:"column:state" json:"state" comment:"状态:1-启用,2-禁用"`
	PrivilegeLevelSn   string    `gorm:"column:privilege_level_sn" json:"privilege_level_sn" comment:"级别sn:对应tl_privilege_level表,权限类型为数据权限时必须填充"`
	InteriorSort       string    `gorm:"column:interior_sort" json:"interior_sort" comment:"组排序值"`
	CreateTime         time.Time `gorm:"column:create_time" json:"create_time" comment:"创建时间"`
	UpdateTime         time.Time `gorm:"column:update_time" json:"update_time" comment:"更新时间"`
	DeleteTime         time.Time `gorm:"column:delete_time" json:"delete_time" comment:"删除时间"`
}
