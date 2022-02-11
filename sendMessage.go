package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var UserName string

func initUser() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Insert Username: ")
	UserName, _ = reader.ReadString('\n')
	UserName = strings.TrimRight(UserName, "\r\n")
	fmt.Println("===Tahnks enjoy your chating ===")
}
func main() {
	initUser()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5555/")
	failOnError(err, "Faiiled to connect to RabitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open achanel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"golang-queue", //name
		false,          //durable
		false,          //delete when unused
		false,          //exclusive
		false,          //no-wait
		nil,            //arguments
	)
	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("[message]=> ")
		mPayload, _ := reader.ReadString('\n')
		mPayload = fmt.Sprintf("%s: %s", UserName, mPayload)

		err = ch.Publish("", //exchange
			q.Name, //routing keyha
			false,  //mandatory
			false,  //immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(mPayload),
			})
		failOnError(err, "failed to publish a message")

	}

}
