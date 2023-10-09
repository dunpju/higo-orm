package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-orm/test/model/School"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type YY struct {
}

func (this *YY) String() string {
	return "yy"
}

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
	fmt.Println(School.TableName())
	fmt.Println(School.TableName().Alias("a"))
	fmt.Println(School.SchoolName)
	fmt.Println(School.SchoolName.AS("j"))
	fmt.Println(School.SchoolName.Pre("A"))
	fmt.Println(School.SchoolName.Pre("A").DESC())
	fmt.Println(School.SchoolName.Pre("A").AS("sn"))
	fmt.Println(School.SchoolName.COUNT())
	fmt.Println(School.SchoolName.COUNT().Pre("B"))
	fmt.Println(School.SchoolName.COUNT().AS("G"))
	fmt.Println(School.SchoolName.SUM().String())
	fmt.Println(&YY{})
	res := make(map[string]interface{})
	School.New().Select().Where("schoolId", "=", 1).First(&res)
	fmt.Println(res)
	res = make(map[string]interface{})
	School.New().Raw("select * from ts_user").Get(&res)
	fmt.Println(res)
	school := School.New(School.WithSchoolId(130), School.WithSchoolName("小学"))
	fmt.Println(school)
	School.New().Alias("A").Select().Where(School.SchoolId, "=", 1).Get(&res)
	err = School.New().Begin().Transaction(func(tx *gorm.DB) error {
		_, rowsAffected := School.New().
			TX(tx).
			Update().
			Set(School.UserName, "33").
			Where(School.SchoolId, "=", 1).
			Exec()
		fmt.Println(rowsAffected)
		school := School.New().
			TX(tx).
			Insert().
			Columns(School.SchoolName, School.Ip, School.Port, School.UserName, School.Password, School.CreateTime, School.UpdateTime)
		school.Values(rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now())
		school.Values(rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now())
		_, affected := school.Save()
		fmt.Println(affected)
		_, affected1 := School.New().
			TX(tx).
			Insert().
			Columns(School.SchoolName, School.Ip, School.Port, School.UserName, School.Password, School.CreateTime, School.UpdateTime).
			Values(rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now()).
			Values(rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now()).
			Values(rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now()).
			Save()
		fmt.Println(affected1)
		_, lastInsertId := School.New().
			TX(tx).
			Insert().
			Column(School.SchoolName, rand.Intn(6)).
			Column(School.Ip, rand.Intn(6)).
			Column(School.Port, rand.Intn(6)).
			Column(School.UserName, rand.Intn(6)).
			Column(School.Password, rand.Intn(6)).
			Column(School.CreateTime, time.Now()).
			Column(School.UpdateTime, time.Now()).
			LastInsertId()
		fmt.Println(lastInsertId)
		School.New().TX(tx).Raw("UPDATE school SET userName = '33ff' WHERE (schoolId = ?)", lastInsertId).Exec()
		//School.Delete().TX(tx).Where(School.SchoolId, "=", lastInsertId).Exec()
		//return fmt.Errorf("测试事务")
		return nil
	})
	fmt.Println(err)
}
