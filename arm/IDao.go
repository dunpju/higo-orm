package arm

import "gorm.io/gorm"

type IDao interface {
	Add() (gormDB *gorm.DB, lastInsertId int64)
	Update() *gorm.DB
}
