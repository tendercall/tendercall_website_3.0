package main

import (
	"net/http"

	"tendercall-website.com/main/database"
	"tendercall-website.com/main/service/enquiry/handler"
	"tendercall-website.com/main/service/middleware"
)

func main() {
	database.Initdb()
	//Enquiry router
	http.Handle("/enquiry", middleware.AuthMiddleware(http.HandlerFunc(handler.EnquiryHandler)))
	http.Handle("/enquirys", middleware.AuthMiddleware(http.HandlerFunc(handler.GetEnquiryByIdHandler)))

	http.ListenAndServe(":8080", nil)
}
