package graphite

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	tcpClient net.Conn
	addr string
}

func NewGraphiteClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		tcpClient: conn,
		addr:      addr,
	}, nil
}

func (gc *Client) AddMetric(path, value string) error {
	_, err := gc.tcpClient.Write([]byte(fmt.Sprintf("%s %s %d", path, value, time.Now().Unix())))
	if err != nil {
		return err
	}

	return nil
}