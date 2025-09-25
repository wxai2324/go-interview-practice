package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

// Product represents a product in the inventory system
type Product struct {
	ID       int64
	Name     string
	Price    float64
	Quantity int
	Category string
}

// ProductStore manages product operations
type ProductStore struct {
	db *sql.DB
}

// NewProductStore creates a new ProductStore with the given database connection
func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{db: db}
}

// InitDB sets up a new SQLite database and creates the products table
func InitDB(dbPath string) (*sql.DB, error) {
	// The table should have columns: id, name, price, quantity, category

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL ,
    price REAL NOT NULL ,
    quantity integer not null default 0,
    category TEXT
)
`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	res, err := ps.db.Exec("insert into products(name,price,quantity,category) values (?,?,?,?);", product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = id
	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {

	row := ps.db.QueryRow("select * from products where id=?", id)
	p := &Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, err
	}
	return p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	res, err := ps.db.Exec("update products set name=?,price=?,quantity=?,category=? where id=?", product.Name, product.Price, product.Quantity, product.Category, product.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product with ID %d not found", product.ID)
	}
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	res, err := ps.db.Exec("delete from products where id = ? ", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	var (
		result []*Product
		rows   *sql.Rows
		err    error
	)
	if category != "" {
		rows, err = ps.db.Query("select id,name,price,quantity,category from products where category = ?", category)
	} else {
		rows, err = ps.db.Query("select id,name,price,quantity,category from products")
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := &Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
	// TODO: For each product ID in the updates map, update its quantity
	// TODO: If any update fails, roll back the transaction
	// TODO: Otherwise, commit the transaction
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
	}()
	stmt, err := tx.Prepare("update products set quantity = ? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id, quantity := range updates {
		result, err := stmt.Exec(quantity, id)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return fmt.Errorf("product with ID %d not found", id)
		}
	}

	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
	db, err := InitDB("products.db")
	if err != nil {
		panic(err)
	}
	product := NewProductStore(db)
	product.CreateProduct(&Product{Name: "Apple", Price: 1.5, Quantity: 10, Category: "Fruit"})
}
