package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (c Courses) String() string {
	var buffer bytes.Buffer
	for _, course := range c {
		buffer.WriteString(fmt.Sprintf("<b>%s</b><br>\nDataInizio: %s<br>\nComune: %s<br>\n<hr>\n", course.Titolo, course.DataInizio, course.Comune))
	}
	return buffer.String()
}

func NewCoursesFromReader(reader io.Reader) (Courses, error) {
	var courses Courses
	err := json.NewDecoder(reader).Decode(&courses)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (c Courses) Upcoming() Courses {
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
	fmt.Println("############ Starting")

	notificator := NewArrayNotificator(
		NewConsoleNotificator(),
		NewIFTTTNotificator("cLw9xVm1E5_4jobIOI6QIN"),
	)

	fmt.Println("############ Fetch")
	resp, err := http.Get("https://www.federclimb.it/formazione/index.php?option=com_gare&task=getCorsi&_=1637958764145")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("############ Decode")
	courses, err := NewCoursesFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("############ Notify")
	notificator.Notify("Prossimi corsi FederClimb", courses.Upcoming().String())
}
