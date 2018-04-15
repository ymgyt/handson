package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/k0kubun/pp"
)

type Product struct {
	gorm.Model
	Code  string
	Price int
}

func main() {

	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:3306)/gorm")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&Product{})

	var p Product
	_ = p
	q := db.Table("product").Select("count(*)").Where("price > ?", 1).QueryExpr()
	qq := db.Where("code = ?", 10).Joins("LEFT JOIN (?)", q).Table("hoge").QueryExpr()

	pp.Println(qq)
}
