package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Handler struct {}

func CreateHandler() (*Handler, error) {
	return &Handler{}, nil
}

type PayloadInterface interface {
	Valid() error
}

// CreateUserRequestBody defines the format of the /createUser request body.
type CreateUserRequestBody struct {
	// example: 123456
	UserId int `json:"user_id"`
	// example: Joe Smith
	Name string `json:"name"`
	// example: 1990-08-06
	DateOfBirth string `json:"date_of_birth"`
	// example: 1642612034
	CreatedOn int `json:"created_on"`
}

// CreateUserRequestBodyList defines the format of the /createUser request body.
type CreateUserRequestBodyList struct {
	list []CreateUserRequestBody
}

type CreateUserResponseBody struct {
	// example: 123456
	UserId int `json:"user_id"`
	// example: Joe Smith
	Name string `json:"name"`
	// example: Tuesday
	DayOfWeek string `json:"date_of_week"`
	// example: 1642612034
	CreateOn string `json:"create_on_rfc"`
}

func (users CreateUserRequestBodyList) Valid() error {
	for _, user := range users.list {
		if user.UserId == 0 {
			return fmt.Errorf("user_id is required")
		}
		if user.Name == "" {
			return fmt.Errorf("name is required")
		}
		if user.DateOfBirth == "" {
			return fmt.Errorf("date_of_birth is required")
		}
		if user.CreatedOn == 0 {
			return fmt.Errorf("created_on is required")
		}
	}

	return nil
}


func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		// validate auth

		h.handleCreateUserRequest(writer, request)
	default:
		WriteStatusCode(writer, http.StatusMethodNotAllowed)
	}

}

func (h *Handler) handleCreateUserRequest(writer http.ResponseWriter, request *http.Request) {
	payload := &CreateUserRequestBodyList{}

	if err := ReadAndParseRequestBody(writer, request, payload); err != nil {
		return
	}

	var responseList []CreateUserResponseBody

	for _, user := range payload.list {
		dayOfWeek, err := getDayFromDate(user.DateOfBirth)
		if err != nil {
			log.
				WithField("func", "handleCreateUserRequest").
				WithField("statusCode", http.StatusUnprocessableEntity).
				WithError(err).
				Error("Invalid request body")
			WriteStatusCode(writer, http.StatusUnprocessableEntity)
			return
		}

		createdOnRFC3339 := time.Unix(1392899576, 0).Format(time.RFC3339)

		userResponse := CreateUserResponseBody{
			UserId:    user.UserId,
			Name:      user.Name,
			DayOfWeek: dayOfWeek,
			CreateOn:  createdOnRFC3339,
		}

		responseList = append(responseList, userResponse)
	}

	writeCreateUserResponse(writer, responseList)
}

func ReadAndParseRequestBody(w http.ResponseWriter, r *http.Request, payload PayloadInterface) (err error) {
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		log.
			WithField("func", "ReadAndParseRequestBody").
			WithField("statusCode", http.StatusBadRequest).
			WithError(err).
			Error("Error reading request body")

		WriteStatusCode(w, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, payload)
	if err != nil {
		log.
			WithField("func", "ReadAndParseRequestBody").
			WithField("statusCode", http.StatusBadRequest).
			WithError(err).
			Error("Error parsing JSON body")
		WriteStatusCode(w, http.StatusBadRequest)
		return
	}

	err = payload.Valid()

	if err != nil {
		log.
			WithField("func", "ReadAndParseRequestBody").
			WithField("statusCode", http.StatusUnprocessableEntity).
			WithError(err).
			Error("Invalid request body")
		WriteStatusCode(w, http.StatusUnprocessableEntity)
		return
	}

	return
}

func writeCreateUserResponse(w http.ResponseWriter, responseList []CreateUserResponseBody) {
	jsonResp, err := json.Marshal(responseList)

	if err != nil {
		log.WithField("func", "writeCreateUserResponse").
			WithField("statusCode", http.StatusInternalServerError).
			WithField("createUserResponse", responseList).
			WithError(err).
			Error("Failed to create response")
		WriteStatusCode(w, http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}


// WriteStatusCode writes the status code to the header
func WriteStatusCode(writer http.ResponseWriter, statusCode int) int {
	writer.WriteHeader(statusCode)
	return statusCode
}

func getDayFromDate(date string) (string, error) {
	t, err := time.Parse("01-02-2006", date)
	if err != nil {
		return "", fmt.Errorf("could not get weekday")
	}
	return t.Weekday().String(), nil
}