package School

import (
	"github.com/dunpju/higo-orm/arm"
	"github.com/dunpju/higo-orm/test/model"
	"time"
)

const (
	SchoolId   arm.Fields = "schoolId"   //主键
	SchoolName arm.Fields = "schoolName" //学校名称
	Ip         arm.Fields = "ip"         //海康存储ip地址
	Port       arm.Fields = "port"       //海康存储端口
	UserName   arm.Fields = "userName"   //海康存储用户名
	Password   arm.Fields = "password"   //海康存储用户密码
	IsDelete   arm.Fields = "isDelete"   //是否删除:1-否,2-是
	CreateTime arm.Fields = "createTime" //创建时间
	UpdateTime arm.Fields = "updateTime" //更新时间
	DeleteTime arm.Fields = "deleteTime" //删除时间
)

type Model struct {
	*model.BaseModel
	SchoolId   int64     `gorm:"column:schoolId" json:"schoolId" comment:"主键"`
	SchoolName string    `gorm:"column:schoolName" json:"schoolName" comment:"学校名称"`
	Ip         string    `gorm:"column:ip" json:"ip" comment:"海康存储ip地址"`
	Port       string    `gorm:"column:port" json:"port" comment:"海康存储端口"`
	UserName   string    `gorm:"column:userName" json:"userName" comment:"海康存储用户名"`
	Password   string    `gorm:"column:password" json:"password" comment:"海康存储用户密码"`
	IsDelete   int       `gorm:"column:isDelete" json:"isDelete" comment:"是否删除:1-否,2-是"`
	CreateTime time.Time `gorm:"column:createTime" json:"createTime" comment:"创建时间"`
	UpdateTime time.Time `gorm:"column:updateTime" json:"updateTime" comment:"更新时间"`
	DeleteTime time.Time `gorm:"column:deleteTime" json:"deleteTime" comment:"删除时间"`
}

func New(properties ...arm.IProperty) *Model {
	return (&Model{BaseModel: model.NewBaseModel()}).New(properties...)
}

func TableName() *arm.TableName {
	return arm.NewTableName("school")
}

func (this *Model) New(properties ...arm.IProperty) *Model {
	arm.Properties(properties).Apply(this)
	arm.ApplyModel(this)
	return this
}

func (this *Model) TableName() *arm.TableName {
	return TableName()
}

func (this *Model) Apply(model *arm.Model) {
	this.Model = model
}

func (this *Model) Exist() bool {
	return this.SchoolId > 0
}

func WithSchoolId(schoolId int64) arm.IProperty {
	return arm.SetProperty(func(model arm.IModel) {
		model.(*Model).SchoolId = schoolId
	})
}
