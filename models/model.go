/*
 *  MIT License
 *
 * Copyright (c) 2025 Jonas Kaninda
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

package models

import "time"

// **************** Models ************************
type ResponseDto struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Response[T any] struct {
	ResponseDto
	Data T `json:"data,omitempty"`
}

func SuccessResponse[T any](message string, data T) Response[T] {
	return Response[T]{
		ResponseDto: ResponseDto{
			Success: true,
			Message: message,
		},
		Data: data,
	}
}
func ErrorResponse(message string, err error) Response[any] {
	return Response[any]{
		ResponseDto: ResponseDto{
			Success: false,
			Message: message,
			Details: err.Error(),
		},
	}
}
func ErrorResponseData(message string, err error, data any) Response[any] {
	return Response[any]{
		ResponseDto: ResponseDto{
			Success: false,
			Message: message,
			Details: err,
		},
		Data: data,
	}
}

type Book struct {
	Id        int       `json:"id"`
	Title     string    `json:"title" form:"title"  max:"50" required:"true" description:"Book name"`
	Price     int       `json:"price" form:"price" query:"price" yaml:"price" required:"true" description:"Book price"`
	Year      int       `json:"year" form:"year" query:"year" yaml:"year" required:"true" description:"Book year of publication"`
	Author    string    `json:"author" form:"author" query:"author" yaml:"author" required:"false" description:"Book author"`
	Country   string    `json:"country" form:"country" query:"country" yaml:"country" required:"false" description:"Book country of origin"`
	ImageLink string    `json:"imageLink" form:"imageLink" query:"imageLink" yaml:"imageLink" required:"false" description:"Book image link"`
	Language  string    `json:"language" form:"language" query:"language" yaml:"language" required:"false" description:"Book language"`
	Link      string    `json:"link" form:"link" query:"link" yaml:"link" required:"false" description:"Book link"`
	Pages     int       `json:"pages" form:"pages" query:"pages" yaml:"pages" required:"false" description:"Number of pages in the book"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt" query:"createdAt" yaml:"createdAt" required:"false" description:"Book creation date"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt" query:"updatedAt" yaml:"updatedAt" required:"false" description:"Book last update date"`
}

type AuthRequest struct {
	Username string `json:"username" required:"true" description:"Username for authentication" minLength:"3" pattern:"^[a-z]+$"`
	Password string `json:"password" required:"true" description:"Password for authentication"`
}
type Auth struct {
	Token     string   `json:"token,omitempty"`
	ExpiresAt int64    `json:"expires,omitempty"`
	USer      UserInfo `json:"user"`
}
type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type WhoAmIResponse struct {
	Host        string   `json:"host"`
	RealIp      string   `json:"realIp"`
	CurrentUser UserInfo `json:"currentUser"`
}
type WhoAmIRequest struct {
	Email string `header:"current_user_email" format:"email"`
	Name  string `header:"current_user_name"`
	Role  string `header:"current_user_role" enum:"ADMIN,admin,USER,user" doc:"One of ADMIN,admin,USER or user"`
}
type BookResponse = Response[Book]
type UserResponse = Response[UserInfo]
type BooksResponse = Response[[]Book]
type AuthResponse = Response[Auth]

type RequestDto struct {
	Platform string `json:"platform" enum:"web,mobile"`
	AppId    string `json:"appId" enum:"com.mobile,com.web"`
	Details  any    `json:"details,omitempty"`
}

type Request[T any] struct {
	RequestDto
	Fields T `json:"fields,omitempty"`
}
