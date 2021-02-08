package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/go-sql-driver/mysql"
)

var (
	host     = os.Getenv("DB_ENDPOINT")
	port     = os.Getenv("DB_PORT")
	database = os.Getenv("DB_NAME")
	user     = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
)

var connectionString = fmt.Sprintf(
	"%s:%s@tcp(%s:%s)/%s?allowNativePasswords=true", user, password, host, port, database,
)

type errorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

type breed struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return get(req)
	case "POST":
		return create(req)
	case "PUT":
		return update(req)
	case "DELETE":
		return delete(req)
	default:
		return unhandledMethod()
	}
}

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}

func get(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	queryid := req.QueryStringParameters["id"]
	id, err := strconv.Atoi(queryid)
	if err == nil {
		result, err := serviceFetchOne(id)
		if err != nil {
			return apiResponse(http.StatusBadRequest, errorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}
	result, err := serviceFetchAll()
	if err != nil {
		return apiResponse(http.StatusBadRequest, errorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func create(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	result, err := serviceCreate(req)
	if err != nil {
		return apiResponse(http.StatusBadRequest, errorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusCreated, result)
}

func update(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	result, err := serviceUpdate(req)
	if err != nil {
		return apiResponse(http.StatusBadRequest, errorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func delete(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	queryid := req.QueryStringParameters["id"]
	id, err := strconv.Atoi(queryid)
	err = serviceDelete(id)
	if err != nil {
		return apiResponse(http.StatusBadRequest, errorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, nil)
}

func unhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, "method Not allowed")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func serviceFetchOne(id int) (*breed, error) {
	if id == 0 {
		return nil, errors.New("Invalid ID")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	b := new(breed)
	row := db.QueryRow("SELECT id, name FROM breed WHERE id= ?", id)
	err = row.Scan(&b.ID, &b.Name)
	if err != nil && err != sql.ErrNoRows {
		checkError(err)
	}
	return b, nil
}

func serviceFetchAll() ([]*breed, error) {
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	results, err := db.Query("SELECT id, name FROM breed")
	checkError(err)
	bs := []*breed{}
	for results.Next() {
		var b = new(breed)
		err = results.Scan(&b.ID, &b.Name)
		if err != nil && err != sql.ErrNoRows {
			checkError(err)
		}
		bs = append(bs, b)
	}
	return bs, nil
}

func serviceCreate(req events.APIGatewayProxyRequest) (*breed, error) {
	b := new(breed)
	err := json.Unmarshal([]byte(req.Body), &b)
	if err != nil {
		return nil, errors.New("Invalid Data")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	res, err := db.Exec("INSERT INTO breed (name) VALUES (?);", b.Name)
	checkError(err)
	bID, err := res.LastInsertId()
	checkError(err)
	b, err = serviceFetchOne(int(bID))
	return b, nil
}

func serviceUpdate(req events.APIGatewayProxyRequest) (*breed, error) {
	b := new(breed)
	err := json.Unmarshal([]byte(req.Body), &b)
	if err != nil {
		return nil, errors.New("Invalid Data")
	}
	if b.ID == 0 {
		return nil, errors.New("Invalid ID")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	rows, err := db.Exec("UPDATE breed SET name = ? WHERE id = ?;", b.Name, b.ID)
	rowCount, err := rows.RowsAffected()
	if err != nil || rowCount == 0 {
		return nil, errors.New("Could Not Update")
	}
	b, err = serviceFetchOne(b.ID)
	return b, nil
}

func serviceDelete(id int) error {
	if id == 0 {
		return errors.New("Invalid ID")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	rows, err := db.Exec("DELETE FROM breed WHERE id = ?", id)
	rowCount, err := rows.RowsAffected()
	if err != nil || rowCount == 0 {
		return errors.New("Could Not Delete")
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
