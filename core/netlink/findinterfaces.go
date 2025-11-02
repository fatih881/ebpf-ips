package netlink

import (
	"log"
	"net"
)

func FindInterfaces() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Can't find interfaces", err)
		return nil, err
	}
	return FilterInterfaces(interfaces)

}

// For now,we are removing down interfaces and loopback interfaces.In future,we may add a function to automatically attach XDP when a down interface is up.
func FilterInterfaces(interfaces []net.Interface) ([]string, error) {
	var availableinterfaces []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		availableinterfaces = append(availableinterfaces, iface.Name)
	}
	log.Printf("info: detected interfaces: %s", availableinterfaces)
	return availableinterfaces, nil
}
