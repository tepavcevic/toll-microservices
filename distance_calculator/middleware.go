package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tepavcevic/toll-microservices/types"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(n CalculatorServicer) *LogMiddleware {
	return &LogMiddleware{
		next: n,
	}
}

func (l *LogMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(started time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(started),
			"err":  err,
			"dist": dist,
		}).Info("calculate distance")
	}(time.Now())

	dist, err = l.next.CalculateDistance(data)
	return
}
