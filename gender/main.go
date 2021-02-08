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

type gender struct {
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

func serviceFetchOne(id int) (*gender, error) {
	if id == 0 {
		return nil, errors.New("Invalid ID")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	g := new(gender)
	row := db.QueryRow("SELECT id, name FROM gender WHERE id= ?", id)
	err = row.Scan(&g.ID, &g.Name)
	if err != nil && err != sql.ErrNoRows {
		checkError(err)
	}
	return g, nil
}

func serviceFetchAll() ([]*gender, error) {
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	results, err := db.Query("SELECT id, name FROM gender")
	checkError(err)
	gs := []*gender{}
	for results.Next() {
		var g = new(gender)
		err = results.Scan(&g.ID, &g.Name)
		if err != nil && err != sql.ErrNoRows {
			checkError(err)
		}
		gs = append(gs, g)
	}
	return gs, nil
}

func serviceCreate(req events.APIGatewayProxyRequest) (*gender, error) {
	g := new(gender)
	err := json.Unmarshal([]byte(req.Body), &g)
	if err != nil {
		return nil, errors.New("Invalid Data")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	res, err := db.Exec("INSERT INTO gender (name) VALUES (?);", g.Name)
	checkError(err)
	gID, err := res.LastInsertId()
	checkError(err)
	g, err = serviceFetchOne(int(gID))
	return g, nil
}

func serviceUpdate(req events.APIGatewayProxyRequest) (*gender, error) {
	g := new(gender)
	err := json.Unmarshal([]byte(req.Body), &g)
	if err != nil {
		return nil, errors.New("Invalid Data")
	}
	if g.ID == 0 {
		return nil, errors.New("Invalid ID")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	rows, err := db.Exec("UPDATE gender SET name = ? WHERE id = ?;", g.Name, g.ID)
	rowCount, err := rows.RowsAffected()
	if err != nil || rowCount == 0 {
		return nil, errors.New("Could Not Update")
	}
	g, err = serviceFetchOne(g.ID)
	return g, nil
}

func serviceDelete(id int) error {
	if id == 0 {
		return errors.New("Invalid ID")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	rows, err := db.Exec("DELETE FROM gender WHERE id = ?", id)
	rowCount, err := rows.RowsAffected()
	if err != nil || rowCount == 0 {
		return errors.New("Could Not Delete")
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
