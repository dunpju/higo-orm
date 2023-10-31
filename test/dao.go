package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-orm/test/dao"
	"github.com/dunpju/higo-orm/test/entity/SchoolEntity"
	"gorm.io/gorm"
	"time"
)

func main() {
	_, err := him.DbConfig(him.DefaultConnect).
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
		SetSlowThreshold(1).
		SetColorful(true).
		Init()
	if err != nil {
		panic(err)
	}

	schoolDao := dao.NewSchoolDao()

	model := schoolDao.GetBySchoolId(1)
	fmt.Println(model)
	models := schoolDao.GetBySchoolIds([]int64{2, 3})
	fmt.Println(models)
	schoolEntity := SchoolEntity.New()
	schoolEntity.SchoolName = "SchoolName" + time.Now().Format(time.DateTime)
	schoolEntity.Ip = "Ip" + time.Now().Format(time.DateTime)
	schoolEntity.Port = "Port" + time.Now().Format(time.DateTime)
	schoolEntity.UserName = "UserName" + time.Now().Format(time.DateTime)
	schoolEntity.Password = "Password" + time.Now().Format(time.DateTime)
	schoolDao.SetData(schoolEntity).Add()

	models = schoolDao.GetBySchoolIds([]int64{5, 6})
	fmt.Println(models)

	schoolEntity = SchoolEntity.New()
	schoolEntity.SchoolName = "SchoolName" + time.Now().Format(time.DateTime)
	schoolEntity.Ip = "Ip" + time.Now().Format(time.DateTime)
	schoolEntity.Port = "Port" + time.Now().Format(time.DateTime)
	schoolEntity.UserName = "UserName" + time.Now().Format(time.DateTime)
	schoolEntity.Password = "Password" + time.Now().Format(time.DateTime)
	_, schoolEntity.SchoolId = schoolDao.SetData(schoolEntity).Add()

	SchoolEntity.FlagUpdate.Apply(schoolEntity)
	schoolEntity.SchoolName = "SchoolName" + time.Now().Format(time.DateTime)
	schoolEntity.Ip = "Ip" + time.Now().Format(time.DateTime)
	schoolDao.SetData(schoolEntity).Update()

	err = schoolDao.Begin().Transaction(func(tx *gorm.DB) error {
		schoolEntity = SchoolEntity.New()
		schoolEntity.SchoolName = "SchoolName" + time.Now().Format(time.DateTime)
		schoolEntity.Ip = "Ip" + time.Now().Format(time.DateTime)
		schoolEntity.Port = "Port" + time.Now().Format(time.DateTime)
		schoolEntity.UserName = "UserName" + time.Now().Format(time.DateTime)
		schoolEntity.Password = "Password" + time.Now().Format(time.DateTime)
		_, schoolEntity.SchoolId = schoolDao.TX(tx).SetData(schoolEntity).Add()

		SchoolEntity.FlagUpdate.Apply(schoolEntity)
		schoolEntity.SchoolName = "SchoolName" + time.Now().Format(time.DateTime)
		schoolEntity.Ip = "Ip" + time.Now().Format(time.DateTime)
		schoolDao.TX(tx).SetData(schoolEntity).Update()

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
