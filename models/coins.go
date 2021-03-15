package models

import (
	"time"
	//"context"
)

type Coin struct {
	Name              string    `json:"Name" bson:"Name"`
	Symbol            string    `json:"Symbol" bson:"Symbol"`
	MarketCap         string    `json:"MarketCap" bson:"MarketCap"`
	Price             string    `json:"Price" bson:"Price"`
	Volume24h         string    `json:"Volume24h" bson:"Volume24h"`
	Change1h          string    `json:"Change1h" bson:"Change1h"`
	Change24h         string    `json:"Change24h" bson:"Change24h"`
	Change7d          string    `json:"Change7d" bson:"Change7d"`
	CirculatingSupply string    `json:"CirculatingSupply" bson:"CirculatingSupply"`
	CreatedAt         time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt         time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
