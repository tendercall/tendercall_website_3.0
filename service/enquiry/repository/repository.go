package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"tendercall-website.com/main/service/models"
)

var DB *sql.DB

func PostEnquiry(email, message, enquiry_type, enquiry_id string, created_date, updated_date time.Time) (uint, error) {
	var id uint

	currentTime := time.Now()

	err := DB.QueryRow("INSERT INTO enquiry (email, message, enquiry_type, enquiry_id, created_date, updated_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", email, message, enquiry_type, enquiry_id, currentTime, currentTime).Scan(&id)
	if err != nil {
		return 0, err
	}

	fmt.Println("Post Successful")
	return id, nil
}

func GetEnquiry() ([]models.Enquiry, error) {
	rows, err := DB.Query("SELECT id, email, message, enquiry_type, enquiry_id, created_date, updated_date FROM enquiry")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enquiries []models.Enquiry
	for rows.Next() {
		var enquiry models.Enquiry
		err := rows.Scan(&enquiry.ID, &enquiry.Email, &enquiry.Message, &enquiry.EnquiryType, &enquiry.EnquiryID, &enquiry.CreatedDate, &enquiry.UpdatedDate)
		if err != nil {
			return nil, err
		}
		enquiries = append(enquiries, enquiry)
	}

	fmt.Println("Get Successful")
	return enquiries, nil
}

func GetEnquiryById(id uint) (*models.Enquiry, error) {
	var enquiry models.Enquiry

	err := DB.QueryRow("SELECT id, email, message, enquiry_type, enquiry_id, created_date, updated_date FROM enquiry WHERE id=$1", id).Scan(&enquiry.ID, &enquiry.Email, &enquiry.Message, &enquiry.EnquiryType, &enquiry.EnquiryID, &enquiry.CreatedDate, &enquiry.UpdatedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no enquiry found with id %d", id)
		}
		log.Println("Error", err)
		return nil, err
	}

	return &enquiry, nil
}

func PutEnquiry(id uint, email, message, enquiry_type, enquiry_id string, updated_date time.Time) error {
	result, err := DB.Exec("UPDATE enquiry SET email=$1, message=$2, enquiry_type=$3, enquiry_id=$4, updated_date=$5 where id=$6", email, message, enquiry_type, enquiry_id, time.Now(), id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to determine affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("enquiry not found")
	}

	fmt.Println("Update successfull")

	return nil
}

func DeleteEnquiry(id uint) error {
	result, err := DB.Exec("DELETE FROM enquiry WHERE id=$1", id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to determine affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(" enquiry found")
	}

	fmt.Println("Delete successfull")

	return nil
}
