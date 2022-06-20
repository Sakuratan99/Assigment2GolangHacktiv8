package handler

import (
	entity "Assigment2Golang/Entity"
	"context"
	_ "context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	_ "time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type ItemHandlerInterface interface {
	ItemsHandler(w http.ResponseWriter, r *http.Request)
}

type ItemHandler struct {
	db *sql.DB
}

// ItemsHandler implements ItemHandlerInterface

func NewItemHandler(db *sql.DB) ItemHandlerInterface {
	return &ItemHandler{db: db}
}

var (
	db *sql.DB

	err error
)

func (h *ItemHandler) ItemsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Println(id)
	switch r.Method {
	case http.MethodGet:
		h.getItemsHandler(w, r)
	case http.MethodPost:
		//users
		h.createItemsHandler(w, r)
	case http.MethodPut:
		//users/{id}
		h.UpdateOrderById(w, r, id)
	case http.MethodDelete:
		h.DeleteOrderbyId(w, r, id)
	}
}

func (h *ItemHandler) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	queryString := `select
		o.order_id as order_id
		,o.customer_name
		,o.ordered_at
		,json_agg(json_build_object(
			'item_id',i.item_id
			,'item_code',i.item_code
			,'description',i.description
			,'quantity',i.quantity
			,'order_id',i.order_id
		)) as items
	from orders o join items i
	on o.order_id = i.order_id
	group by o.order_id`
	rows, err := h.db.QueryContext(ctx, queryString)
	if err != nil {
		fmt.Println("query row error", err)
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		var o entity.Order
		var itemsStr string
		if serr := rows.Scan(&o.Order_id, &o.Customer_Name, &o.Ordered_At, &itemsStr); serr != nil {
			fmt.Println("Scan error", serr)
		}
		var items []entity.Item
		if err := json.Unmarshal([]byte(itemsStr), &items); err != nil {
			fmt.Errorf("Error when parsing items")
		} else {
			o.Item = append(o.Item, items...)
		}
		orders = append(orders, &o)
	}
	jsonData, _ := json.Marshal(&orders)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *ItemHandler) createItemsHandler(w http.ResponseWriter, r *http.Request) {
	var newOrder entity.Order
	json.NewDecoder(r.Body).Decode(&newOrder)
	fmt.Println(newOrder)
	sqlStatment := `insert into orders
	(Customer_Name,Ordered_At)
	values ($1 ,$2) returning order_id ;`
	ctx := context.Background()
	var id int
	err := h.db.QueryRowContext(ctx, sqlStatment, newOrder.Customer_Name, time.Now()).Scan(&id)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(newOrder.Item); i++ {
		var items entity.Item
		items.ItemCode = newOrder.Item[i].ItemCode
		items.Description = newOrder.Item[i].Description
		items.Quantity = newOrder.Item[i].Quantity
		query := `insert into items 
		(item_code,description,quantity,order_id)
		values ($1,$2,$3,$4) `

		_, err := h.db.Exec(query, items.ItemCode, items.Description, items.Quantity, id)
		if err != nil {
			panic(nil)
		}
	}

	w.Write([]byte(fmt.Sprint("Create user rows ")))
	return
}

func (h *ItemHandler) UpdateOrderById(w http.ResponseWriter, r *http.Request, id string) {
	if id != "" { // get by id
		var newOrder entity.Order
		json.NewDecoder(r.Body).Decode(&newOrder)
		sqlstatment := `
		update orders set customer_name = $1 , ordered_at = $2 where order_id = $3;`

		res, err := h.db.Exec(sqlstatment,
			newOrder.Customer_Name,
			time.Now(),
			id,
		)
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(newOrder.Item); i++ {
			var items entity.Item
			items.Item_Id = newOrder.Item[i].Item_Id
			items.ItemCode = newOrder.Item[i].ItemCode
			items.Description = newOrder.Item[i].Description
			items.Quantity = newOrder.Item[i].Quantity
			query := `update items set item_code = $1, description = $2, quantity = $3 where order_id = $4 and item_id = $5`

			_, err := h.db.Exec(query, items.ItemCode, items.Description, items.Quantity, id, items.Item_Id)
			if err != nil {
				panic(nil)
			}
		}
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}

		w.Write([]byte(fmt.Sprint("User  update ", count)))
		return
	}
}

func (h *ItemHandler) DeleteOrderbyId(w http.ResponseWriter, r *http.Request, id string) {
	sqlstament := `DELETE from orders where Order_id = $1;`
	if idInt, err := strconv.Atoi(id); err == nil {
		sqlstament2 := `DELETE from items where Order_id = $1;`
		_, err2 := h.db.Exec(sqlstament2, idInt)
		if err2 != nil {
			panic(err)
		}
		res, err := h.db.Exec(sqlstament, idInt)
		if err != nil {
			panic(err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}

		w.Write([]byte(fmt.Sprint("Delete user rows ", count)))
		return
	}
}
