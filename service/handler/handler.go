package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nfnt/resize"
	"tendercall-website.com/main/service/helper"
	"tendercall-website.com/main/service/models"
)

// Enquiry Handler
func EnquiryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		PostEnquiryHandler(w, r)
	} else if r.Method == http.MethodGet {
		GetEnquiryHandler(w, r)
	} else if r.Method == http.MethodPut {
		PutEnquiryHandler(w, r)
	} else if r.Method == http.MethodDelete {
		DeleteEnquiryHandler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func PostEnquiryHandler(w http.ResponseWriter, r *http.Request) {
	var enquiry models.Enquiry

	if err := json.NewDecoder(r.Body).Decode(&enquiry); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	id, err := helper.PostEnquiry(enquiry.Email, enquiry.Message, enquiry.EnquiryType, enquiry.EnquiryID, enquiry.CreatedDate, enquiry.UpdatedDate)
	if err != nil {
		http.Error(w, "Error while posting enquiry", http.StatusInternalServerError)
	}

	enquiry.ID = id

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(enquiry)
}

func GetEnquiryHandler(w http.ResponseWriter, r *http.Request) {
	enquiry, err := helper.GetEnquiry()
	if err != nil {
		http.Error(w, "Error while getting enquiry", http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(enquiry)
}

func GetEnquiryByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL query parameter
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Retrieve support from helper by ID
	enquiry, err := helper.GetEnquiryById(uint(id))
	if err != nil {
		if err.Error() == fmt.Sprintf("no enquiry found with id %d", id) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(enquiry)
}

func PutEnquiryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var enquiry models.Enquiry
	if err := json.NewDecoder(r.Body).Decode(&enquiry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the id from the query parameter
	enquiry.ID = uint(id)

	// Set the updated date to the current time
	enquiry.UpdatedDate = time.Now()

	// Update the enquiry in the helper
	err = helper.PutEnquiry(enquiry.ID, enquiry.Email, enquiry.Message, enquiry.EnquiryType, enquiry.EnquiryID, enquiry.UpdatedDate)
	if err != nil {
		if err.Error() == "Enquiry not found" {
			http.Error(w, "Enquiry not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("Failed to update enquiry: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(enquiry); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteEnquiryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = helper.DeleteEnquiry(uint(id))
	if err != nil {
		if err.Error() == "Enquiry not found" {
			http.Error(w, "Enquiry not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("Failed to update enquiry: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// Testimonal Handler
func TestimonialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		PostTestimonialHandler(w, r)
	} else if r.Method == http.MethodGet {
		GetTestimonialHandler(w, r)
	} else if r.Method == http.MethodPut {
		PutTestimonialHandler(w, r)
	} else if r.Method == http.MethodDelete {
		DeleteTestimonialHandler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func PostTestimonialHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max file size
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	var Testimonial models.Testimonial
	Testimonial.Description = r.FormValue("description")
	Testimonial.Name = r.FormValue("name")
	Testimonial.Position = r.FormValue("position")

	// Process uploaded image
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		fmt.Println("Error uploading file:", err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file content", http.StatusInternalServerError)
		fmt.Println("Error reading file content:", err)
		return
	}

	// Resize image if it exceeds 3MB
	if len(fileBytes) > 3*1024*1024 {
		img, _, err := image.Decode(bytes.NewReader(fileBytes))
		if err != nil {
			http.Error(w, "Error decoding image", http.StatusInternalServerError)
			fmt.Println("Error decoding image:", err)
			return
		}

		newImage := resize.Resize(800, 0, img, resize.Lanczos3)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, newImage, nil)
		if err != nil {
			http.Error(w, "Error encoding compressed image", http.StatusInternalServerError)
			fmt.Println("Error encoding compressed image:", err)
			return
		}
		fileBytes = buf.Bytes()
	}

	// Upload image to AWS S3
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Replace with your AWS region
		Credentials: credentials.NewStaticCredentials(
			"AKIAYS2NVN4MBSHP33FF",                     // Replace with your AWS access key ID
			"aILySGhiQAB7SaFnqozcRZe1MhZ0zNODLof2Alr4", // Replace with your AWS secret access key
			""), // Optional token, leave blank if not using
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	svc := s3.New(sess)
	imageKey := fmt.Sprintf("TestimonalImages/%d.jpg", time.Now().Unix()) // Adjust key as needed
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("tendercall-db"), // Replace with your S3 bucket name
		Key:    aws.String(imageKey),
		Body:   bytes.NewReader(fileBytes),
	})
	if err != nil {
		log.Printf("Failed to upload image to S3: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Construct imageURL assuming it's from your S3 bucket
	imageURL := fmt.Sprintf("https://tendercall-db.s3.amazonaws.com/%s", imageKey)

	id, err := helper.PostTestimonial(imageURL, Testimonial.Description, Testimonial.Name, Testimonial.Position, Testimonial.CreatedDate, Testimonial.UpdatedDate)
	if err != nil {
		http.Error(w, "Error while posting Testimonal", http.StatusInternalServerError)
	}

	Testimonial.ID = id

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(Testimonial)
}

func GetTestimonialHandler(w http.ResponseWriter, r *http.Request) {
	Testimonial, err := helper.GetTestimonial()
	if err != nil {
		http.Error(w, "Error while getting Testimonial", http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(Testimonial)
}

func GetTestimonialByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL query parameter
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Retrieve support from helper by ID
	testimonial, err := helper.GetTestimonialById(uint(id))
	if err != nil {
		if err.Error() == fmt.Sprintf("no Testimonial found with id %d", id) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(testimonial)
}

func PutTestimonialHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var Testimonial models.Testimonial
	if err := json.NewDecoder(r.Body).Decode(&Testimonial); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the id from the query parameter
	Testimonial.ID = uint(id)

	// Set the updated date to the current time
	Testimonial.UpdatedDate = time.Now()

	// Update the testimonial in the helper
	err = helper.PutTestimonial(Testimonial.ID, Testimonial.Image, Testimonial.Description, Testimonial.Name, Testimonial.Position, Testimonial.UpdatedDate)
	if err != nil {
		if err.Error() == "Testimonial not found" {
			http.Error(w, "Testimonial not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("Failed to update Testimonial: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Testimonial); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteTestimonialHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = helper.DeleteTestimonial(uint(id))
	if err != nil {
		if err.Error() == "Testimonial not found" {
			http.Error(w, "Testimonial not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("Failed to delete Testimonial: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// Faq Handler
func FaqHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		PostFaqHandler(w, r)
	} else if r.Method == http.MethodGet {
		GetFaqHandler(w, r)
	} else if r.Method == http.MethodPut {
		PutFaqHandler(w, r)
	} else if r.Method == http.MethodDelete {
		DeleteFaqHandler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func PostFaqHandler(w http.ResponseWriter, r *http.Request) {
	var faq models.Faq

	if err := json.NewDecoder(r.Body).Decode(&faq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	id, err := helper.PostFaq(faq.Question, faq.Answer, faq.CreatedDate, faq.UpdatedDate)
	if err != nil {
		http.Error(w, "Error while posting enquiry", http.StatusInternalServerError)
	}

	faq.ID = id

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(faq)
}

func GetFaqHandler(w http.ResponseWriter, r *http.Request) {
	faq, err := helper.GetFaq()
	if err != nil {
		http.Error(w, "Error while getting Testimonial", http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(faq)
}

func GetFaqByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL query parameter
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Retrieve support from helper by ID
	faq, err := helper.GetFaqById(uint(id))
	if err != nil {
		if err.Error() == fmt.Sprintf("no faq found with id %d", id) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(faq)
}

func PutFaqHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var faq models.Faq
	if err := json.NewDecoder(r.Body).Decode(&faq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the id from the query parameter
	faq.ID = uint(id)

	// Set the updated date to the current time
	faq.UpdatedDate = time.Now()

	// Update the faq in the helper
	err = helper.PutFaq(faq.ID, faq.Question, faq.Answer, faq.UpdatedDate)
	if err != nil {
		if err.Error() == "Faq not found" {
			http.Error(w, "Faq not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("Failed to update Faq: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(faq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteFaqHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = helper.DeleteFaq(uint(id))
	if err != nil {
		if err.Error() == "Faq not found" {
			http.Error(w, "Faq not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("Failed to delete faq: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// Log Handler
func LogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		PostLogHandler(w, r)
	} else if r.Method == http.MethodGet {
		GetLogHandler(w, r)
	} else if r.Method == http.MethodPut {
		PutLogHandler(w, r)
	} else if r.Method == http.MethodDelete {
		DeleteLogHandler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func PostLogHandler(w http.ResponseWriter, r *http.Request) {
	var log models.Log

	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	id, err := helper.PostLog(log.Function, log.LogMessage, log.IP, log.CreatedDate, log.UpdatedDate)
	if err != nil {
		http.Error(w, "Error while posting enquiry", http.StatusInternalServerError)
	}

	log.ID = id

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(log)
}

func GetLogHandler(w http.ResponseWriter, r *http.Request) {
	log, err := helper.GetLog()
	if err != nil {
		http.Error(w, "Error while getting Testimonial", http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(log)
}

func GetLogByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL query parameter
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Retrieve support from helper by ID
	log, err := helper.GetLogById(uint(id))
	if err != nil {
		if err.Error() == fmt.Sprintf("no log found with id %d", id) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(log)
}

func PutLogHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var log models.Log
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the id from the query parameter
	log.ID = uint(id)

	// Set the updated date to the current time
	log.UpdatedDate = time.Now()

	// Update the log in the helper
	err = helper.PutLog(log.ID, log.Function, log.LogMessage, log.IP, log.UpdatedDate)
	if err != nil {
		if err.Error() == "Log not found" {
			http.Error(w, "Log not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("Failed to update log: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(log); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteLogHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = helper.DeleteLog(uint(id))
	if err != nil {
		if err.Error() == "Log not found" {
			http.Error(w, "Log not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("Failed to delete Log: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
