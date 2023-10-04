package model

import (
	"github.com/dunpju/higo-orm/arm"
	"github.com/dunpju/higo-orm/him"
)

type BaseModel struct {
	*arm.Model
}

func NewBaseModel() *BaseModel {
	return &BaseModel{}
}

func (this *BaseModel) Connection() string {
	return him.DefaultConnect
}
