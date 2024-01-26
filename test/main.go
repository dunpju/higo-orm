package main

import (
	"fmt"
	"time"
)

func main() {
	var value interface{}
	value = "111111111111122222222222222222222222"
	switch value.(type) {
	case string:
		fmt.Println(fmt.Sprintf("%v", value))
	case int, int64:
		fmt.Println(fmt.Sprintf("%v", value))
	case float32, float64:
		fmt.Println(fmt.Sprintf("%v", value))
	}

	/*
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
		fmt.Println(squirrel.Select("*").From("tt").Where("a = ? OR b = ?", "a", "b ").ToSql())*/

	// select_run(db)

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
