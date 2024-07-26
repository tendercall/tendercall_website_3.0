package helper

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"tendercall-website.com/main/service/models"
)

var DB *sql.DB

// Enquiry POST, GET, PUT and DELETE
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

// Testimnaol POST, GET, PUT and DELETE
func PostTestimonial(image, description, name, position string, created_date, updated_date time.Time) (uint, error) {
	var id uint

	currentTime := time.Now()

	err := DB.QueryRow("INSERT INTO testimonial (image, description, name, position, created_date, updated_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", image, description, name, position, currentTime, currentTime).Scan(&id)
	if err != nil {
		return 0, err
	}

	fmt.Println("Post Successful")
	return id, nil
}

func GetTestimonial() ([]models.Testimonial, error) {
	rows, err := DB.Query("SELECT id, image, description, name, position, created_date, updated_date FROM testimonial")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testimonials []models.Testimonial
	for rows.Next() {
		var testimonial models.Testimonial
		err := rows.Scan(&testimonial.ID, &testimonial.Image, &testimonial.Description, &testimonial.Name, &testimonial.Position, &testimonial.CreatedDate, &testimonial.UpdatedDate)
		if err != nil {
			return nil, err
		}
		testimonials = append(testimonials, testimonial)
	}

	fmt.Println("Get Successful")
	return testimonials, nil
}

func GetTestimonialById(id uint) (*models.Testimonial, error) {
	var testimonial models.Testimonial

	err := DB.QueryRow("SELECT id, image, description, name, position, created_date, updated_date FROM testimonial WHERE id=$1", id).Scan(&testimonial.ID, &testimonial.Image, &testimonial.Description, &testimonial.Name, &testimonial.Position, &testimonial.CreatedDate, &testimonial.UpdatedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no testimonial found with id %d", id)
		}
		log.Println("Error", err)
		return nil, err
	}

	return &testimonial, nil
}

func PutTestimonial(id uint, image, description, name, position string, updated_date time.Time) error {
	result, err := DB.Exec("UPDATE testimonial SET image=$1, description=$2, name=$3, position=$4, updated_date=$5 where id=$6", image, description, name, position, time.Now(), id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to determine affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("testimonial not found")
	}

	fmt.Println("Update successfull")

	return nil
}

func DeleteTestimonial(id uint) error {
	result, err := DB.Exec("DELETE FROM testimonial WHERE id=$1", id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to determine affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(" testimonial found")
	}

	fmt.Println("Delete successfull")

	return nil
}

// Faq POST, GET, PUT and DELETE
func PostFaq(question, answer string, created_date, updated_date time.Time) (uint, error) {
	var id uint

	currentTime := time.Now()

	err := DB.QueryRow("INSERT INTO faq (question, answer, created_date, updated_date) VALUES ($1, $2, $3, $4) RETURNING id", question, answer, currentTime, currentTime).Scan(&id)
	if err != nil {
		return 0, err
	}

	fmt.Println("Post Successful")
	return id, nil
}

func GetFaq() ([]models.Faq, error) {
	rows, err := DB.Query("SELECT id, question, answer, created_date, updated_date FROM faq")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faqs []models.Faq
	for rows.Next() {
		var faq models.Faq
		err := rows.Scan(&faq.ID, &faq.Question, &faq.Answer, &faq.CreatedDate, &faq.UpdatedDate)
		if err != nil {
			return nil, err
		}

		faqs = append(faqs, faq)
	}

	return faqs, nil
}

func GetFaqById(id uint) (*models.Faq, error) {
	var faq models.Faq

	err := DB.QueryRow("SELECT id, question, answer, created_date, updated_date FROM faq WHERE id=$1", id).Scan(&faq.ID, &faq.Question, &faq.Answer, &faq.CreatedDate, &faq.UpdatedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no testimonial found with id %d", id)
		}
		log.Println("Error", err)
		return nil, err
	}

	return &faq, nil
}

func PutFaq(id uint, question, answer string, updated_date time.Time) error {
	result, err := DB.Exec("UPDATE faq SET question=$1, answer=$2, updated_date=$3 where id=$4", question, answer, time.Now(), id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to determine affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("testimonial not found")
	}

	fmt.Println("Update successfull")

	return nil
}

func DeleteFaq(id uint) error {
	result, err := DB.Exec("DELETE FROM faq WHERE id=$1", id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to determine affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("faq found")
	}

	fmt.Println("Delete successfull")

	return nil
}

// Log POST, GET, PUT and DELETE
func PostLog(function, log_message, ip string, created_date, updated_date time.Time) (uint, error) {
	var id uint

	currentTime := time.Now()

	err := DB.QueryRow("INSERT INTO log (function, log_message, ip, created_date, updated_date) VALUES ($1, $2, $3, $4, $5) RETURNING id", function, log_message, ip, currentTime, currentTime).Scan(&id)
	if err != nil {
		return 0, err
	}

	fmt.Println("Post Successful")

	return id, nil
}

func GetLog() ([]models.Log, error) {
	rows, err := DB.Query("SELECT id, function, log_message, ip, created_date, updated_date FROM log")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		err := rows.Scan(&log.ID, &log.Function, &log.LogMessage, &log.IP, &log.CreatedDate, &log.UpdatedDate)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	fmt.Println("Get Successful")

	return logs, nil
}

func GetLogById(id uint) (*models.Log, error) {
	var logs models.Log
	err := DB.QueryRow("SELECT id, function, log_message, ip, created_date, updated_date FROM log WHERE id=$1", id).Scan(&logs.ID, &logs.Function, &logs.LogMessage, &logs.IP, &logs.CreatedDate, &logs.UpdatedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no log found with id %d", id)
		}
		log.Println("Error", err)
		return nil, err
	}

	fmt.Println("Get By Id Successful")

	return &logs, nil
}

func PutLog(id uint, function, log_message, ip string, updated_date time.Time) error {
	result, err := DB.Exec("UPDATE log SET function=$1, log_message=$2, ip=$3, updated_date=$4 WHERE id=$5", function, log_message, ip, time.Now(), id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no log found with id %d", id)
	}

	fmt.Println("Update successful")

	return nil
}

func DeleteLog(id uint) error {
	result, err := DB.Exec("DELETE FROM log WHERE id=$1", id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no log found with id %d", id)
	}

	fmt.Println("Delete Successful")

	return nil
}
