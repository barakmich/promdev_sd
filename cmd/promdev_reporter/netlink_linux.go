package main

import (
	"log"
	"net"

	"github.com/jsimonetti/rtnetlink/rtnl"
)

func getRouteToIP(ip net.IP) *net.Interface {
	conn, err := rtnl.Dial(nil)
	if err != nil {
		log.Fatal("can't establish netlink connection: ", err)
	}
	defer conn.Close()
	route, err := conn.RouteGet(ip)
	if err != nil {
		log.Fatalln("can't get route to ip", ip, ":", err)
	}
	return route.Interface
}
