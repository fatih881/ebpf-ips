package netlink

import (
	"net"

	"go.uber.org/zap"
)

func FindInterfaces(logger *zap.Logger) ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Error("failed to get interfaces", zap.Error(err))
		return nil, err
	}
	return FilterInterfaces(interfaces, logger)

}
func FilterInterfaces(interfaces []net.Interface, logger *zap.Logger) ([]string, error) {
	availableinterfaces := []string{}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		availableinterfaces = append(availableinterfaces, iface.Name)
	}
	logger.Info("available interfaces", zap.Strings("interfaces", availableinterfaces))
	return availableinterfaces, nil
}
