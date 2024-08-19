package httpTransport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"goProject-2024/models"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type StudentService interface {
	GetStudent(ctx context.Context, ID string) (models.Student, error)
	PostStudent(ctx context.Context, std models.Student) (models.Student, error)
	UpdateStudent(ctx context.Context, ID string, newStd models.Student) (models.Student, error)
	DeleteStudent(ctx context.Context, ID string) error
	ReadyCheck(ctx context.Context) error
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Invalid request payload"})
		return
	}

	if creds.UserID == "admin" && creds.Password == "mindPalace@23" {
		token, err := generateJWTToken(creds.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Message: "Error generating token"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Message: "Invalid user ID or password"})
	}
}

func extractUserIDFromToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}

func (h *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	std, err := h.Service.GetStudent(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(std); err != nil {
		panic(err)
	}
}

type PostStudentRequest struct {
	FName       string    `json:"first_name" validate:"required"`
	LName       string    `json:"last_name" validate:"required"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
	Email       string    `json:"email" validate:"required"`
	Address     string    `json:"address" validate:"required"`
	Gender      string    `json:"gender" validate:"required"`
	Age         int       `json:"age" validate:"required"`
}

// PostStudent - adds a new student
/* func (h *Handler) PostStudent(w http.ResponseWriter, r *http.Request) {
	var postStdReq PostStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&postStdReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(postStdReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	std := models.Student{
		FName:       postStdReq.FName,
		LName:       postStdReq.LName,
		DateOfBirth: postStdReq.DateOfBirth,
		Email:       postStdReq.Email,
		Address:     postStdReq.Address,
		Gender:      postStdReq.Gender,
		Age:         postStdReq.Age,
		CreatedBy:   postStdReq.CreatedBy,
		CreatedOn:   postStdReq.CreatedOn,
	}
	std, err = h.Service.PostStudent(r.Context(), std)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(std); err != nil {
		panic(err)
	}
} */
func (h *Handler) PostStudent(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Message: "Unauthorized"})
		return
	}

	var postStdReq PostStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&postStdReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(postStdReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	std := models.Student{
		FName:       postStdReq.FName,
		LName:       postStdReq.LName,
		DateOfBirth: postStdReq.DateOfBirth,
		Email:       postStdReq.Email,
		Address:     postStdReq.Address,
		Gender:      postStdReq.Gender,
		Age:         postStdReq.Age,
		CreatedBy:   userID,
		CreatedOn:   time.Now(),
	}

	std, err = h.Service.PostStudent(r.Context(), std)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(std); err != nil {
		panic(err)
	}
}

type UpdateStudentRequest struct {
	FName       *string    `json:"first_name" validate:"omitempty"`
	LName       *string    `json:"last_name" validate:"omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth" validate:"omitempty"`
	Email       *string    `json:"email" validate:"omitempty"`
	Address     *string    `json:"address" validate:"omitempty"`
	Gender      *string    `json:"gender" validate:"omitempty"`
	Age         *int       `json:"age" validate:"omitempty"`
}

// UpdateStudent - updates a student by ID
/* func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]

	var updateStdRequest UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&updateStdRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(updateStdRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve the existing student record from the database
	existingStudent, err := h.Service.GetStudent(r.Context(), studentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update only the fields that are provided in the request
	if updateStdRequest.FName != nil {
		existingStudent.FName = *updateStdRequest.FName
	}
	if updateStdRequest.LName != nil {
		existingStudent.LName = *updateStdRequest.LName
	}
	if updateStdRequest.DateOfBirth != nil {
		existingStudent.DateOfBirth = *updateStdRequest.DateOfBirth
	}
	if updateStdRequest.Email != nil {
		existingStudent.Email = *updateStdRequest.Email
	}
	if updateStdRequest.Address != nil {
		existingStudent.Address = *updateStdRequest.Address
	}
	if updateStdRequest.Gender != nil {
		existingStudent.Gender = *updateStdRequest.Gender
	}
	if updateStdRequest.Age != nil {
		existingStudent.Age = *updateStdRequest.Age
	}
	if updateStdRequest.UpdatedBy != nil {
		existingStudent.UpdatedBy = *updateStdRequest.UpdatedBy
	}
	existingStudent.UpdatedOn = time.Now()

	// Save the updated student record to the database
	updatedStudent, err := h.Service.UpdateStudent(r.Context(), studentID, existingStudent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedStudent); err != nil {
		panic(err)
	}
} */

func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]

	var updateStdRequest UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&updateStdRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(updateStdRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := extractUserIDFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existingStudent, err := h.Service.GetStudent(r.Context(), studentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if updateStdRequest.FName != nil {
		existingStudent.FName = *updateStdRequest.FName
	}
	if updateStdRequest.LName != nil {
		existingStudent.LName = *updateStdRequest.LName
	}
	if updateStdRequest.DateOfBirth != nil {
		existingStudent.DateOfBirth = *updateStdRequest.DateOfBirth
	}
	if updateStdRequest.Email != nil {
		existingStudent.Email = *updateStdRequest.Email
	}
	if updateStdRequest.Address != nil {
		existingStudent.Address = *updateStdRequest.Address
	}
	if updateStdRequest.Gender != nil {
		existingStudent.Gender = *updateStdRequest.Gender
	}
	if updateStdRequest.Age != nil {
		existingStudent.Age = *updateStdRequest.Age
	}

	existingStudent.UpdatedBy = userID
	existingStudent.UpdatedOn = time.Now()

	updatedStudent, err := h.Service.UpdateStudent(r.Context(), studentID, existingStudent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedStudent); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteStudent(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
