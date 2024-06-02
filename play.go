package main

import (
	"time"
)

func main() {
	MSKLoc, _ := time.LoadLocation("Europe/Moscow")

	println(time.Now().String())
	println(time.Now().Local().String())
	println(time.Now().In(MSKLoc).String())
	println(time.Now().In(time.UTC).String())
}
