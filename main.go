package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	port := flag.Int("port", 9000, "The port. In 'server' mode we listen on this port, in client mode we connect to it instead.")
	host := flag.String("host", "localhost", "The hostname or IP address to bind to.")
	server := flag.Bool("server", false, "Run in server mode")

	flag.Parse()

	fmt.Printf("port = %d\n", *port)
	fmt.Printf("host = %s\n", *host)
	fmt.Printf("server = %t\n", *server)

	if *server {
		runServer(*port)
	} else {
		runClient(*host, *port)
	}
}

func runClient(host string, port int) {
	requests := 0
	totalTime := time.Duration(0)
	min := time.Duration(time.Hour)
	max := time.Duration(0)
	remoteAddress := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", remoteAddress)
	if err != nil {
		fmt.Printf("Failed to connect to remote host %s\n", remoteAddress)
		os.Exit(-1)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for i := 1; i <= 100; i++ {
		req := fmt.Sprintf("%d\n", i)
		start := time.Now()
		_, err := writer.Write([]byte(req))
		if err != nil {
			fmt.Printf("Failed to write to remote host %s\n", remoteAddress)
		}
		err = writer.Flush()
		if err != nil {
			fmt.Printf("Failed to write to remote host %s\n", remoteAddress)
		}
		res, err := reader.ReadString('\n')
		end := time.Now()
		requests = requests + 1
		if err != nil {
			fmt.Printf("Error reading response from remote host %s\n", remoteAddress)
			return
		}
		if req != res {
			fmt.Printf("Unexpected response from remote host %s\n", remoteAddress)
		}
		latency := end.Sub(start)
		if latency > max {
			max = latency
		}
		if latency < min {
			min = latency
		}
		totalTime = totalTime + latency
		fmt.Printf("Lat: %.3f ms  Min: %.3f ms  Max: %.3f ms  Avg: %.3f ms\n",
			millis(latency), millis(min), millis(max), millis(totalTime)/float64(requests))
		time.Sleep(time.Second)
	}
}

func millis(t time.Duration) float64 {
	a := float64(t)
	b := float64(time.Millisecond)
	r := a / b
	return r
}

func runServer(port int) {
	localAddress := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", localAddress)
	fmt.Printf("listening on %s\n", localAddress)
	if err != nil {
		panic(err)
	}

	for {
		localConnection, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		// Handle the actual forwarding to the remote
		go runEchoService(localConnection)
	}

}

func runEchoService(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// read client request data
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read data, err:", err)
			}
			return
		}
		line := fmt.Sprintf("%s", bytes)
		fmt.Printf("response: %s", line)
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Printf("failed to write response, err:%v\n", err)
			return
		}
	}
}
