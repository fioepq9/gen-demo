package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	UUID    string `gorm:"column:uuid"`
	Name    string `gorm:"column:name"`
	Age     int    `gorm:"column:age"`
	Version int    `gorm:"column:version"`
}

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"0.0.0.0",
		"3306",
		"douyin_12306")
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

}
