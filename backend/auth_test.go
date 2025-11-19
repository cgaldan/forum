package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal("Failed to open test database:", err)
	}

	if err = initDatabase(); err != nil {
		t.Fatal("Failed to initialize test database:", err)
	}
}

func teardownTestDB() {
	if db != nil {
		db.Close()
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestRegisterHandler(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	tests := []struct {
		name           string
		payload        RegisterRequest
		expectedStatus bool
		expectedError  string
	}{
		{
			name: "Valid registration",
			payload: RegisterRequest{
				Nickname:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Age:       25,
				Gender:    "male",
			},
			expectedStatus: true,
		},
		{
			name: "Short nickname",
			payload: RegisterRequest{
				Nickname:  "ab",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Age:       25,
				Gender:    "male",
			},
			expectedStatus: false,
			expectedError:  "nickname must be at least 3 characters",
		},
		{
			name: "Invalid email",
			payload: RegisterRequest{
				Nickname:  "testuser",
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Age:       25,
				Gender:    "male",
			},
			expectedStatus: false,
			expectedError:  "valid email is required",
		},
		{
			name: "Short password",
			payload: RegisterRequest{
				Nickname:  "testuser",
				Email:     "test@example.com",
				Password:  "pass",
				FirstName: "Test",
				LastName:  "User",
				Age:       25,
				Gender:    "male",
			},
			expectedStatus: false,
			expectedError:  "password must be at least 6 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			registerHandler(w, req)

			var response AuthResponse
			json.NewDecoder(w.Body).Decode(&response)

			if response.Success != tt.expectedStatus {
				t.Errorf("Expected success=%v, got success=%v", tt.expectedStatus, response.Success)
			}

			if !tt.expectedStatus && response.Message != tt.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedError, response.Message)
			}

			if tt.expectedStatus && response.Token == "" {
				t.Error("Expected token to be returned on successful registration")
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// First, register a user
	registerPayload := RegisterRequest{
		Nickname:  "logintest",
		Email:     "login@example.com",
		Password:  "password123",
		FirstName: "Login",
		LastName:  "Test",
		Age:       30,
		Gender:    "female",
	}

	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	registerHandler(w, req)

	tests := []struct {
		name           string
		payload        LoginRequest
		expectedStatus bool
	}{
		{
			name: "Login with nickname",
			payload: LoginRequest{
				Identifier: "logintest",
				Password:   "password123",
			},
			expectedStatus: true,
		},
		{
			name: "Login with email",
			payload: LoginRequest{
				Identifier: "login@example.com",
				Password:   "password123",
			},
			expectedStatus: true,
		},
		{
			name: "Invalid password",
			payload: LoginRequest{
				Identifier: "logintest",
				Password:   "wrongpassword",
			},
			expectedStatus: false,
		},
		{
			name: "Non-existent user",
			payload: LoginRequest{
				Identifier: "nonexistent",
				Password:   "password123",
			},
			expectedStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			loginHandler(w, req)

			var response AuthResponse
			json.NewDecoder(w.Body).Decode(&response)

			if response.Success != tt.expectedStatus {
				t.Errorf("Expected success=%v, got success=%v", tt.expectedStatus, response.Success)
			}

			if tt.expectedStatus && response.Token == "" {
				t.Error("Expected token to be returned on successful login")
			}
		})
	}
}

func TestValidateRegistration(t *testing.T) {
	tests := []struct {
		name      string
		req       RegisterRequest
		expectErr bool
	}{
		{
			name: "Valid registration",
			req: RegisterRequest{
				Nickname:  "validuser",
				Email:     "valid@example.com",
				Password:  "password123",
				FirstName: "Valid",
				LastName:  "User",
				Age:       25,
				Gender:    "male",
			},
			expectErr: false,
		},
		{
			name: "Age too young",
			req: RegisterRequest{
				Nickname:  "younguser",
				Email:     "young@example.com",
				Password:  "password123",
				FirstName: "Young",
				LastName:  "User",
				Age:       10,
				Gender:    "male",
			},
			expectErr: true,
		},
		{
			name: "Age too old",
			req: RegisterRequest{
				Nickname:  "olduser",
				Email:     "old@example.com",
				Password:  "password123",
				FirstName: "Old",
				LastName:  "User",
				Age:       150,
				Gender:    "male",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegistration(tt.req)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error=%v, got error=%v", tt.expectErr, err != nil)
			}
		})
	}
}
