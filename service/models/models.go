package models

import "time"

type Enquiry struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	Message     string    `json:"message"`
	EnquiryType string    `json:"enquiry_type"`
	EnquiryID   string    `json:"enquiry_id"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type Testimonial struct {
	ID          uint      `json:"id"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	Position    string    `json:"position"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type Faq struct {
	ID          uint      `json:"id"`
	Question    string    `json:"question"`
	Answer      string    `json:"answer"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type Log struct {
	ID          uint      `json:"id"`
	Function    string    `json:"function"`
	LogMessage  string    `json:"log_message"`
	IP          string    `json:"ip"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}
