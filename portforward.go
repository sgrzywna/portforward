package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/net/proxy"
)

func main() {
	flag.Usage = usage

	proxyAddrString := flag.String("proxy", "", "SOCKS5 proxy address")

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	localAddrString := flag.Arg(0)
	remoteAddrString := flag.Arg(1)

	localAddr, err := net.ResolveTCPAddr("tcp", localAddrString)
	if localAddr == nil {
		log.Fatalf("net.ResolveTCPAddr failed: %s", err)
	}

	local, err := net.ListenTCP("tcp", localAddr)
	if local == nil {
		log.Fatalf("portforward: %s", err)
	}

	if *proxyAddrString == "" {
		log.Printf("portforward listen on %s forward to %s", localAddr, remoteAddrString)
	} else {
		log.Printf("portforward listen on %s forward to %s through %s", localAddr, remoteAddrString, *proxyAddrString)
	}

	for {
		conn, err := local.Accept()
		if conn == nil {
			log.Fatalf("accept failed: %s", err)
		}
		if *proxyAddrString == "" {
			go directForward(conn, remoteAddrString)
		} else {
			go proxyForward(conn, *proxyAddrString, remoteAddrString)
		}
	}
}

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] localaddr:port remoteaddr:port\n", os.Args[0])
	flag.PrintDefaults()
}

func directForward(local net.Conn, remoteAddr string) {
	remote, err := net.Dial("tcp", remoteAddr)
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

func proxyForward(local net.Conn, proxyAddrString, remoteAddr string) {
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
