package main

import (
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/net/proxy"
)

func main() {
	if len(os.Args) != 4 {
		log.Println("usage: proxy-forward local proxy remote")
		log.Printf("example: proxy-forward localhost:3080 localhost:1080 192.168.0.1:3080")
		os.Exit(1)
	}
	localAddrString := os.Args[1]
	proxyAddrString := os.Args[2]
	remoteAddrString := os.Args[3]
	localAddr, err := net.ResolveTCPAddr("tcp", localAddrString)
	if localAddr == nil {
		log.Fatalf("net.ResolveTCPAddr failed: %s", err)
	}
	local, err := net.ListenTCP("tcp", localAddr)
	if local == nil {
		log.Fatalf("portforward: %s", err)
	}
	log.Printf("portforward listen on %s", localAddr)
	for {
		conn, err := local.Accept()
		if conn == nil {
			log.Fatalf("accept failed: %s", err)
		}
		go forward(conn, proxyAddrString, remoteAddrString)
	}
}

func forward(local net.Conn, proxyAddrString, remoteAddr string) {
	dialer, err := proxy.SOCKS5("tcp", proxyAddrString, nil, proxy.Direct)
	if err != nil {
		log.Printf("SOCKS5 proxy failed: %v\n", err)
		local.Close()
		return
	}
	remote, err := dialer.Dial("tcp", remoteAddr)
	if remote == nil {
		log.Printf("remote dial failed: %v\n", err)
		local.Close()
		return
	}
	go func() {
		defer local.Close()
		io.Copy(local, remote)
	}()
	go func() {
		defer remote.Close()
		io.Copy(remote, local)
	}()
	log.Printf("forward %s to %s", local.RemoteAddr(), remoteAddr)
}
