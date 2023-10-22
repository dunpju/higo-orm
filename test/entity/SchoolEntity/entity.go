package SchoolEntity

import (
	"github.com/dunpju/higo-orm/arm"
	"time"
)

const (
	FlagDelete arm.Flag = iota + 1
	FlagUpdate
)

type Entity struct {
	isEdit     bool
	flag       arm.Flag
	SchoolId   int64     `json:"schoolId"   comment:"主键"`
	SchoolName string    `json:"schoolName" comment:"学校名称"`
	Ip         string    `json:"ip"         comment:"海康存储ip地址"`
	Port       string    `json:"port"       comment:"海康存储端口"`
	UserName   string    `json:"userName"   comment:"海康存储用户名"`
	Password   string    `json:"password"   comment:"海康存储用户密码"`
	IsDelete   int       `json:"isDelete"   comment:"是否删除:1-否,2-是"`
	CreateTime time.Time `json:"createTime" comment:"创建时间"`
	UpdateTime time.Time `json:"updateTime" comment:"更新时间"`
	DeleteTime time.Time `json:"deleteTime" comment:"删除时间"`
}

func New() *Entity {
	tn := time.Now()
	return &Entity{CreateTime: tn, UpdateTime: tn}

}

func (this *Entity) IsEdit() bool {
	return this.isEdit
}

func (this *Entity) Edit(isEdit bool) {
	this.isEdit = isEdit
}

func (this *Entity) SetFlag(flag arm.Flag) {
	this.flag = flag
	this.isEdit = true
}

func (this *Entity) Flag() arm.Flag {
	return this.flag
}

func (this *Entity) PrimaryEmpty() bool {
	return this.SchoolId == 0
}