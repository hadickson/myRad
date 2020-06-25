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
read (GET) - return a JSON list of all events in the database
create (POST) - add a new event in the database
update (PUT) – update an event in the database
delete (DELETE) – delete an event in the database

// Each event consists of the following fields:
// organizer – string – the name of the organizer organizing the event, e.g. “Plug In America”
// venue – string – the name if the venue, e.g. “New York Auto Show”
// date – string – date of the event, e.g. “June 1, 2020” (we’ll keep this as a string for simplicity)
// An event will have the key of organizer&venu.

// The AWS RDS database is MySQL. 
// The schema for the database is as follows
Schema Name is myrad_schema
Table Name is events
Fields are:
  organizer varchar(30) PK NN
  venue     varchar(45) PK NN
  date      varchar(25)

// A few examples of URLS are:

GET https://5v1qow85u9.execute-api.us-east-2.amazonaws.com/RND/events
POST https://5v1qow85u9.execute-api.us-east-2.amazonaws.com/RND/events?venue="Tesla Auto Show"&date="July 10, 2021"&organizer="ZappyRide"
PUT https://5v1qow85u9.execute-api.us-east-2.amazonaws.com/RND/events?venue="NYC Auto Show"&date="July 12, 2021"&organizer="Howard Org 2"
DEL https://5v1qow85u9.execute-api.us-east-2.amazonaws.com/RND/events?date="July 10, 2021"&organizer="ZappyRide"

// GET OUTPUT EXAMPLE IS AS FOLLOWS:
{
{"organizer":"BOB","venue":"NYC Auto Show","date":"August 22, 2021"},
{"organizer":"BOB","venue":"NYC Auto Show 2","date":"August 22, 2021"},
{"organizer":"BOB","venue":"NYC Auto Show 3","date":"August 22, 2021"},
{"organizer":"BOB1","venue":"NYC Auto Show 3","date":"July 10 2021"},
{"organizer":"Howard Org","venue":"NYC Auto Show","date":"July 10, 2021"},
{"organizer":"Howard Org 2","venue":"NYC Auto Show","date":"July 12, 2021"},
{"organizer":"Howard Org 3","venue":"NYC Auto Show","date":"6/01/2021"}
}

// POST OUTPUT EXAMPLE IS AS FOLLOWS:
{"Status":  "CREATE SUCCESSFUL"}
{"organizer":"\"ZappyRide\"","venue":"\"Tesla Auto Show\"","date":"\"July 10, 2021\""}

// PUT OUTPUT EXAQMPLE IS AS FOLLOWS:
{"Status":  "UPDATE SUCCESSFUL"}
{"organizer":"\"Howard Org 2\"","venue":"\"NYC Auto Show\"","date":"\"July 12, 2021\""}

//DELETE OUTPUT EXAMPLE IS AS FOLLOWS:
{"Status: " "DELETE SUCCESSFUL"}
{"organizer":"\"ZappyRide\"","venue":"\"Tesla Auto Show\"","date":""}