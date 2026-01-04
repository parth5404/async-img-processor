package client

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"tcp/test/constants"
)

func ReadAndSend() {
	conn, err := Conn()
	if err != nil {
		log.Fatal("connection error:", err)
	}
	defer conn.Close()

	entries, err := os.ReadDir(constants.Src)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".jpg" && ext != ".png" {
			continue
		}

		path := filepath.Join(constants.Src, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		size := len(data)
		// if size > constants.MAX_FILE_SIZE {
		// 	log.Println("skip large file:", e.Name())
		// 	continue
		// }

		if err := Send(conn, size, data); err != nil {
			log.Println("send error:", err)
			break
		}
	}
}
