package openvpn_status

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseStatusLog(filePath string) ([]clientInfo, error) {
	conn, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	if scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		if words[0] == "OpenVPN" && words[1] == "CLIENT" && words[2] == "LIST" {
			return parseStatusLogV1(scanner), nil
		} else if words[0] == "TITLE,OpenVPN" {
			return parseStatusLogV2(scanner), nil
		} else if words[0] == "TITLE" && words[1] == "OpenVPN" {
			return parseStatusLogV3(scanner), nil
		} else {
			return nil,
				fmt.Errorf("the status log file is invalid")
		}

	}

	return nil,
		fmt.Errorf("the status log file is invalid")
}

func parseStatusLogV1(scanner *bufio.Scanner) []clientInfo {
	checkClientListHeader := func(headers []string) bool {
		var clientListHeaderColumns = [5]string{
			"Common Name",
			"Real Address",
			"Bytes Received",
			"Bytes Sent",
			"Connected Since",
		}
		for i, v := range headers {
			if v != clientListHeaderColumns[i] {
				return false
			}
		}
		return true

	}
	clientListHeader := false

	var clients []clientInfo

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		const dateLayout = "Mon Jan 2 15:04:05 2006"
		if checkClientListHeader(parts) {
			clientListHeader = true
			continue
		}
		if clientListHeader && len(parts) == 5 {
			in, _ := strconv.Atoi(parts[2])
			out, _ := strconv.Atoi(parts[3])
			connectedSince, _ := time.Parse(dateLayout, parts[4])

			c := clientInfo{
				CommonName:     parts[0],
				BytesReceived:  in,
				BytesSent:      out,
				ConnectedSince: connectedSince.Unix(),
			}
			clients = append(clients, c)
		} else {
			clientListHeader = false
		}
	}
	return clients
}

func parseStatusLogV2(scanner *bufio.Scanner) []clientInfo {
	var clients []clientInfo

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		switch parts[0] {
		case "HEADER":
		case "END":
			break
		default:
			switch statusType := parts[0]; statusType {
			case "CLIENT_LIST":
				in, _ := strconv.Atoi(parts[5])
				out, _ := strconv.Atoi(parts[6])
				connectedSinceUnix, _ := strconv.Atoi(parts[8])

				c := clientInfo{
					CommonName:     parts[1],
					BytesReceived:  in,
					BytesSent:      out,
					ConnectedSince: int64(connectedSinceUnix),
				}
				clients = append(clients, c)
			}
		}
	}
	return clients
}

func parseStatusLogV3(scanner *bufio.Scanner) []clientInfo {
	var clients []clientInfo

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		switch parts[0] {
		case "HEADER":
		case "END":
			break
		default:
			switch statusType := parts[0]; statusType {
			case "CLIENT_LIST":

				// v3 use only space for missing field of ipv6.
				// This makes it error-prone to parse status v3.
				// As a workaround,lets first check if ipv6 is
				// available in the line to deduce the right index
				// for other relevant fields
				var fieldIndexIPv6 int
				if r := net.ParseIP(parts[4]); r == nil {
					fieldIndexIPv6 = 3
				} else {
					fieldIndexIPv6 = 4
				}

				in, _ := strconv.Atoi(parts[fieldIndexIPv6+1])
				out, _ := strconv.Atoi(parts[fieldIndexIPv6+2])
				connectedSinceUnix, _ := strconv.Atoi(parts[fieldIndexIPv6+4])

				c := clientInfo{
					CommonName:     parts[1],
					BytesReceived:  in,
					BytesSent:      out,
					ConnectedSince: int64(connectedSinceUnix),
				}
				clients = append(clients, c)
			}
		}
	}
	return clients
}
