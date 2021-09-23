package router

import (
	"github.com/gorilla/mux"
	"github.com/prest/prest/config"
	"github.com/prest/prest/controllers"
	"github.com/prest/prest/middlewares"
	"github.com/urfave/negroni"
)

var router *mux.Router

func initRouter() {
	router = mux.NewRouter().StrictSlash(true)
}

// GetRouter reagister all routes
func GetRouter() *mux.Router {
	if router == nil {
		initRouter()
	}

	// if auth is enabled
	if config.PrestConf.AuthEnabled {
		router.HandleFunc("/auth", controllers.Auth).Methods("POST")
	}
	router.HandleFunc("/p/databases", controllers.GetDatabases).Methods("GET")
	router.HandleFunc("/p/schemas", controllers.GetSchemas).Methods("GET")
	router.HandleFunc("/p/tables", controllers.GetTables).Methods("GET")
	router.HandleFunc("/p/_QUERIES/{queriesLocation}/{script}", controllers.ExecuteFromScripts)
	router.HandleFunc("/p/{database}/{schema}", controllers.GetTablesByDatabaseAndSchema).Methods("GET")
	router.HandleFunc("/p/show/{database}/{schema}/{table}", controllers.ShowTable).Methods("GET")
	crudRoutes := mux.NewRouter().PathPrefix("/").Subrouter().StrictSlash(true)
	crudRoutes.HandleFunc("/p/{database}/{schema}/{table}", controllers.SelectFromTables).Methods("GET")
	crudRoutes.HandleFunc("/p/{database}/{schema}/{table}", controllers.InsertInTables).Methods("POST")
	crudRoutes.HandleFunc("/p/batch/{database}/{schema}/{table}", controllers.BatchInsertInTables).Methods("POST")
	crudRoutes.HandleFunc("/p/{database}/{schema}/{table}", controllers.DeleteFromTable).Methods("DELETE")
	crudRoutes.HandleFunc("/p/{database}/{schema}/{table}", controllers.UpdateTable).Methods("PUT", "PATCH")
	router.PathPrefix("/").Handler(negroni.New(
		middlewares.AccessControl(),
		middlewares.AuthMiddleware(),
		negroni.Wrap(crudRoutes),
	))

	return router
}

// Routes for pREST
func Routes() *negroni.Negroni {
	n := middlewares.GetApp()
	n.UseHandler(GetRouter())
	return n
}
