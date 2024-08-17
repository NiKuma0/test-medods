package main

import "time"

func sendMsg(msg string) { println(msg) }

func main() {
	defer sendMsg("msg")
	defer sendMsg("msg1")

	<-time.After(10 * time.Millisecond)
}
