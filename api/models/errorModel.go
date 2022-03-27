package models

// Custom Error Model
type AppError struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// Implementing the error interface method to make VendError of type error
func (me AppError) Error() string {
	return me.ErrorMessage
}
