package main

import (
	"github.com/ag9920/db-demo/gendemo/dal/model"
	"gorm.io/gen"
)

func main() {

	g := gen.NewGenerator(gen.Config{
		OutPath: "../../dal/query",
		Mode:    gen.WithDefaultQuery,
	})

	g.ApplyBasic(model.Passport{}, model.User{})

	g.ApplyInterface(func(model.Method) {}, model.User{})
	g.ApplyInterface(func(model.UserMethod) {}, model.User{})

	g.Execute()
}
