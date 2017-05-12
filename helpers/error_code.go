package helpers

// const
const (

	//Common error
	ErrorMissingParameter = 10

	//Error about User
	ErrorExistEmail            = 101
	ErrorExistUsername         = 102
	ErrorExistEmailAndUsername = 103
	ErrorExistStudentCode      = 113
)

//ErrorDetail struct for a Error
type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//NewErrorDetail func to create  and return a new ErrorDetail
func NewErrorDetail(c int, m string) ErrorDetail {
	return ErrorDetail{Code: c, Message: m}
}

//Errors struct include a or a few errors
type Errors struct {
	Code         int           `json:"code"`
	Message      string        `json:"message"`
	ErrorDetails []ErrorDetail `json:"error_details"`
}
