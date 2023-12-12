package main

import (
	msgserver "codeskserver/comsgserver"
	"codeskserver/log"
	"codeskserver/storage"
	types "codeskserver/types"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Create a loop to run until interrupted
	loopRun := true
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		loopRun = false
	}()

	log.Debug("msgserver start...")

	// Create a file with a connection string
	storage.CreateFile("sqlite3cipher://sqlcipher.db?_pragma_key=zhangxin&_pragma_cipher_page_size=4096")

	// AutoMigrate the Client struct
	storage.Instance().AutoMigrate(&types.Client{})

	// Initialize the msgserver
	msgserver.Instance().Init()
	// Start the msgserver
	msgserver.Instance().Start()
	// Run server until interrupted
	for {
		if !loopRun {
			break
		}
		time.Sleep(2 * time.Second)

		// Publish a test message
		data := generateData(60)
		payload, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error encoding data:", err)
			continue
		}
		msgserver.Instance().PublishMsg("testtopic/123", string(payload), 1)
	}

	log.Debug("msgserver try to close...")
	// Stop the msgserver
	msgserver.Instance().Stop()
	log.Debug("msgserver closed...")
}

type DataItem struct {
	ID        int
	IP        string
	Sending   int
	TimeStamp int64
}

func generateData(len int) []DataItem {
	var data []DataItem

	for i := 0; i < len; i++ {
		data = append(data, DataItem{
			ID:        i + 1,
			IP:        fmt.Sprintf("192.168.100.%d", i+2),
			Sending:   rand.Intn(20) + 1,
			TimeStamp: time.Now().UnixNano() / int64(time.Millisecond),
		})
	}

	return data
}
