package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
)

var (
	consumerKey       = getEnvSecret("consumerKey")
	consumerSecret    = getEnvSecret("consumerSecret")
	accessToken       = getEnvSecret("accessToken")
	accessTokenSecret = getEnvSecret("accessTokenSecret")
)

func getEnvSecret(env string) string {
	v := os.Getenv(env)
	if v == "" {
		fmt.Fprintf(os.Stderr, "could not find env variable %q", env)
		os.Exit(1)
	}
	return v
}

func main() {
	track := flag.String("track", "", "string to track")
	flag.Parse()
	if *track == "" {
		fmt.Printf("track cannot be empty string\n")
		os.Exit(1)
	}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	api.SetLogger(&logger{logrus.New()})
	f := url.Values{}
	f.Set("track", *track)
	stream := api.PublicStreamFilter(f)
	for c := range stream.C {
		if tweet, ok := c.(anaconda.Tweet); ok {
			fmt.Printf("%s\n", tweet.Text)
		}
	}
}

type logger struct {
	*logrus.Logger
}

func (l *logger) Critical(args ...interface{}) {
	logrus.Error(args)
}

func (l *logger) Criticalf(format string, args ...interface{}) {
	logrus.Errorf(format, args)
}

func (l *logger) Notice(args ...interface{}) {
	logrus.Info(args)
}

func (l *logger) Noticef(format string, args ...interface{}) {
	logrus.Infof(format, args)
}
