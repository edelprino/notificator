package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Course struct {
	Titolo     string `json:"denominazione"`
	DataInizio string `json:"data_inizio"`
	Comune     string `json:"comune"`
}

type Courses []Course

func (c Courses) stringify() string {
	var buffer bytes.Buffer
	for _, course := range c {
		buffer.WriteString(fmt.Sprintf("%s\n", course.Titolo))
		buffer.WriteString(fmt.Sprintf("DataInizio: %s\n", course.DataInizio))
		buffer.WriteString(fmt.Sprintf("Comune: %s\n", course.Comune))
		buffer.WriteString(fmt.Sprintf("-------\n"))
	}
	return buffer.String()
}

func (c Courses) html() string {
	var buffer bytes.Buffer
	for _, course := range c {
		buffer.WriteString(fmt.Sprintf("<b>%s</b><br>", course.Titolo))
		buffer.WriteString(fmt.Sprintf("DataInizio: %s<br>", course.DataInizio))
		buffer.WriteString(fmt.Sprintf("Comune: %s<br>", course.Comune))
		buffer.WriteString(fmt.Sprintf("<hr>"))
	}
	return buffer.String()
}

func (c Courses) removeCourserAlreadyStarted() Courses {
	var courses Courses
	for _, course := range c {
		when, err := time.Parse("2006-01-02", course.DataInizio)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if time.Now().After(when) {
			continue
		}
		courses = append(courses, course)
	}
	return courses
}

func main() {
	fmt.Println("Starting")
	resp, err := http.Get("https://www.federclimb.it/formazione/index.php?option=com_gare&task=getCorsi&_=1637958764145")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Decode")

	var courses Courses
	err = json.NewDecoder(resp.Body).Decode(&courses)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(courses.removeCourserAlreadyStarted().stringify())

	sendNotification("Corsi FederClimb", courses.removeCourserAlreadyStarted().html())

}

func sendNotification(title string, message string) {
	type Email struct {
		Title   string `json:"value1"`
		Message string `json:"value2"`
	}

	email := Email{title, message}

	json, err := json.Marshal(email)
	if err != nil {
		log.Fatal(err)
	}

	_, err = http.Post("https://maker.ifttt.com/trigger/email/with/key/cLw9xVm1E5_4jobIOI6QIN", "application/json", bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Notification sent")
}
