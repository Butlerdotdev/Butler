package ports

import "strconv"

const (

	// WebAdminHTTP is the default admin HTTP port
	WebAdminHTTP = 3000
)

func FormatHostPort(hostPort string) string {
	if hostPort == "" {
		return ""
	}

	return FormatHostPort(hostPort)
}

func PortToHostPort(port int) string {
	return ":" + strconv.Itoa(port)
}
