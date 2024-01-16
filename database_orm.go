package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetConnectionORM() *gorm.DB {
	db := GetConnection()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "main.", // schema name
			SingularTable: false,
		},
	})

	if err != nil {
		panic(err)
	}

	return gormDb
}

func doMigration(db *gorm.DB) {
	db.AutoMigrate(Album{})
}
