package main

import (
	"fmt"
	"os"

	"github.com/eliasbokreta/deprek8/cmd"
	"github.com/eliasbokreta/deprek8/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if err := utils.LoadConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
