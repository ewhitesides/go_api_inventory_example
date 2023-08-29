package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT * FROM products WHERE id=%v", p.ID)
	row := db.QueryRow(query)
	err := row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errors.New("product not found")
		default:
			return err
		}
	}
	return nil
}

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf(
		"INSERT INTO products(name, quantity, price) VALUES('%s', %v, %v)",
		p.Name,
		p.Quantity,
		p.Price,
	)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil
	}
	p.ID = int(id)
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf(
		"UPDATE products SET name='%s', quantity=%v, price=%v WHERE id=%v",
		p.Name,
		p.Quantity,
		p.Price,
		p.ID,
	)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM products WHERE id=%v", p.ID)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
