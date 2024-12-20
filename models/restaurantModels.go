package models

import (
	"github.com/lib/pq"
	"time"
)
import "github.com/shopspring/decimal"

// Restaurant represents a Restaurant in the handlers
type Restaurant struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	OwnerID   uint      `db:"owner_id"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	//Dishes    []Dish    `db:"dishes"`
	Dishes JSONDish `db:"dishes"`
}

type JSONDish []Dish
type RestaurantRequestBody struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Dish represents a dish in the handlers
type Dish struct {
	ID           uint            `db:"id"`
	Name         string          `db:"name"`
	Price        decimal.Decimal `db:"price"`
	RestaurantID uint            `db:"restaurant_id"`
	OwnerID      uint            `db:"owner_id"`
	Tags         pq.StringArray  `db:"tags" swaggertype:"array,string"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at"`
}
type DishRequestBody struct {
	Name         string          `json:"name"`
	Price        decimal.Decimal `json:"price"`
	RestaurantID uint            `json:"restaurant_id"`
	Tags         []string        `json:"tags"`
}
