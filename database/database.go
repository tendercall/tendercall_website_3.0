package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"tendercall-website.com/main/service/helper"
)

func Initdb() {
	// PostgreSQL connection parameters
	const (
		host     = "ep-winter-mud-a2uj3wgl.eu-central-1.pg.koyeb.app"
		port     = 5432
		user     = "koyeb-adm"
		password = "dRODw38HAjxK"
		dbname   = "koyebdb"
	)

	// Construct the connection string
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)

	// Attempt to connect to the database
	var err error
	helper.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = helper.DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Database connection established")

	// createTable := `CREATE TABLE IF NOT EXISTS testimonial (
	// id SERIAL PRIMARY KEY,
	// image VARCHAR(256),
	// description VARCHAR(256),
	// name VARCHAR(256),
	// position VARCHAR(256),
	// created_date TIMESTAMP NOT NULL DEFAULT NOW(),
	// updated_date TIMESTAMP NOT NULL DEFAULT NOW()
	// )`

	// _, err = helper.DB.Exec(createTable)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Table created successfully")

	//Example ALTER TABLE statement to add a new column
	// 	query := `
	// 		ALTER TABLE rating
	// 		ALTER COLUMN rating INTEGER;
	// 	`

	// 	// Execute the ALTER TABLE statement
	// 	_, err = helper.DB.Exec(query)
	// 	if err != nil {
	// 		log.Fatalf("Error executing ALTER TABLE statement: %v\n", err)
	// 	}

	// 	fmt.Println("Column added successfully.")
	// }
}
