package main

import (
	"context"
	"errors"
	"fmt"
	"gorm-script/dal/model"
	"gorm-script/dal/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var ctx = context.Background()

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

	q := query.Use(db)
	rand.Seed(time.Now().Unix())
	err = Upsert(ctx, q, model.People{
		UUID: "Alice",
		Name: "Alice",
		Age:  17,
	})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 98; i++ {
		err = Upsert(ctx, q, model.People{
			UUID: RandStringBytes(64),
			Name: RandStringBytes(64),
			Age:  rand.Int63n(88),
		})
		if err != nil {
			panic(err)
		}
	}
	err = Upsert(ctx, q, model.People{
		UUID: "Alice",
		Name: "Jack",
		Age:  18,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(q.WithContext(ctx).People.GetMaxVersionCount())
}

func Upsert(ctx context.Context, q *query.Query, people model.People) error {
	p := q.People
	_people, err := p.WithContext(ctx).Where(p.UUID.Eq(people.UUID)).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = p.WithContext(ctx).Create(&people)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	people.Version = _people.Version + 1
	_, err = p.WithContext(ctx).Where(p.UUID.Eq(people.UUID)).Updates(&people)
	if err != nil {
		return err
	}
	return nil
}
