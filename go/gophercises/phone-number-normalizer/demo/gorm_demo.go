package demo

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Name   string
	Gender string `gorm:"default:Male"`
	Number string
}

func GORMDemo() error {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "mydb"
	)

	// Connect data source
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Create corresponding tables
	// AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
	// It will change existing column’s type if its size, precision, nullable changed.
	// It WON’T delete unused columns to protect your data.
	err = db.AutoMigrate(&Person{})
	if err != nil {
		return err
	}

	// Insert records
	persons := []Person{
		{Name: "ZHANG San", Gender: "Male", Number: "1234567890"},
		{Name: "LI Si", Gender: "Male", Number: "123 456 7891"},
		{Name: "WANG Wu", Gender: "Female", Number: "(123) 456 7892"},
		{Name: "ZHAO Liu", Gender: "Male", Number: "(123) 456-7893"},
		{Name: "ZHENG Qi", Gender: "Female", Number: "123-456-7894"},
	}

	result := db.Create(&persons)
	if result.Error != nil {
		return result.Error
	}

	// Query or select records
	var personFirst, personTake, personLast, personID Person

	db.First(&personFirst)
	fmt.Println(personFirst)

	db.Take(&personTake)
	fmt.Println(personTake)

	db.Last(&personLast)
	fmt.Println(personLast)

	db.Last(&personID, 3)
	fmt.Println(personID)

	var personsFind []Person

	db.Where("gender = ?", "Male").Find(&personsFind)
	fmt.Println(personsFind)

	// Update record
	db.Model(&personID).Update("gender", "Male")
	
	db.Last(&personID, 3)
	fmt.Println(personID)

	// Delete records
	db.Where("1 = 1").Delete(&Person{})

	// Drop table
	db.Exec(`drop table people`)

	return nil
}
