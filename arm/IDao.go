package arm

import "gorm.io/gorm"

type IDao interface {
	SetModel(model IModel)
	GetModel() IModel
	Add() (gormDB *gorm.DB, lastInsertId int64)
	Update() *gorm.DB
	CheckError(gormDB *gorm.DB)
}
