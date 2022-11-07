package api

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserHappyPath(t *testing.T) {
	h, err := CreateUserHandler()
	if err != nil {
		t.Fatalf("error creating handler")
	}

	var body = []CreateUserRequestBody{}

	body = append(body, CreateUserRequestBody{
		UserId: 1234,
		DateOfBirth: "1998-02-03",
		Name: "John",
		CreatedOn: 1234578564,
	},
	)

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/createUser", &b)
	h.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// TODO: read response body into a CreateUserResponseBody and assert the value each field
	assert.Equal(t, "[{\"user_id\":1234,\"name\":\"John\",\"date_of_week\":\"Thursday\",\"create_on_rfc\":\"2014-02-20T07:32:56-05:00\"}]", responseRecorder.Body.String())
}


func TestCreateNoUserId(t *testing.T) {
	h, err := CreateUserHandler()
	if err != nil {
		t.Fatalf("error creating handler")
	}

	var body = []CreateUserRequestBody{}

	body = append(body, CreateUserRequestBody{
		DateOfBirth: "1998-02-03",
		Name: "John",
		CreatedOn: 1234578564,
	},
	)

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/createUser", &b)
	h.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusUnprocessableEntity, responseRecorder.Code)
	assert.Equal(t, "user_id is required", responseRecorder.Body.String())
}

func TestCreateNoName(t *testing.T) {
	h, err := CreateUserHandler()
	if err != nil {
		t.Fatalf("error creating handler")
	}

	var body = []CreateUserRequestBody{}

	body = append(body, CreateUserRequestBody{
		UserId: 1234,
		DateOfBirth: "1998-02-03",
		CreatedOn: 1234578564,
	},
	)

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/createUser", &b)
	h.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusUnprocessableEntity, responseRecorder.Code)
	assert.Equal(t, "name is required", responseRecorder.Body.String())
}

// TODO: add tests for other validation errors

func TestBadEntity(t *testing.T) {
	h, err := CreateUserHandler()
	if err != nil {
		t.Fatalf("error creating handler")
	}

	var body = CreateUserRequestBody{
		Name: "Joe Smith",
		UserId: 1234,
		DateOfBirth: "1998-02-03",
		CreatedOn: 1234578564,
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/createUser", &b)
	h.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}