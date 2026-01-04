package client

import "net"

func Conn() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		return nil, err
	}
	return conn, nil
}
