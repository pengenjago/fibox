package response

import "github.com/gofiber/fiber/v3"

// Response is the standard API response structure
type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
}

type Pagination struct {
	PageNo      int `json:"pageNo"`
	PageSize    int `json:"pageSize"`
	PageTotal   int `json:"pageTotal"`
	TotalRecord int `json:"totalRecord"`
}

// Success sends a success response
func Success(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPagination sends a success response with pagination info
func SuccessWithPagination(c fiber.Ctx, message string, data interface{}, pageNo, pageSize, totalRecord int) error {
	totalPage := 0

	if pageSize > 0 {
		totalPage = totalRecord / pageSize
		if totalRecord%pageSize > 0 {
			totalPage++
		}
	}

	if pageNo > totalPage {
		pageNo = totalPage
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: Pagination{
			PageNo:    pageNo,
			PageSize:  pageSize,
			PageTotal: totalPage,
		},
	})
}

// Created sends a created response
func Created(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a bad request error response
func BadRequest(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Message: message,
	})
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Message: message,
	})
}

// Forbidden sends a forbidden error response
func Forbidden(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success: false,
		Message: message,
	})
}

// NotFound sends a not found error response
func NotFound(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Message: message,
	})
}

// InternalError sends an internal server error response
func InternalError(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Message: message,
		Success: false,
	})
}
