package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tepavcevic/toll-microservices/types"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(n Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: n,
	}
}

func (m *LogMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(started time.Time) {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"took": time.Since(started),
		}).Info("aggregating distance")
	}(time.Now())

	err = m.next.AggregateDistance(dist)
	return
}
