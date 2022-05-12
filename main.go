package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gookit/color"
)

type Config struct {
	Url      string `json:"url"`
	School   string `json:"school"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//go:embed "config.json"
var f []byte

func main() {
	config := Config{}
	err := json.Unmarshal(f, &config)
	if err != nil {
		panic(err)
	}

	//get client
	webuntis := NewWebUntis(config.Url, config.School)
	_, err = webuntis.Authenticate(config.Username, config.Password)
	if err != nil {
		panic(err)
	}
	defer webuntis.Logout()

	//get day from arguments
	var res map[string]any
	args := os.Args[1:]
	if len(args) > 0 {
		i, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Input is not a number")
			os.Exit(0)
		}
		res, err = webuntis.GetTimetableForStudent(i)
		if err != nil {
			panic(err)
		}
	} else {
		res, err = webuntis.GetTimetableForStudent(0)
	}
	if err != nil {
		panic(err)
	}

	if res["result"] == nil {
		fmt.Println("Got no data")
		os.Exit(0)
	}

	timetable := res["result"].([]any)

	//no data, abort
	if len(timetable) == 0 {
		fmt.Println("Got no data")
		os.Exit(0)
	}

	sort.Slice(timetable, func(p, q int) bool {
		return timetable[p].(map[string]any)["startTime"].(float64) < timetable[q].(map[string]any)["startTime"].(float64)
	})
	//remove broken entries
	//might break stuff but too lazy to fix
	for i, e := range timetable {
		if len(e.(map[string]any)["su"].([]any)) == 0 || len(e.(map[string]any)["ro"].([]any)) == 0 {
			timetable = removeFromSlice(timetable, i)
		}
	}
	//merge if next is same
	for i, e := range timetable {
		if i < len(timetable)-1 {
			if e.(map[string]any)["su"].([]any)[0].(map[string]any)["id"].(float64) == timetable[i+1].(map[string]any)["su"].([]any)[0].(map[string]any)["id"].(float64) {
				e.(map[string]any)["endTime"] = timetable[i+1].(map[string]any)["endTime"]
				timetable = removeFromSlice(timetable, i+1)
			}
		}
	}
	//get subjects
	res, err = webuntis.GetSubjects()
	if err != nil {
		panic(err)
	}
	subjects := res["result"].([]any)
	//get rooms
	res, err = webuntis.GetRooms()
	if err != nil {
		panic(err)
	}
	rooms := res["result"].([]any)
	//print stuff

	//date
	dateStr := fmt.Sprint(int(timetable[0].(map[string]any)["date"].(float64)))
	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Format("02.01.2006"))
	fmt.Println()

	for _, e := range timetable {
		var name string
		var room string
		var time string

		//getSubjects
		id := e.(map[string]any)["su"].([]any)[0].(map[string]any)["id"].(float64)
		for _, s := range subjects {
			if s.(map[string]any)["id"].(float64) == id {
				name = s.(map[string]any)["name"].(string)
			}
		}
		//get room
		id = e.(map[string]any)["ro"].([]any)[0].(map[string]any)["id"].(float64)
		for _, s := range rooms {
			if s.(map[string]any)["id"].(float64) == id {
				room = s.(map[string]any)["name"].(string)
			}
		}
		//get time
		times := []string{
			fmt.Sprint(e.(map[string]any)["startTime"].(float64)),
			fmt.Sprint(e.(map[string]any)["endTime"].(float64)),
		}
		for i, t := range times {
			times[i] = t[:len(t)-2] + ":" + t[len(t)-2:]
		}

		time = fmt.Sprintf("%s - %s", times[0], times[1])

		//check cancelled
		if _, ok := e.(map[string]any)["code"]; ok {
			if e.(map[string]any)["code"].(string) == "cancelled" {
				name = fmt.Sprintf("%s%s - cancelled%s", "<red>", name, "</>")
			}
		}

		color.Println(name)
		color.Gray.Println(room)
		color.Gray.Println(time)
		fmt.Println()
	}
}

func removeFromSlice(slice []any, s int) []any {
	return append(slice[:s], slice[s+1:]...)
}
