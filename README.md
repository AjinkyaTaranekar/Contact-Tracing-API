# Contact Tracing API

## List of Tasks Done
- Completion Percentage :- 75 %
- Total 5 endpoints are working. 
- Make the server thread safe (I didn't know how to proceed with it in existing code, learned there are ways by locking function when one requestion is fetched but couldn't implemented it).
- Added pagination to the list endpoint.
- Add unit tests (added code, but failed to run)

## How to run the application

1. Clone the application.

2. There is a dependency of mongo db driver which need to be imported before running the application. Please get the dependency through the following commands -

    ```shell
        go get "go.mongodb.org/mongo-driver/mongo"
    ```

3. To run the application, please use the following command -

    ```shell
        go run .
    ```

> Note: By default the port number its being run on is **5005**.

> Note: The Mongo DB is setup on my Azure VM.

## Endpoints Description

### Get All Users

```
    URL - *http://localhost:5005/allUsers*
    Method - GET
```

### Get User By ID

```JSON
    URL - *http://localhost:5005/user/<user_id_here>*
    Method - GET
```

### Create User

```JSON
    URL - *http://localhost:5005/user*
    Method - POST
    Body - (content-type = application/json)
    {
    	"name":"John Doe",
    	"emailAddress":"john.doe@gmail.com",
    	"phoneNo":"1234567890",
    	"dateOfBirth":"31-12-2019",
    }
```

### Create Contact

```JSON
    URL - *http://localhost:5005/contacts*
    Method - POST
    Body - (content-type = application/json)
    {
    	"_idOne":"5f83284980ndfs42",
    	"_idTwo":"5f83284980ndfs46"
    }
```

### List all primary contacts within the last 14 days of infection

```JSON
    URL - *http://localhost:5005/contacts?user=<user id>&infection_timestamp=<timestamp>*
    Method - GET
```


## Test Driven Development Description

To run all the unit test cases, please do the following -

1. `go run .`
2. `go test -v`

