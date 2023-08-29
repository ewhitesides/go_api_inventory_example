package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize(DbUser string, DbPass string, DbHost string, DbName string) error {
	// Connect to database
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s", DbUser, DbPass, DbHost, DbName)

	var err error
	app.DB, err = sql.Open("mysql", connStr)
	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(
		http.ListenAndServe(address, app.Router),
	)
}

func (app *App) sendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func (app *App) sendError(w http.ResponseWriter, statusCode int, err error) {
	errMsg := map[string]string{"error": err.Error()}
	app.sendResponse(w, statusCode, errMsg)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, err)
		return
	}
	app.sendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idval, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.sendError(w, http.StatusBadRequest, err)
		return
	}

	p := product{ID: idval}
	err = p.getProduct(app.DB)
	if err != nil {
		app.sendError(w, http.StatusNotFound, err)
		return
	}

	app.sendResponse(w, http.StatusOK, p)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {

	var p product

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		app.sendError(w, http.StatusBadRequest, err)
		return
	}

	err = p.createProduct(app.DB)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, err)
		return
	}

	app.sendResponse(w, http.StatusCreated, p)
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idval, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.sendError(w, http.StatusBadRequest, err)
		return
	}

	var p product
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		app.sendError(w, http.StatusBadRequest, err)
		return
	}

	p.ID = idval
	err = p.updateProduct(app.DB)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, err)
		return
	}

	app.sendResponse(w, http.StatusOK, p)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idval, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.sendError(w, http.StatusBadRequest, err)
		return
	}

	p := product{ID: idval}
	err = p.deleteProduct(app.DB)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, err)
		return
	}

	result := map[string]string{
		"result": "if resource existed, it has been deleted",
	}
	app.sendResponse(w, http.StatusOK, result)
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}
