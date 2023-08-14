// Package gen @Author: youngalone [2023/8/7]
package gen

import (
	"github.com/RaymondCode/simple-demo/database/mysql"
	"gorm.io/gen"
)

func Setup() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./dao",
		ModelPkgPath:      "./model",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable:     false,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	})
	g.UseDB(mysql.DB)
	allModel := g.GenerateAllTable()
	g.ApplyBasic(allModel...)
	g.Execute()
}
