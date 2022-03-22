package tool

import "net"

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

