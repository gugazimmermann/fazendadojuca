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
	Name string `json:"name,omitempty"`
}

type breed struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type purityLevel struct {
	ID    int    `json:"id,omitempty"`
	Level string `json:"level,omitempty"`
}

type animal struct {
	ID           int         `json:"id,omitempty"`
	Name         string      `json:"name"`
	Gender       gender      `json:"gender"`
	Breed        breed       `json:"breed"`
	PurityLevel  purityLevel `json:"purity_level"`
	Number       string      `json:"number"`
	Registry     string      `json:"registry"`
	Origin       string      `json:"origin"`
	Father       int         `json:"father"`
	Mother       int         `json:"mother"`
	Insemination int         `json:"insemination"`
	Birth        string      `json:"birth"`
	Death        string      `json:"death"`
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

func serviceFetchOne(id int) (*animal, error) {
	if id == 0 {
		return nil, errors.New("Invalid ID")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	a := new(animal)
	row := db.QueryRow(`
	SELECT
		a.id,
		a.name,
		g.id,
		g.name,
		b.id,
		b.name,
		p.id,
		p.level,
		a.number,
		a.registry,
		a.origin,
		a.father,
		a.mother,
		a.insemination,
		a.birth,
		a.death
	FROM animal a
		JOIN gender g ON g.id = a.gender_id
		JOIN breed b ON b.id  = a.breed_id
		JOIN purity_level p ON p.id = a.purity_level_id
	WHERE a.id= ?`,
		id)
	err = row.Scan(
		&a.ID,
		&a.Name,
		&a.Gender.ID,
		&a.Gender.Name,
		&a.Breed.ID,
		&a.Breed.Name,
		&a.PurityLevel.ID,
		&a.PurityLevel.Level,
		&a.Number,
		&a.Registry,
		&a.Origin,
		&a.Father,
		&a.Mother,
		&a.Insemination,
		&a.Birth,
		&a.Death,
	)
	if err != nil && err != sql.ErrNoRows {
		checkError(err)
	}
	return a, nil
}

func serviceFetchAll() ([]*animal, error) {
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	results, err := db.Query(`
	SELECT
		a.id,
		a.name,
		g.id,
		g.name,
		b.id,
		b.name,
		p.id,
		p.level,
		a.number,
		a.registry,
		a.origin,
		a.father,
		a.mother,
		a.insemination,
		a.birth,
		a.death
	FROM animal a
		JOIN gender g ON g.id = a.gender_id
		JOIN breed b ON b.id  = a.breed_id
		JOIN purity_level p ON p.id = a.purity_level_id`)
	checkError(err)
	as := []*animal{}
	for results.Next() {
		var a = new(animal)
		err = results.Scan(
			&a.ID,
			&a.Name,
			&a.Gender.ID,
			&a.Gender.Name,
			&a.Breed.ID,
			&a.Breed.Name,
			&a.PurityLevel.ID,
			&a.PurityLevel.Level,
			&a.Number,
			&a.Registry,
			&a.Origin,
			&a.Father,
			&a.Mother,
			&a.Insemination,
			&a.Birth,
			&a.Death,
		)
		if err != nil && err != sql.ErrNoRows {
			checkError(err)
		}
		as = append(as, a)
	}
	return as, nil
}

func serviceCreate(req events.APIGatewayProxyRequest) (*animal, error) {
	a := new(animal)
	err := json.Unmarshal([]byte(req.Body), &a)
	if err != nil {
		return nil, errors.New("Invalid Data")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	res, err := db.Exec(`
	INSERT INTO animal (
		name,
		gender_id,
		breed_id,
		purity_level_id,
		number,
		registry,
		origin,
		father,
		mother,
		insemination,
		birth,
		death
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		&a.Name,
		&a.Gender.ID,
		&a.Breed.ID,
		&a.PurityLevel.ID,
		&a.Number,
		&a.Registry,
		&a.Origin,
		&a.Father,
		&a.Mother,
		&a.Insemination,
		&a.Birth,
		&a.Death)
	checkError(err)
	aID, err := res.LastInsertId()
	checkError(err)
	a, err = serviceFetchOne(int(aID))
	return a, nil
}

func serviceUpdate(req events.APIGatewayProxyRequest) (*animal, error) {
	a := new(animal)
	err := json.Unmarshal([]byte(req.Body), &a)
	if err != nil {
		return nil, errors.New("Invalid Data")
	}
	if a.ID == 0 {
		return nil, errors.New("Invalid ID")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	rows, err := db.Exec(`
	UPDATE animal SET 
		name = ?,
		gender_id = ?,
		breed_id = ?,
		purity_level_id = ?,
		number = ?,
		registry = ?,
		origin = ?,
		father = ?,
		mother = ?,
		insemination = ?,
		birth = ?,
		death = ?
	WHERE id = ?;`,
		&a.Name,
		&a.Gender.ID,
		&a.Breed.ID,
		&a.PurityLevel.ID,
		&a.Number,
		&a.Registry,
		&a.Origin,
		&a.Father,
		&a.Mother,
		&a.Insemination,
		&a.Birth,
		&a.Death,
		&a.ID)
	rowCount, err := rows.RowsAffected()
	if err != nil || rowCount == 0 {
		return nil, errors.New("Could Not Update")
	}
	a, err = serviceFetchOne(a.ID)
	return a, nil
}

func serviceDelete(id int) error {
	if id == 0 {
		return errors.New("Invalid ID")
	}
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	rows, err := db.Exec("DELETE FROM animal WHERE id = ?", id)
	rowCount, err := rows.RowsAffected()
	if err != nil || rowCount == 0 {
		return errors.New("Could Not Delete")
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
