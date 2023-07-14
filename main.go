package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]int)

type Item struct {
	ShortDescription string      `json:"shortDescription" binding:"required"`
	Price            json.Number `json:"price" binding:"required,numeric,gte=0"`
}

type Id struct {
	Id string `json:"id"`
}

type Points struct {
	Points int `json:"points"`
}

type Receipt struct {
	Retailer     string      `json:"retailer" binding:"required"`
	PurchaseDate string      `json:"purchaseDate" binding:"required"`
	PurchaseTime string      `json:"purchaseTime" binding:"required"`
	Items        []Item      `json:"items" binding:"required"`
	Total        json.Number `json:"total" binding:"required,numeric,gte=0"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	receipts := r.Group("/receipts")

	receipts.POST("/process", func(c *gin.Context) {

		var receipt Receipt
		err := c.Bind(&receipt)

		if err == nil {
			fmt.Println(receipt)
			id := uuid.New().String()
			db[id], err = calculatePoints(&receipt)

			if err == nil {
				c.JSON(http.StatusOK, Id{id})
			} else {
				c.Status(http.StatusBadRequest)
			}

		} else {
			fmt.Println(err)
			c.Status(http.StatusBadRequest)
		}
	})

	receipts.GET("/:id/points", func(c *gin.Context) {
		id := c.Param("id")
		points, ok := db[id]

		if ok {
			c.JSON(http.StatusOK, Points{points})
		} else {
			c.Status(http.StatusBadRequest)
		}
	})

	return r
}

// One point for every alphanumeric character in the retailer name.
func calculateRetailer(retailer string) int {
	points := 0
	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points += 1
		}
	}
	return points
}

// Perform calculations on the 'total' field
func calculateTotal(totalFloat float64) int {
	points := 0
	if int(totalFloat*100)%100 == 0 { // 50 points if the total is a round dollar amount with no cents.
		points += 50
	}
	if int(totalFloat*100)%25 == 0 { // 25 points if the total is a multiple of 0.25.
		points += 25
	}
	return points
}

// 5 points for every two items on the receipt.
func calculateEveryTwo(itemsLength int) int {
	return (itemsLength / 2) * 5
}

// Perform calculations on each item in the 'items' field
func calculateTrimmedLength(items []Item) (int, error) {
	points := 0
	for _, item := range items {
		shortDescription := item.ShortDescription
		shortDescriptionLen := len(strings.TrimSpace(shortDescription))

		if shortDescriptionLen%3 == 0 { // If the trimmed length of the item description is a multiple of 3...
			totalFloat, err := item.Price.Float64()
			if err == nil {
				points += int(math.Ceil(totalFloat * 0.2)) // multiply the price by 0.2 and round up...
			} else {
				return -1, err
			}
		}
	}
	return points, nil
}

// 6 points if the day in the purchase date is odd.
func calculateOddPurchaseDate(purchaseDate string) (int, error) {
	lastDigitOfDate := purchaseDate[len(purchaseDate)-1:]
	parsedLastDigit, err := strconv.Atoi(lastDigitOfDate)
	if err == nil {
		if parsedLastDigit%2 == 1 {
			return 6, nil
		}
		return 0, nil
	}
	return -1, err
}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func calculatePurchaseTime(purchaseTime string) (int, error) {
	strippedPurchaseTime := strings.Replace(purchaseTime, ":", "", -1)
	purchaseTimeInt, err := strconv.Atoi(strippedPurchaseTime)
	if err == nil {
		if purchaseTimeInt > 1400 && purchaseTimeInt < 1600 {
			return 10, nil
		}
		return 0, nil
	}
	return -1, err
}

func calculatePoints(receipt *Receipt) (int, error) {
	points := 0
	retailer := receipt.Retailer
	total := receipt.Total
	items := receipt.Items
	purchaseDate := receipt.PurchaseDate
	purchaseTime := receipt.PurchaseTime

	// Add points based on the retailer name
	points += calculateRetailer(retailer)
	fmt.Println(points, "points")

	// Add points based on the total field of the receipt
	totalFloat, err := total.Float64()
	if err == nil {
		points += calculateTotal(totalFloat)
	} else {
		return -1, err
	}
	fmt.Println(points, "points")

	// Add points based on the number of items on the receipt
	points += calculateEveryTwo(len(items))
	fmt.Println(points, "points")

	// Add points based on trimmed item descriptions
	trimmedLengthPoints, err := calculateTrimmedLength(items)
	if err == nil {
		points += trimmedLengthPoints
	} else {
		println(err)
		return -1, err
	}
	fmt.Println(points, "points")

	// Add points based on purchase date
	purchaseDatePoints, err := calculateOddPurchaseDate(purchaseDate)
	if err == nil {
		points += purchaseDatePoints
	} else {
		println(err)
		return -1, err
	}
	fmt.Println(points, "points")

	// Add points based on purchase time
	purchaseTimePoints, err := calculatePurchaseTime(purchaseTime)
	if err == nil {
		points += purchaseTimePoints
	} else {
		println(err)
		return -1, err
	}
	fmt.Println(points, "points")

	return points, nil
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
