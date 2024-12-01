package service

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
)

type BeanstalkClient struct {
	Conn *beanstalk.Conn
}

func NewBeanstalkClient() (*BeanstalkClient, error) {
	conn, err := beanstalk.Dial("tcp", "host.docker.internal:11300")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Beanstalk: %w", err)
	}
	return &BeanstalkClient{Conn: conn}, nil
}