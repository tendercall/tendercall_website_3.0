package router

import (
	"net/http"

	"tendercall-website.com/main/service/handler"
	"tendercall-website.com/main/service/middleware"
)

func Route() {
	//Enquiry router
	http.Handle("/enquiry", middleware.AuthMiddleware(http.HandlerFunc(handler.EnquiryHandler)))
	http.Handle("/enquirys", middleware.AuthMiddleware(http.HandlerFunc(handler.GetEnquiryByIdHandler)))

	//Testimonal router
	http.Handle("/testimonial", middleware.AuthMiddleware(http.HandlerFunc(handler.TestimonialHandler)))
	http.Handle("/testimonials", middleware.AuthMiddleware(http.HandlerFunc(handler.GetTestimonialByIdHandler)))
}
