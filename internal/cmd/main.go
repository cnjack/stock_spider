package main

import (
	"os"
	"stock/internal/apis"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})

	apis.Route(":8080")
}
