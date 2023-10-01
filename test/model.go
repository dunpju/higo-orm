package main

import (
	"fmt"
	"github.com/dunpju/higo-orm/test/model/School"
)

type YY struct {
}

func (this *YY) String() string {
	return "yy"
}

func main() {
	fmt.Println(School.TableName())
	fmt.Println(School.TableName().Alias("a"))
	fmt.Println(School.SchoolName.AS("j"))
	fmt.Println(&YY{})
	fmt.Println(School.Select())
	fmt.Println(School.Raw())
}
