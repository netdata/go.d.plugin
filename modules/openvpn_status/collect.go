package openvpn_status

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type clientInfo struct {
	CommonName     string
	BytesReceived  int
	BytesSent      int
	ConnectedSince int64
}

func (o *OpenVPNStatus) collect() (map[string]int64, error) {
	var err error

	mx := make(map[string]int64)

	clients, err := parseStatusLog(o.StatusPath)
	if err != nil {
		o.Errorf("%v", err)
		return nil, err
	}
	collectTotalStats(mx, clients)
	if o.perUserMatcher != nil {
		o.collectUsers(mx, clients)
	}

	return mx, nil
}

func parseStatusLog(filePath string) ([]clientInfo, error) {
	conn, err := os.Open(filePath)
	defer conn.Close()

	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	if scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "OpenVPN CLIENT LIST") {
			return parseStatusLogV1(scanner), nil
		} else if strings.Contains(line, "TITLE,Open") {
			return parseStatusLogV2(scanner), nil
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
		switch parts := strings.Split(line, ","); parts[0] {
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

func collectTotalStats(mx map[string]int64, clients []clientInfo) {
	bytesIn := 0
	bytesOut := 0
	for _, c := range clients {
		bytesIn += c.BytesReceived
		bytesOut += c.BytesSent
	}
	mx["clients"] = int64(len(clients))
	mx["bytes_in"] = int64(bytesIn)
	mx["bytes_out"] = int64(bytesOut)
}

func (o *OpenVPNStatus) collectUsers(mx map[string]int64, clients []clientInfo) {
	now := time.Now().Unix()

	for _, user := range clients {
		name := user.CommonName
		if !o.perUserMatcher.MatchString(name) {
			continue
		}
		if !o.collectedUsers[name] {
			o.collectedUsers[name] = true
			if err := o.addUserCharts(name); err != nil {
				o.Warning(err)
			}
		}
		mx[name+"_bytes_in"] = int64(user.BytesReceived)
		mx[name+"_bytes_out"] = int64(user.BytesSent)
		mx[name+"_connection_time"] = now - user.ConnectedSince
	}
}

func (o *OpenVPNStatus) addUserCharts(userName string) error {
	cs := userCharts.Copy()

	for _, chart := range *cs {
		chart.ID = fmt.Sprintf(chart.ID, userName)
		chart.Fam = fmt.Sprintf(chart.Fam, userName)

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, userName)
		}
		chart.MarkNotCreated()
	}
	return o.charts.Add(*cs...)
}
