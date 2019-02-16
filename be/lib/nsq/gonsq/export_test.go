package gonsq

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"

	nsq "github.com/nsqio/go-nsq"
)

var nullLogger = log.New(ioutil.Discard, "", log.LstdFlags)

func TFuncPatch() {
	newProducer = func(addr string, config *nsq.Config) (*nsq.Producer, error) {

		if strings.Contains(addr, "error") {
			return nil, errors.New("Test error")
		}

		dummyProd, _ := nsq.NewProducer(addr, config)
		dummyProd.SetLogger(nullLogger, nsq.LogLevelInfo)

		return dummyProd, nil
	}
}
