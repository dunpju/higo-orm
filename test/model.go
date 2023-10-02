package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-orm/test/model/School"
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
	fmt.Println(School.SchoolName.AS("j"))
	fmt.Println(&YY{})
	res := make(map[string]interface{})
	School.Select().Where("schoolId", "=", 1).First(&res)
	fmt.Println(res)
	res = make(map[string]interface{})
	School.Raw("select * from ts_user").Get(&res)
	fmt.Println(res)
	School
	fmt.Println(res)
}
