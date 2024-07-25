package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"tendercall-website.com/main/service/enquiry/repository"
	"tendercall-website.com/main/service/models"
)

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

	id, err := repository.PostEnquiry(enquiry.Email, enquiry.Message, enquiry.EnquiryType, enquiry.EnquiryID, enquiry.CreatedDate, enquiry.UpdatedDate)
	if err != nil {
		http.Error(w, "Error while posting enquiry", http.StatusInternalServerError)
	}

	enquiry.ID = id

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(enquiry)
}

func GetEnquiryHandler(w http.ResponseWriter, r *http.Request) {
	enquiry, err := repository.GetEnquiry()
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

	// Retrieve support from repository by ID
	enquiry, err := repository.GetEnquiryById(uint(id))
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

	// Update the Offer in the repository
	err = repository.PutEnquiry(enquiry.ID, enquiry.Email, enquiry.Message, enquiry.EnquiryType, enquiry.EnquiryID, enquiry.UpdatedDate)
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

	err = repository.DeleteEnquiry(uint(id))
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
