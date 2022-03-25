package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

type ServerAuthrozationEntry struct {
	Data        TokenData `json:"data"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
}
type TokenData struct {
	ExpirationTime int    `json:"expirationTime"`
	ServerTime     int    `json:"serverTime"`
	Token          string `json:"token"`
}
type ExifShare struct {
	ExifStatus bool
	mu         sync.Mutex
}

func (exifShare *ExifShare) getExitStatus() bool {
	exifStatus := false
	exifShare.mu.Lock()
	exifStatus = exifShare.ExifStatus
	exifShare.mu.Unlock()
	return exifStatus
}
func (exifShare *ExifShare) setExitStatus(exifStatus bool) {
	exifShare.mu.Lock()
	exifShare.ExifStatus = exifStatus
	exifShare.mu.Unlock()
}

func main() {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Println("hostname:", name)
	exifShare := ExifShare{ExifStatus: false}
	sigint := make(chan os.Signal, 1)
	idleConnsClosed := make(chan struct{})

	go func() {

		signal.Notify(sigint, os.Interrupt)
		<-sigint
		exifShare.setExitStatus(true)
		log.Printf("%v", "shutdown Goroutine: end 1 \n")
		closeConn, err := net.Dial("tcp", ":9000")
		// log.Printf("closeConn: %v", err)
		if err == nil {
			closeConn.Close()
		}
		// We received an interrupt signal, shut down.
		// if err := srv.Shutdown(context.Background()); err != nil {
		// 	// Error from closing listeners, or context timeout:
		// 	log.Printf("server Shutdown: %v", err)
		// }

		log.Printf("%v", "shutdown Goroutine: end 2 \n")
		close(idleConnsClosed)

	}()

	creatUserDataFolder(GetRunningPath() + getUserDataFolder())
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		// handle error
		log.Fatalf("Listen %v", err)
	}
	log.Printf("m1 0\n")
	writelog(getLogFileName(), "m1 0\n")
	for {
		conn, err := ln.Accept()
		if exifShare.getExitStatus() {
			log.Printf("exit\n")
			writelog(getLogFileName(), "exit\n")
			break
		}
		log.Printf("m1 1\n")
		writelog(getLogFileName(), "m1 1\n")
		if err != nil {
			// handle error
			log.Printf("m1 2 %v\n", err)
			writelog(getLogFileName(), fmt.Sprintf("m1 2 %v\n", err))
			continue
		}
		go handleConnection(conn)
	}

	<-idleConnsClosed
	log.Println("end")
}
func handleConnection(c net.Conn) {
	log.Printf("m1 3\n")
	writelog(getLogFileName(), "m1 3\n")
	defer c.Close()
	hasSet := false
	// var Message message
	bytes := make([]byte, 1024)
	for {
		nRead, err := c.Read(bytes)
		if nRead < 0 {
			log.Printf("m1 0 %v\n", err)
			return
		}
		if err != nil {
			log.Printf("m1 4 %v\n", err)
			writelog(getLogFileName(), fmt.Sprintf("m1 4 %v\n", err))
			return
		}
		// if !hasSet {
		// log.Printf("m1 5 \n")
		// Message.readMsg(bytes[0:nRead])
		// log.Printf("Message=%+v\n\n", Message)
		// data := Message.controlMsg(`{"sleep\":\"1\",\"startTime\":52800,\"endTime\":54000,\"type\":1}`)
		// log.Printf("data [%v]%x\n", len(data), data)
		// log.Printf("data [%v]%q\n\n", len(data), data)
		// c.Write(data)
		// 	hasSet = true

		// }

		log.Printf("[%v]%q\n", nRead, bytes[0:nRead])
		writelog(getLogFileName(), fmt.Sprintf("[%v]%q\n", nRead, bytes[0:nRead]))
		// b := []byte{0xaa, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb0, 0xf8, 0x93, 0x11, 0x49, 0x58, 0x00, 0x3d, 0x00, 0x00, 0x02, 0x7b, 0x22, 0x62, 0x72, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x65, 0x73, 0x73, 0x22, 0x3a, 0x22, 0x30, 0x22, 0x2c, 0x22, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x32, 0x7d, 0xff, 0x23, 0x45, 0x4e, 0x44, 0x23}
		// log.Printf("b [%v]%x\n", len(b), b)
		// log.Printf("b [%v]%q\n\n", len(b), b)
		// b2 := []byte{0xaa, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb0, 0xf8, 0x93, 0x11, 0x49, 0x58, 0x00, 0x3d, 0x00, 0x00, 0x02, 0x7b, 0x22, 0x62, 0x72, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x65, 0x73, 0x73, 0x22, 0x3a, 0x22, 0x31, 0x30, 0x30, 0x22, 0x2c, 0x22, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x32, 0x7d, 0xff, 0x23, 0x45, 0x4e, 0x44, 0x23}
		// log.Printf("b2 [%v]%x\n", len(b2), b2)
		// log.Printf("b2 [%v]%q\n\n", len(b2), b2)
		//time is sum of minus
		if !hasSet {
			b3 := []byte{0xaa, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb0, 0xf8, 0x93, 0x11, 0x49, 0x58, 0x00, 0x5a, 0x00, 0x00, 0x02, 0x7b, 0x22, 0x73, 0x6c, 0x65, 0x65, 0x70, 0x22, 0x3a, 0x22, 0x31, 0x22, 0x2c, 0x22, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x3a, '8', '1', '0', '0', '0', 0x2c, 0x22, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x3a, '2', '5', '2', '0', '0', 0x2c, 0x22, 0x74, 0x79, 0x70, 0x65, 0x22, ':', '1', '}', 0xff, '#', 'E', 'N', 'D', '#'}
			log.Printf("b3[%v]%x\n", len(b3), b3)
			log.Printf("b3[%v]%q\n\n", len(b3), b3)
			c.Write(b3)
			hasSet = true

		}
	}

}
func getTodayDate() string {
	t := time.Now().UTC()
	return t.Format("2006-01-02")
}
func getFormatTimeNowUTC() string {
	currentTime := time.Now().UTC()
	return "[ " + currentTime.Format("15:04:05") + " ] "
}
func getLogFileName() string {
	logFileName := getTodayDate()
	return GetRunningPath() + getUserDataFolder() + logFileName + "-m1Log.txt"
}
func GetRunningPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("getCurrent err %v\n", err)
	}
	return dir
}
func writelog(logFileName, logMsg string) {
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer f.Close()

	if _, err = f.WriteString(getFormatTimeNowUTC() + logMsg + "\n"); err != nil {
		log.Fatal(err)

	}

}
func creatUserDataFolder(DataBase string) {
	log.Printf("DataBase %v\n", DataBase)
	if _, err := os.Stat(DataBase); os.IsNotExist(err) {
		err := os.Mkdir(DataBase, 0755)
		if err != nil {
			log.Fatalf("creatUserDataFolder err %v\n", err)
		}
	}

}
func getUserDataFolder() string {

	return "/m1TcpData/"
}
