package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetServiceID(t *testing.T) {
	mockResponse := `{
		"data": {
			"account": {
				"services": {
					"nodes": [
						{"id": "1", "name": "your_service_name"}
					]
				}
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	defer func() { baseURL = originalBaseURL }()
	baseURL = server.URL

	serviceID, err := getServiceID("your_service_name")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedID := "1"
	if serviceID != expectedID {
		t.Errorf("Expected serviceID to be %s, got %s", expectedID, serviceID)
	}
}

func TestGetTags(t *testing.T) {
	mockResponse := `{
		"data": {
			"account": {
				"service": {
					"tags": [
						{"id": "1", "name": "tag1"},
						{"id": "2", "name": "tag2"}
					]
				}
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	defer func() { baseURL = originalBaseURL }()
	baseURL = server.URL

	tags, err := getTags("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedTags := []string{"tag1", "tag2"}
	for i, tag := range tags {
		if tag != expectedTags[i] {
			t.Errorf("Expected tag to be %s, got %s", expectedTags[i], tag)
		}
	}
}

func TestGetOnCallOrManager(t *testing.T) {
	mockResponse := `{
		"data": {
			"account": {
				"service": {
					"onCallRotation": {
						"onCallUser": {"name": "on_call_user"}
					},
					"owner": {
						"manager": {"name": "manager_name"}
					}
				}
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	defer func() { baseURL = originalBaseURL }()
	baseURL = server.URL

	name, err := getOnCallOrManager("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedName := "on_call_user"
	if name != expectedName {
		t.Errorf("Expected name to be %s, got %s", expectedName, name)
	}

	// Test case where there is no on-call user, only manager
	mockResponseNoOnCall := `{
		"data": {
			"account": {
				"service": {
					"onCallRotation": {
						"onCallUser": {"name": ""}
					},
					"owner": {
						"manager": {"name": "manager_name"}
					}
				}
			}
		}
	}`

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponseNoOnCall))
	}))
	defer server.Close()

	name, err = getOnCallOrManager("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedName = "manager_name"
	if name != expectedName {
		t.Errorf("Expected name to be %s, got %s", expectedName, name)
	}
}
