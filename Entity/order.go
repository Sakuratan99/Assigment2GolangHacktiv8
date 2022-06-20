package entity

import "time"

type Order struct {
	Order_id      int       `json:"orderid"`
	Customer_Name string    `json:"customerName"`
	Ordered_At    time.Time `json:"ordereAt"`
	Item          []Item    `json:"items"`
}

type Item struct {
	Item_Id     int    `json:"lineItemId"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Order_Id    int    `json:"orderid"`
}
