package dao

import (
	"fmt"
	"github.com/dunpju/higo-orm/arm"
	"github.com/dunpju/higo-orm/event"
	"github.com/dunpju/higo-orm/exception/DaoException"
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-orm/test/entity/SchoolEntity"
	"github.com/dunpju/higo-orm/test/model/School"
	"gorm.io/gorm"
)

type SchoolDao struct {
	*arm.BaseDao
	model *School.Model
}

func NewSchoolDao() *SchoolDao {
	schoolDao := &SchoolDao{model: School.New()}
	schoolDao.BaseDao = arm.NewBaseDao(schoolDao)
	return schoolDao
}

func (this *SchoolDao) SetModel(model arm.IModel) {
	this.model = model.(*School.Model)
}

func (this *SchoolDao) GetModel() arm.IModel {
	return this.model
}

func (this *SchoolDao) Model() *School.Model {
	return School.New()
}

func (this *SchoolDao) Models() []*School.Model {
	return make([]*School.Model, 0)
}

func (this *SchoolDao) TX(tx *gorm.DB) *SchoolDao {
	this.model.TX(tx)
	return this
}

func (this *SchoolDao) SetData(entity *SchoolEntity.Entity) arm.IDao {
	return this.model.Builder(this, func() {
		if !entity.PrimaryEmpty() || entity.IsEdit() { //编辑
			if !this.GetBySchoolId(entity.SchoolId).Exist() {
				DaoException.Throw("不存在", 0)
			}
			this.model.Where(School.SchoolId, "=", entity.SchoolId)
			if entity.Equals(SchoolEntity.FlagDelete) {
				// todo::填充修改字段
				this.model.Set(School.DeleteTime, entity.UpdateTime)
			} else if entity.Equals(SchoolEntity.FlagUpdate) {
				// todo::填充修改字段
				this.model.Set(School.SchoolName, "update_"+entity.SchoolName) //学校名称
				this.model.Set(School.Ip, "update_"+entity.Ip)                 //海康存储ip地址
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
	})
}

func Init() {
	event.AddEvent(event.BeforeInsert, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println(111)
		fmt.Println(data.Sql)
	})
}

// Add 添加
func (this *SchoolDao) Add() (gormDB *gorm.DB, lastInsertId int64) {
	gormDB, lastInsertId = this.model.Insert().LastInsertId()
	this.CheckError(gormDB)
	return
}

// Update 更新
func (this *SchoolDao) Update() (*gorm.DB, int64) {
	gormDB, affected := this.model.Update().Exec()
	this.CheckError(gormDB)
	return gormDB, affected
}

// GetBySchoolId id查询
func (this *SchoolDao) GetBySchoolId(schoolId int64) *School.Model {
	model := this.Model()
	gormDB := this.model.Select().Where(School.SchoolId, "=", schoolId).First(&model)
	this.CheckError(gormDB)
	return model
}

// GetBySchoolIds id集查询
func (this *SchoolDao) GetBySchoolIds(schoolIds []int64) []*School.Model {
	models := this.Models()
	gormDB := this.model.Select(School.SchoolId, School.SchoolName).WhereIn(School.SchoolId, schoolIds).Get(&models)
	this.CheckError(gormDB)
	return models
}

// DeleteBySchoolId 硬删除
func (this *SchoolDao) DeleteBySchoolId(schoolId int64) *gorm.DB {
	gormDB, _ := this.model.Delete().Where(School.SchoolId, "=", schoolId).Exec()
	this.CheckError(gormDB)
	return gormDB
}

// Paginate 列表
func (this *SchoolDao) Paginate(perPage, page uint64) him.IPaginate {
	models := this.Models()
	gormDB, paginate := this.model.Select().Paginate(page, perPage, &models)
	this.CheckError(gormDB)
	return paginate
}
