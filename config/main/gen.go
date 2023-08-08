package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gen/examples/dal/model"
	"gorm.io/gorm"
)

func main() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../dal", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Initialize a *gorm.DB instance
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(model.Customer{}, model.CreditCard{}, model.Bank{}, model.Passport{})

	// Generate default DAO interface for those generated structs from database
	companyGenerator := g.GenerateModelAs("company", "MyCompany"),
		g.ApplyBasic(
			g.GenerateModel("users"),
			companyGenerator,
			g.GenerateModelAs("people", "Person",
				gen.FieldIgnore("deleted_at"),
				gen.FieldNewTag("age", `json:"-"`),
			),
		)

	// Execute the generator
	g.Execute()
}
