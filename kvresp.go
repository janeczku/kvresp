// Package kvresp provides answers to queries from the Lua generic UDP Question/Answer stuff in 
// https://github.com/PowerDNS/pdns/blob/master/pdns/kv-example-script.lua
package main

import (
      "flag"
      "fmt"
      "log"
      "net"
      "strings"
)

var (
      listenAddr string
      verbose bool
)

type message struct{
      command string
      value   string
}

func init() {
      flag.StringVar(&listenAddr, "listen", ":5555",
            "host:port to listen on")
      flag.BoolVar(&verbose, "verbose", false,
            "Be more verbose")
}

func handleConnection(conn *net.UDPConn) {
      buf := make([]byte, 1500)
      resp := "0"

      n, addr, err := conn.ReadFromUDP(buf[0:])
      if err != nil {
            log.Println("Error receiving: ", err)
            return
      }

      p := string(buf[0:n])
      if verbose {
            log.Println("Got packet: ", p)
      }

      msg, err := parsePacket(&p)
      if err != nil {
            log.Println(err)
            resp = "???"
      } else {
            switch msg.command {
            case "DOMAIN":
                  if check := strings.Contains(msg.value, "xxx"); check {
                        resp = "1"
                  }
            case "IP":
                  if check := strings.Contains(msg.value, "127.0.0.1"); check {
                        resp = "1"
                  }
            default:
                  resp = "???"
            }
      }

      _, err = conn.WriteToUDP([]byte(resp), addr)
      if err != nil {
            log.Println("Error sending: ", err)
      }

      if verbose {
            log.Println("Sent reply: ", resp)
      }
}

func parsePacket(p *string) (*message, error) {
      parts := strings.Fields(*p)
      if len(parts) < 2 {
            return nil, fmt.Errorf("Invalid packet: %s", p)
      }
      return &message{parts[0], parts[1]}, nil
}

func main() {
      flag.Parse()

      udpAddr, err := net.ResolveUDPAddr("udp4", listenAddr)
      if err != nil {
            log.Fatal(err)
      }

      l, err := net.ListenUDP("udp", udpAddr)
      if err != nil {
            log.Fatal(err)
      }

      log.Println("UDP server up and listening on ", listenAddr)
      defer l.Close()

      for {
            handleConnection(l)
      }
}