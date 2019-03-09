package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: portforward local remote...")
	}
	localAddrString := os.Args[1]
	remoteAddrString := os.Args[2]
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
		go forward(conn, remoteAddrString)
	}
}

func forward(local net.Conn, remoteAddr string) {
	remote, err := net.Dial("tcp", remoteAddr)
	if remote == nil {
		log.Printf("remote dial failed: %v\n", err)
		local.Close()
		return
	}
	go func() {
		defer local.Close()
		//remote.SetReadTimeout(120*1E9)
		io.Copy(local, remote)
	}()
	go func() {
		defer remote.Close()
		//local.SetReadTimeout(120*1E9)
		io.Copy(remote, local)
	}()
	log.Printf("forward %s to %s", local.RemoteAddr(), remoteAddr)
}
