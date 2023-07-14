package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Runs a post and get for the given request body, expecting a 200 OK and returns the number of points
func runPostAndGetOk(t *testing.T, testResourcePath string) int {
	router := setupRouter()

	w := httptest.NewRecorder()
	file, _ := ioutil.ReadFile(testResourcePath)

	req, _ := http.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(file))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var id Id
	json.Unmarshal(w.Body.Bytes(), &id)
	uuid := id.Id

	req2, _ := http.NewRequest(http.MethodGet, "/receipts/"+uuid+"/points", bytes.NewBuffer(file))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var points Points
	assert.Equal(t, http.StatusOK, w2.Code)
	json.Unmarshal(w2.Body.Bytes(), &points)
	return points.Points
}

// Runs a post with a malformed body and expects a 400 back
func runPostAndGetBadRequest(t *testing.T, testResourcePath string) {
	router := setupRouter()

	w := httptest.NewRecorder()
	file, _ := ioutil.ReadFile(testResourcePath)

	req, _ := http.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(file))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Nil(t, w.Body.Bytes())
}

// Provided test case 1
func TestProvidedTestCase1(t *testing.T) {
	points := runPostAndGetOk(t, "testResources/test1.json")
	assert.Equal(t, 28, points)
}

// Provided test case 2
func TestProvidedTestCase2(t *testing.T) {
	points := runPostAndGetOk(t, "testResources/test2.json")
	assert.Equal(t, 109, points)
}

// Test for 99 points
func Test99Points(t *testing.T) {
	points := runPostAndGetOk(t, "testResources/99points.json")
	assert.Equal(t, 99, points)
}

// Test no items, this should return 0 points
func TestNoItems(t *testing.T) {
	points := runPostAndGetOk(t, "testResources/noItems.json")
	assert.Equal(t, 0, points)
}

// Test an incomplete JSON body
func TestBadJson(t *testing.T) {
	runPostAndGetBadRequest(t, "testResources/malformedInput.json")
}

// Test a request with invalid prices
func TestInvalidPrices(t *testing.T) {
	runPostAndGetBadRequest(t, "testResources/testPricesAreNotNumbers.json")
}

// Test a request with a missing field
func TestMissingPurchaseTime(t *testing.T) {
	runPostAndGetBadRequest(t, "testResources/missingPurchaseTime.json")
}

// Test getting an ID that doesn't exist
func TestGetForInvalidId(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/receipts/"+uuid.New().String()+"/points", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Nil(t, w.Body.Bytes())
}
