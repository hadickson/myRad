//*********************************************************************
// Copyright ZappyRide
//*********************************************************************
//
// file: myrad/main.go
// Author: Howard Dickson
// Date: 06/22/2020
//
// myrad is an API that allows access to a simple publically accessible
// REST API with one endpoint named: /events. The events endpoint will
// be able to perform the following functions via standard HTTP verbs (GET, POST, etc.):
// read - return a JSON list of all events in the database
// create - add a new event in the database
// update – update an event in the database
// delete – delete an event in the database

// Each event consists of the following fields:
// organizer – string – the name of the organizer organizing the event, e.g. “Plug In America”
// venue – string – the name if the venue, e.g. “New York Auto Show”
// date – string – date of the event, e.g. “June 1, 2020” (we’ll keep this as a string for simplicity)
// An event will have the key of organizer&venu.
package main

import (
	"errors"
	"fmt"

	//"github.com/aws/aws-sdk-go/service/rds"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// The event structure. This is used to map the data retrieved to json format
type selectEvent struct {
	Organizer string `json:"organizer"`
	Venue     string `json:"venue"`
	Date      string `json:"date"`
}

// This function will parameterize the Query string into the appropriate
// internal variables. It will also perform some high level check of
// the query values to guard against SQL Injection (SQLi) type attacks.
func getQuery(queryVal map[string]string, argc int) (string, string, string, error) {

	var argn int
	var organizer, venue, date string

	for k, v := range queryVal {
		argn++

		if argn > argc {
			return organizer, venue, date, errors.New("myrad: bad arg list for method")
		}

		// Check for common web attack formatting moves

		badActor := ";|*|'|/|>|<|\\"
		if strings.ContainsAny(v, badActor) {
			return organizer, venue, date, errors.New("myrad: bad arg list for method. Characters not allowed")
		}

		if k == "organizer" {
			organizer = v
		} else if k == "venue" {
			venue = v
		} else if k == "date" {
			date = v
		}

	}

	return organizer, venue, date, nil
}

// Create the Database Connection

func createConn() (*sql.DB, error) {

	db, err := sql.Open("mysql", "admin:myRad123@tcp(myrad-db.comtrskermlg.us-east-2.rds.amazonaws.com:3306)/myrad_schema")

	if err != nil {
		return db, err
	}

	return db, nil
}

// RequestHandler will manage the request based on the invoked
// method.
func RequestHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var eResp, organizer, venue, date string
	var err error
	var apiResp events.APIGatewayProxyResponse
	var db *sql.DB
	var rows *sql.Rows
	var radEvent selectEvent

	db, err = createConn()

	if err != nil {
		eResp = fmt.Sprintf("{\"DB_CONNECTION_FAILED\":  \"err %s\"}", err)
		apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
		return apiResp, nil
	}

	switch request.HTTPMethod {
	case "GET":
		organizer, venue, date, err = getQuery(request.QueryStringParameters, 0)
		if err != nil {
			eResp = fmt.Sprintf("{\"GET_QRY_FAILURE_NO_ARGS_REQUIRED\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		rows, err = db.Query("select * from myrad_schema.events;")
		if err != nil {
			eResp = fmt.Sprintf("{\"SELECT_ALL_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		jsonStr := "{\n"
		addNL := false

		for rows.Next() {

			if addNL == false {
				addNL = true
			} else {
				jsonStr = jsonStr + ",\n"
			}

			err = rows.Scan(&radEvent.Organizer, &radEvent.Venue, &radEvent.Date)
			if err != nil {
				eResp = fmt.Sprintf("{\"CURSOR_FAILURE\":  \"err %s\"}", err)
				apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
				return apiResp, nil
			}

			jsonByteData, err := json.Marshal(radEvent)
			if err != nil {
				eResp = fmt.Sprintf("{\"JSON_MARSHALL_FAILURE\":  \"err %s\"}", err)
				apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
				return apiResp, nil
			}

			jsonStr = jsonStr + string(jsonByteData)

		}

		jsonStr = jsonStr + "\n}"

		apiResp = events.APIGatewayProxyResponse{Body: jsonStr, StatusCode: 200}
		return apiResp, nil
	case "POST":
		organizer, venue, date, err = getQuery(request.QueryStringParameters, 3)
		if err != nil {
			eResp = fmt.Sprintf("{\"POST_QRY_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		insertStr := fmt.Sprintf("INSERT INTO myrad_schema.events VALUES ( %s, %s, %s )", organizer, venue, date)
		_, err = db.Query(insertStr)
		if err != nil {
			eResp = fmt.Sprintf("{\"INSERT_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		radEvent.Organizer = organizer
		radEvent.Venue = venue
		radEvent.Date = date

		jsonByteData, err := json.Marshal(radEvent)
		if err != nil {
			eResp = fmt.Sprintf("{\"JSON_MARSHALL_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		apiResp = events.APIGatewayProxyResponse{Body: "{\"Status\":  \"CREATE SUCCESSFUL\"}\n" + string(jsonByteData), StatusCode: 200}
		return apiResp, nil
	case "DELETE":
		organizer, venue, date, err = getQuery(request.QueryStringParameters, 2)
		if err != nil {
			eResp = fmt.Sprintf("{\"DELETE_QRY_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		// Check if the row exists before trying to update it.

		selectStr := fmt.Sprintf("select * from myrad_schema.events where organizer=%s AND venue=%s", organizer, venue)
		row := db.QueryRow(selectStr)
		err = row.Scan(&radEvent.Organizer, &radEvent.Venue, &radEvent.Date)
		if err != nil {
			if err == sql.ErrNoRows {
				eResp = fmt.Sprintf("{\"CURSOR_NOROW_FOUND_TO_DELETE\":  \"err %s\"}", err)
				apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
				return apiResp, nil
			}

			eResp = fmt.Sprintf("{\"CURSOR_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		deleteStr := fmt.Sprintf("DELETE from myrad_schema.events where organizer=%s AND venue=%s", organizer, venue)
		_, err = db.Query(deleteStr)
		if err != nil {
			eResp = fmt.Sprintf("{\"DELETE_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		radEvent.Organizer = organizer
		radEvent.Venue = venue

		jsonByteData, err := json.Marshal(radEvent)
		if err != nil {
			eResp = fmt.Sprintf("{\"JSON_MARSHALL_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		apiResp = events.APIGatewayProxyResponse{Body: "{\"Status: \" \"DELETE SUCCESSFUL\"}\n" + string(jsonByteData), StatusCode: 200}
		return apiResp, nil
	case "PUT":
		organizer, venue, date, err = getQuery(request.QueryStringParameters, 3)
		if err != nil {
			eResp = fmt.Sprintf("{\"PUT_QRY_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		// Check if the row exists before trying to update it.

		selectStr := fmt.Sprintf("select * from myrad_schema.events where organizer=%s AND venue=%s", organizer, venue)
		row := db.QueryRow(selectStr)
		err = row.Scan(&radEvent.Organizer, &radEvent.Venue, &radEvent.Date)
		if err != nil {
			if err == sql.ErrNoRows {
				eResp = fmt.Sprintf("{\"CURSOR_NOROW_FOUND_TO_UPDATE\":  \"err %s\"}", err)
				apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
				return apiResp, nil
			}

			eResp = fmt.Sprintf("{\"CURSOR_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		updateStr := fmt.Sprintf("UPDATE  myrad_schema.events set organizer=%s, venue=%s, date=%s  where organizer=%s AND venue=%s", organizer, venue, date, organizer, venue)
		_, err := db.Query(updateStr)
		if err != nil {
			eResp = fmt.Sprintf("{\"UPDATE_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		// Set the date properly for the return JSON
		radEvent.Date = date

		jsonByteData, err := json.Marshal(radEvent)
		if err != nil {
			eResp = fmt.Sprintf("{\"JSON_MARSHALL_FAILURE\":  \"err %s\"}", err)
			apiResp = events.APIGatewayProxyResponse{Body: eResp, StatusCode: 400}
			return apiResp, nil
		}

		apiResp = events.APIGatewayProxyResponse{Body: "{\"Status\":  \"UPDATE SUCCESSFUL\"}\n" + string(jsonByteData), StatusCode: 200}
		return apiResp, nil
	default:
		apiResp = events.APIGatewayProxyResponse{Body: "Bad Method", StatusCode: 400}
		return apiResp, nil
	}

}

func main() {
	lambda.Start(RequestHandler)
}
