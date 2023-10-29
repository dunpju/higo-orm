package dao

import (
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-orm/test/entity/SchoolEntity"
	"github.com/dunpju/higo-orm/test/model/School"
	"gorm.io/gorm"
)

type SchoolDao struct {
	model *School.Model
}

func NewSchoolDao() *SchoolDao {
	return &SchoolDao{model: School.New()}
}

func (this *SchoolDao) Model() *School.Model {
	return School.New()
}

func (this *SchoolDao) Models() []*School.Model {
	return make([]*School.Model, 0)
}

func (this *SchoolDao) TX(tx *gorm.DB) *SchoolDao {
	this.model.DB().TX(tx)
	return this
}

func (this *SchoolDao) CheckError(gormDB *gorm.DB) {
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
}

func (this *SchoolDao) SetData(entity *SchoolEntity.Entity) *SchoolDao {
	if !entity.PrimaryEmpty() || entity.IsEdit() { //编辑
		if !this.GetBySchoolId(entity.SchoolId).Exist() {
			//DaoException.Throw("不存在", 0)
		}
		this.model.Update().Where(School.SchoolId, "=", entity.SchoolId)
		if SchoolEntity.FlagDelete == entity.Flag() {

		} else if SchoolEntity.FlagUpdate == entity.Flag() {

		}
		this.model.Set(School.UpdateTime, entity.UpdateTime)
	} else { //新增
		this.model.Set(School.SchoolName, entity.SchoolName) //学校名称
		this.model.Set(School.Ip, entity.Ip)                 //海康存储ip地址
		this.model.Set(School.Port, entity.Port)             //海康存储端口
		this.model.Set(School.UserName, entity.UserName)     //海康存储用户名
		this.model.Set(School.Password, entity.Password)     //海康存储用户密码
		this.model.Set(School.CreateTime, entity.CreateTime) //创建时间
		this.model.Set(School.UpdateTime, entity.UpdateTime) //更新时间
	}
	return this
}

// Add 添加
func (this *SchoolDao) Add() (gormDB *gorm.DB, lastInsertId int64) {
	gormDB, lastInsertId = this.model.Insert().LastInsertId()
	this.CheckError(gormDB)
	return
}

// Update 更新
func (this *SchoolDao) Update() *gorm.DB {
	gormDB, _ := this.model.DB().Exec()
	this.CheckError(gormDB)
	return gormDB
}

// GetBySchoolId id查询
func (this *SchoolDao) GetBySchoolId(schoolId int64) *School.Model {
	model := this.Model()
	gormDB := this.model.Select().
		Where(School.SchoolId, "=", schoolId).
		First(&model)
	this.CheckError(gormDB)
	return model
}

// GetBySchoolIds id集查询
func (this *SchoolDao) GetBySchoolIds(schoolIds []int64, fields ...string) []*School.Model {
	models := this.Models()
	gormDB := this.model.Select(fields...).
		WhereIn(School.SchoolId, schoolIds).
		Get(&models)
	this.CheckError(gormDB)
	return models
}

// DeleteBySchoolId 硬删除
func (this *SchoolDao) DeleteBySchoolId(schoolId int64) *gorm.DB {
	gormDB, _ := this.model.Delete().
		Where(School.SchoolId, "=", schoolId).
		Exec()
	this.CheckError(gormDB)
	return gormDB
}

// List 列表
func (this *SchoolDao) List(perPage, page uint64, where map[string]interface{}) him.Paginate {
	var models []struct {
		SchoolId   string `gorm:"column:schoolId" json:"schoolId"`
		SchoolName string `gorm:"column:schoolName" json:"schoolName"`
		Ip         string `gorm:"column:ip" json:"ip"`
		Port       string `gorm:"column:port" json:"port"`
		UserName   string `gorm:"column:userName" json:"userName"`
		Password   string `gorm:"column:password" json:"password"`
	}
	gormDB, paginate := this.model.Select().
		Paginate(page, perPage, &models)
	this.CheckError(gormDB)
	return paginate
}
