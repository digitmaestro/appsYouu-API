//go:generate prisma1 generate
//go:generate go run gqlgen/cmd.go --config gqlgen/gqlgen.yml
//go:generate go run github.com/google/wire/cmd/wire
//go:generate sed -i "/go:generate wire/d" wire_gen.go

package main

import (
	"log"
	"os"
	"time"

	"github.com/steebchen/keskin-api/prisma"
	"github.com/steebchen/keskin-api/server"
)

func main() {
	server.InitTimezoneOnWindows()
	if tz := time.Now().Location().String(); tz != "UTC" {
		panic("timezone is not UTC but " + tz)
	}

	if len(os.Args) > 1 && os.Args[1] == "trigger-cron-jobs" {
		triggerCronJobs()
	} else {
		server, err := Initialize()
		if err != nil {
			panic(err)
		}

		log.Printf("Server is running on: \nhttp://localhost:%s", server.Config.Port)
		log.Printf("Playground is available at: \nhttp://localhost:%s/api/playground", server.Config.Port)
		err = server.Listen()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func triggerCronJobs() {
	options, err := prisma.NewConfig()
	if err != nil {
		panic(err)
	}
	client, err := prisma.NewClient(options)
	if err != nil {
		panic(err)
	}
	app, err := server.NewFirebaseApp()
	if err != nil {
		panic(err)
	}
	messagingClient, err := server.NewFirebaseMessagingClient(app)
	if err != nil {
		panic(err)
	}

	cron := server.NewCronJobs(client, messagingClient)

	cronEntries := cron.Entries()

	for _, entry := range cronEntries {
		entry.Job.Run()
	}

	time.Sleep(60 * time.Second)
}
