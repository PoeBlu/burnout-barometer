// Copyright 2019 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

// Package pkg contains types and methods for interacting with the barometer.
package pkg

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
)

// Request defines the common form parameters when a slash command is invoked.
type Request struct {
	Text      string   // The submitted text in the slash command
	UserID    string   // Slack User ID that submitted the request
	Timestamp string   // Timestamp of the request
	Area      string   // IANA-compliant area
	DB        Database // Database to insert into

	// The parsed log-message to be passed into the database
	item logItem
}

// Process parses the request and stores to the Database. It forms a series of
// methods that first parses the Text into a database-compatible format, then
// converts the Timestamp based on the Area. Lastly, it then inserts the
// generated log-item into the specified database (BigQuery, Postgres, etc.)
func (r *Request) Process() (*Message, error) {
	m, notes, err := r.parseMessage()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Request.parseMessage")
		return nil, err
	}

	ts, err := r.getTimestamp()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Request.getTimestamp")
		return nil, err
	}

	measure, err := strconv.Atoi(m)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("strconv")
		return nil, err
	}

	r.item = logItem{
		Timestamp:  ts,
		UserID:     r.UserID,
		LogMeasure: measure,
		Notes:      notes,
	}

	if err := r.insertToTable(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Request.insertToTable")
		return nil, err
	}

	return r.item.formatReply()
}

// parseMessage extracts the barometer measure and notes from the form text.
func (r *Request) parseMessage() (string, string, error) {
	list := strings.Fields(r.Text)
	measure := list[0]
	notes := strings.Join(list[1:], " ")
	return measure, notes, nil
}

// getTimestamp obtains the timestamp value from the request.
func (r *Request) getTimestamp() (time.Time, error) {
	i, err := strconv.ParseInt(r.Timestamp, 10, 64)
	if err != nil {
		log.Errorf("cannot parse timestamp %s: %v", r.Timestamp, err)
		return time.Time{}, err
	}
	loc, err := time.LoadLocation(r.Area)
	if err != nil {
		log.Errorf("cannot find location: %s", r.Area)
		return time.Time{}, err
	}

	return time.Unix(i, 0).In(loc), nil
}

// insertToTable adds the item entry into the specified database.
func (r *Request) insertToTable() error {
	if err := r.DB.Insert(r.item); err != nil {
		log.Errorf("error in inserting item: %v", err)
		return err
	}

	return nil
}

// logItem is the user log for the barometer. This also serves as
// the schema for the database.
type logItem struct {
	Timestamp  time.Time
	UserID     string
	LogMeasure int
	Notes      string
}

// Save allows us to implement BigQuery's ValueSaver interface.
func (i *logItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"timestamp":   i.Timestamp,
		"user_id":     i.UserID,
		"log_measure": i.LogMeasure,
		"notes":       i.Notes,
	}, "", nil
}

// formatReply prepares the Slack message as a response to a slash command.
func (i *logItem) formatReply() (*Message, error) {
	attach := Attachment{
		Color: "#ef4631",
		Title: "Burnout Barometer",
		Text:  fmt.Sprintf("Acknowledged"),
	}

	msg := &Message{
		ResponseType: "ephemeral",
		Text:         fmt.Sprintf("Received: %d (%s)", i.LogMeasure, i.Notes),
		Attachments:  []Attachment{attach},
	}

	return msg, nil
}

// Message is the Slack message event. see
// https://api.slack.com/docs/message-formatting for more information.
type Message struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

// Attachment defines the message output after running the slash command.
type Attachment struct {
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
}
