package main

import(
	"fmt"
	log "github.com/sirupsen.logrus"
	"github.com/streadway/amqp"
	"os"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main(){
	submit()
	fmt.Println("Running...")
}

func submit(){

	var stat unix.Statfs_t

	wd, err := os.Getwd()

	unix.Statfs(wd, &stat)

	DiskSpace := stat.Bavail * uint64(stat.Bsize)

	fmt.Println("Available Disk Space: "+ DiskSpace)

	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" +rabbit_port + "/")

	if err := nil {
		log.Fatalf("%s: %s","Failed to connect to RabbitMQ",err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err := nil {
		log.Fatalf("%s: %s","Failed to open a Channel",err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"publisher",
		false,
		false,
		false,
		false,
		nil,
	)

	if err := nil {
		log.Fatalf("%s: %s","Failed to declare a queue",err)
	}

	err := ch.publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:		 []byte(DiskSpace),
		}
	)

	if err := nil {
		log.Fatalf("%s: %s","Failed to publish Disk Space",err)
	}

	fmt.Println("Publis Success!")
}