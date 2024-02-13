//go:build !linux
// +build !linux

package main

func getRouteToIP(ip net.IP) *net.Interface {
	r, err := netroute.New()
	if err != nil {
		log.Fatalln("Couldn't get routing table:", err)
	}
	outif, _, _, err := r.Route(ip)
	if err != nil {
		log.Fatalln("Couldn't get route to", ip, ":", err)
	}
	return outif
}
