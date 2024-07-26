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
