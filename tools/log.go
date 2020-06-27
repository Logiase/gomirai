package tools

import "github.com/sirupsen/logrus"

var DefaultLogger *logrus.Logger = logrus.New()

func init() {
	DefaultLogger.SetFormatter(&logrus.JSONFormatter{})
}
