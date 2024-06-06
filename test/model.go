package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/event"
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

func checkError(gormDB *gorm.DB) {
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
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

	// 测试事件
	event.AddEvent(event.BeforeInsert, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("BeforeInsert")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.BeforeUpdate, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("BeforeUpdate")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.BeforeDelete, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("BeforeDelete")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.BeforeRaw, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("BeforeRaw")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.AfterInsert, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("AfterInsert")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.AfterUpdate, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("AfterUpdate")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.AfterDelete, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("AfterDelete")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.AfterRaw, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("AfterRaw")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.BeforeSelect, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("BeforeSelect")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.BeforeCount, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("BeforeSelect")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.AfterSelect, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("AfterSelect")
		fmt.Println(data.Sql)
	})
	event.AddEvent(event.AfterCount, School.New().TableName().String(), func(data event.EventRecord) {
		fmt.Println("AfterCount")
		fmt.Println(data.Sql)
	})

	fmt.Println(School.TableName())
	fmt.Println(School.TableName().Alias("a"))
	fmt.Println(School.SchoolName)
	fmt.Println(School.SchoolName.AS("j"))
	fmt.Println(School.SchoolName.Pre("A"))
	fmt.Println(School.SchoolName.Pre("A").DESC())
	fmt.Println(School.SchoolName.Pre("A").AS("sn"))
	fmt.Println(School.SchoolName.COUNT())
	fmt.Println(School.SchoolName.Pre("B").COUNT())
	fmt.Println(School.SchoolName.COUNT().AS("G"))
	fmt.Println(School.SchoolName.SUM().String())
	fmt.Println(School.SchoolName.IN(1, "g"))
	fmt.Println(&YY{})

	school := School.New(School.WithSchoolId(130), School.WithSchoolName("小学"))
	fmt.Println(school)

	err = School.New().Begin().Transaction(func(tx *gorm.DB) error {
		fmt.Printf("tx %p\n", tx)
		res := make(map[string]interface{})
		tx1 := School.New().TX(tx).Select().Where("schoolId", "=", 1).First(&res)
		fmt.Printf("tx1 %p\n", tx1)
		fmt.Println(res)
		tx1_1 := School.New().TX(tx).Select().Raw("SELECT * FROM `school` WHERE (schoolId = ?) LIMIT 1", 2).First(&res)
		fmt.Printf("tx1_1 %p\n", tx1_1)
		fmt.Println(res)
		res = make(map[string]interface{})
		tx2 := School.New().TX(tx).Raw("select * from ts_user").Get(&res)
		fmt.Printf("tx2 %p\n", tx2)
		fmt.Println(res)
		var res1 []map[string]interface{}
		tx3, paginate := School.New().TX(tx).Select().Paginate(4, 2, &res1)
		fmt.Printf("tx3 %p\n", tx3)
		fmt.Println(paginate.GetItems())
		fmt.Println(paginate.GetTotal())
		fmt.Println(paginate.GetCurrentPage())
		fmt.Println(paginate.GetPerPage())

		tx4 := School.New().TX(tx).Alias("A").Select().Where(School.SchoolId, "=", 1).Get(&res)
		fmt.Printf("tx4 %p\n", tx4)
		fmt.Println(res)

		tx5, rowsAffected := School.New().
			TX(tx).
			Update().
			Set(School.UserName, "33").
			Where(School.SchoolId, "=", 1).
			Exec()
		fmt.Printf("tx5 %p\n", tx5)
		checkError(tx5)
		fmt.Println(rowsAffected)
		school1 := School.New().TX(tx).Update().Where(School.SchoolId, "=", 2)
		school1.Set(School.UserName, "22333")
		school1.Set(School.Ip, "22111")
		tx6, rowsAffected := school1.Exec()
		fmt.Printf("tx6 %p\n", tx6)
		checkError(tx6)
		fmt.Println(rowsAffected)
		school := School.New().
			TX(tx).
			Insert().
			Columns(School.SchoolName, School.Ip, School.Port, School.UserName, School.Password, School.CreateTime, School.UpdateTime)
		school.Values(rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now())
		school.Values(rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now())
		school.OnDuplicateKeyUpdate(School.UpdateTime.VALUES())
		tx7, affected := school.Save()
		fmt.Printf("tx7 %p\n", tx7)
		checkError(tx7)
		fmt.Println(affected)
		tx8, affected1 := School.New().
			TX(tx).
			Insert().
			Columns(School.SchoolName, School.Ip, School.Port, School.UserName, School.Password, School.CreateTime, School.UpdateTime).
			Values(1, rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now()).
			Values(1, rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now()).
			Values(1, rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6), time.Now(), time.Now()).
			OnDuplicateKeyUpdate(School.UpdateTime.VALUES()).
			Save()
		fmt.Printf("tx8 %p\n", tx8)
		checkError(tx8)
		fmt.Println(affected1)
		tx9, lastInsertId := School.New().
			TX(tx).
			Insert().
			Column(School.SchoolName, rand.Intn(6)).
			Column(School.Ip, rand.Intn(6)).
			Column(School.Port, rand.Intn(6)).
			Column(School.UserName, rand.Intn(6)).
			Column(School.Password, rand.Intn(6)).
			Column(School.CreateTime, time.Now()).
			Column(School.UpdateTime, time.Now()).
			OnDuplicateKeyUpdate(School.UpdateTime.VALUES()).
			LastInsertId()
		fmt.Printf("tx9 %p\n", tx9)
		checkError(tx9)
		fmt.Println(lastInsertId)
		tx10, _, _ := School.New().TX(tx).Raw("UPDATE school SET userName = '33ff' WHERE (schoolId = ?)", lastInsertId).Exec()
		fmt.Printf("tx10 %p\n", tx10)
		checkError(tx10)
		tx11, _ := School.New().TX(tx).Delete().Where(School.SchoolId, "=", lastInsertId+1).Exec()
		fmt.Printf("tx11 %p\n", tx11)
		checkError(tx11)
		tx12, _ := School.New().TX(tx).Update().
			CaseWhen(School.Ip.Case(School.SchoolId).When(21, 21).When(22, 22)).
			CaseWhen(School.Port.Case(School.SchoolId).When(21, 21).When(22, 22)).
			WhereIn(School.SchoolId, []int64{21, 22}).
			Exec()
		fmt.Printf("tx12 %p\n", tx12)
		checkError(tx12)
		tx13, _ := School.New().TX(tx).Update().
			Set(School.Ip).
			CaseWhen(School.Ip.Case().When(School.SchoolId.Eq(23), 23).When(School.SchoolId.Eq(24), 24).Else(`'w11'`)).
			CaseWhen(School.Port.Case().When(School.SchoolId.Eq(23), 23).When(School.SchoolId.Eq(24), 24).Else(School.Port)).
			WhereIn(School.SchoolId, []int64{23, 24, 25}).
			Exec()
		fmt.Printf("tx13 %p\n", tx13)
		checkError(tx13)
		//return fmt.Errorf("测试事务")
		return nil
	})
	fmt.Println(err)
}
