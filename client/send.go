package client

import (
	"encoding/binary"
	"log"
	"net"
)

func Send(conn net.Conn, size int, data []byte) error {
	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(size))
	if _, err := conn.Write(hdr[:]); err != nil {
		log.Println("write size error:", err)
		return err
	}

	if _, err := conn.Write(data); err != nil {
		log.Println("write data error:", err)
		return err
	}
	return nil
}
