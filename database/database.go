package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"tendercall-website.com/main/service/enquiry/repository"
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
	repository.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = repository.DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Database connection established")

	// createTable := `CREATE TABLE IF NOT EXISTS enquiry (
	// id SERIAL PRIMARY KEY,
	// email VARCHAR(256),
	// message VARCHAR(256),
	// enquiry_type VARCHAR(256),
	// enquiry_id VARCHAR(256),
	// created_date TIMESTAMP NOT NULL DEFAULT NOW(),
	// updated_date TIMESTAMP NOT NULL DEFAULT NOW()
	// )`

	// _, err = repository.DB.Exec(createTable)
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
	// 	_, err = repository.DB.Exec(query)
	// 	if err != nil {
	// 		log.Fatalf("Error executing ALTER TABLE statement: %v\n", err)
	// 	}

	// 	fmt.Println("Column added successfully.")
	// }
}
