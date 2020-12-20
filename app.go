package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"strings"
	"regexp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURL = "mongodb://20.185.230.90:5010/"
var dbName = "users"

// App ...
type App struct {
	DB     *mongo.Database
}

// Run ...
func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, nil))
}

// Initialize ...
func (app *App) Initialize(_user, _password string) {
	fmt.Println("Starting the application....")

	ctx := dbContext(10)
	app.DB, _ = app.configDB(ctx)
	fmt.Println("Connected to MongoDB!")

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
	http.HandleFunc("/allUsers", app.getUsers)
	http.HandleFunc("/users/", app.getUser)
	http.HandleFunc("/users", app.createUser)
	http.HandleFunc("/contacts", app.contact)
	
}

func (app *App) getUsers(writer http.ResponseWriter, req *http.Request) {
	
	if req.Method == "GET" {
		count, _ := strconv.Atoi(req.FormValue("count"))
		start, _ := strconv.Atoi(req.FormValue("start"))

		if count == 0  || count < 1 {
			count = 10
		}
		if start < 0 {
			start = 0
		}
	
		users, err := getUsers(app.DB, start, count)
		if err != nil {
			respondWithError(writer, http.StatusInternalServerError, err.Error())
			return
		}
	
		respondWithJSON(writer, http.StatusOK, users)
	} else {
		respondWithError(writer, http.StatusInternalServerError, req.Method + " Method not allowed, try with GET")
		return
	}	
}

func (app *App) getUser(writer http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		varID := strings.SplitN(req.URL.Path, "/", 3)[2]
		id, _ := primitive.ObjectIDFromHex(varID)

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
	} else {
		respondWithError(writer, http.StatusInternalServerError, req.Method + " Method not allowed, try with GET")
		return
	}
}

func (app *App) createUser(writer http.ResponseWriter, req *http.Request) {
	
	if req.Method == "POST" {
		var user User
		email, _ := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		phone, _ := regexp.Compile("^[0-9]{10}$")
		decoder := json.NewDecoder(req.Body)

		if err := decoder.Decode(&user); err != nil {
			respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
			return
		}

		defer req.Body.Close()

		if email.MatchString(user.EmailAddress) == false {
			respondWithError(writer, http.StatusInternalServerError, "Please check your email")
			return
		}
		
		if phone.MatchString(user.PhoneNo) == false {
			respondWithError(writer, http.StatusInternalServerError, "Please check your phone")
			return
		}
		result, err := user.createUser(app.DB)
		if err != nil {
			respondWithError(writer, http.StatusInternalServerError, err.Error())
			return
		}

		user.ID = result.InsertedID.(primitive.ObjectID)
		respondWithJSON(writer, http.StatusCreated, user)

	} else {
		respondWithError(writer, http.StatusInternalServerError, req.Method + " Method not allowed, try with POST")
		return
	}
}

func (app *App) contact(writer http.ResponseWriter, req *http.Request) {
	var contact Contact
	
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&contact); err != nil {
			respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer req.Body.Close()
		
		objectIDOne, err := primitive.ObjectIDFromHex(contact.ContactIDOne)
		if err != nil{
			fmt.Println(err)
		}
		
		objectIDTwo, err := primitive.ObjectIDFromHex(contact.ContactIDTwo)
		if err != nil{
			fmt.Println(err)
		}
		
		userOne := User{ID: objectIDOne}
		if err := userOne.getUser(app.DB); err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				respondWithError(writer, http.StatusNotFound, "User "+contact.ContactIDOne+" not found" )
			default:
				respondWithError(writer, http.StatusInternalServerError, err.Error())
			}
			return
		}
		
		userTwo := User{ID: objectIDTwo}
		if err := userTwo.getUser(app.DB); err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				respondWithError(writer, http.StatusNotFound, "User "+contact.ContactIDTwo+" not found" )
			default:
				respondWithError(writer, http.StatusInternalServerError, err.Error())
			}
			return
		}
		
		result, err := contact.createContact(app.DB)
		if err != nil {
			respondWithError(writer, http.StatusInternalServerError, err.Error())
			return
		}
		
		contact.ID = result.InsertedID.(primitive.ObjectID)
		respondWithJSON(writer, http.StatusCreated, contact)

	} else if req.Method == "GET" {	

		user := req.URL.Query()["user"][0]
		infectionTimestamp := req.URL.Query()["infection_timestamp"][0]

		count, _ := strconv.Atoi(req.FormValue("count"))
		start, _ := strconv.Atoi(req.FormValue("start"))
		req.Method = "GET"

		if count < 1 {
			count = 10
		}
		if start < 0 {
			start = 0
		}

		users, err := getContactTracing(app.DB, start, count, user, infectionTimestamp)
		if err != nil {
			respondWithError(writer, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(writer, http.StatusOK, users)
	} else {
		respondWithError(writer, http.StatusInternalServerError, req.Method + " Method not allowed, try with POST")
		return
	}	
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