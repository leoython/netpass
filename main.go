package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var (
	localPort int
	remoteIP string
	remotePort int
	proto string
)

func main()  {
	flag.IntVar(&localPort, "localPort", 0, "The local port")
	flag.StringVar(&remoteIP, "remoteHost", "", "The remote host")
	flag.IntVar(&remotePort, "remotePort", 0, "The remote port")
	flag.StringVar(&proto, "proto", "tcp", "The remote port")

	flag.Parse()

	if localPort == 0 {
		log.Fatalln("The local port is required")
	}
	if remoteIP == "" {
		log.Fatalln("The remote host is required")
	}
	if remotePort == 0 {
		log.Fatalln("The remote port is required")
	}
	if proto == "tcp" {
		passByTCP(localPort, remoteIP, remotePort)
	} else {
		log.Fatalln("netpass does not support protocols other than TCP at this time")
	}
}

func passByTCP(localPort int, remoteIP string, remotePort int) {
	log.Println("[...] netpass start to listen local port")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", localPort))
	if err != nil {
		log.Fatalln(fmt.Sprintf("failed to listener local port, %s", err.Error()))
	}
	for {
		localConn := getLocalConn(listener)
		if localConn == nil {
			continue
		}
		go func(remote string) {
			remoteConn, err := net.Dial("tcp", remote)
			if err != nil {
				_ = localConn.Close()
				log.Println(fmt.Sprintf("connect to remote server failed: %s", err.Error()))
				return
			}
			var wg sync.WaitGroup
			wg.Add(2)
			log.Println(fmt.Sprintf("[...] netpass start to forward traffic, %s:%d <---> %s:%d", "0.0.0.0", localPort, remoteIP, remotePort))
			go passData(localConn, remoteConn, &wg)
			go passData(remoteConn, localConn, &wg)
			wg.Wait()
			localConn.Close()
			remoteConn.Close()
			log.Println(fmt.Sprintf("[ok] netpass stop to forward traffic, %s:%d <---> %s:%d", "0.0.0.0", localPort, remoteIP, remotePort))
		}(fmt.Sprintf("%s:%d", remoteIP, remotePort))
	}
}

func getLocalConn(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		log.Println(fmt.Sprintf("create local conn failed: %s", err.Error()))
		return nil
	}
	log.Println("[...] netpass get a new connection")
	return conn
}

func passData(conn1 net.Conn, conn2 net.Conn, wg *sync.WaitGroup) {
	io.Copy(conn1, conn2)
	wg.Done()
}