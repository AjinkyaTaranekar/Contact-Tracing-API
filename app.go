package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURL = "mongodb://20.185.230.90:5010/"
var dbName = "users"

// App ...
type App struct {
	Router *mux.Router
	DB     *mongo.Database
}

// Run ...
func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

// Initialize ...
func (app *App) Initialize(_user, _password string) {
	fmt.Println("Starting the application....")

	ctx := dbContext(10)
	app.DB, _ = app.configDB(ctx)
	fmt.Println("Connected to MongoDB!")

	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func (app *App) configDB(ctx context.Context) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("Mongo client couldn't connect with background context: %v", err)
	}
	return client.Database(dbName), nil
}

// routing
func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/users", app.getUsers).Methods("GET")
	app.Router.HandleFunc("/users/{id}", app.getUser).Methods("GET")
	app.Router.HandleFunc("/users", app.createUser).Methods("POST")
	app.Router.HandleFunc("/contacts", app.createContact).Methods("POST")
}

func (app *App) getUsers(writer http.ResponseWriter, req *http.Request) {
	count, _ := strconv.Atoi(req.FormValue("count"))
	start, _ := strconv.Atoi(req.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	bs, err := getUsers(app.DB, start, count)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, bs)
}

func (app *App) getUser(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	user := User{ID: id}
	if err := user.getUser(app.DB); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			respondWithError(writer, http.StatusNotFound, "User not found")
		default:
			respondWithError(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(writer, http.StatusOK, user)
}

func (app *App) createUser(writer http.ResponseWriter, req *http.Request) {
	var user User

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer req.Body.Close()

	result, err := user.createUser(app.DB)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusCreated, result)
}

func (app *App) createContact(writer http.ResponseWriter, req *http.Request) {
	var contact Contact

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&contact); err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer req.Body.Close()

	result, err := contact.createContact(app.DB)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusCreated, result)
}

// helpers
func respondWithError(writer http.ResponseWriter, code int, message string) {
	respondWithJSON(writer, code, map[string]string{"error": message})
}

func respondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(response)
}

func dbContext(i time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), i*time.Second)
	return ctx
}