package openvpn_status

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	//"time"
)

type clientInfo struct {
	CommonName    string
	BytesReceived int
	BytesSent     int
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

func (o *OpenVPNStatus) collect() (map[string]int64, error) {
	var err error

	mx := make(map[string]int64)

	clients, err := parseStatusLog(o.StatusPath)
	if err != nil {
		o.Errorf("%v", err)
		return nil, err
	}
	collectTotalStats(mx, clients)

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
		if checkClientListHeader(parts) {
			clientListHeader = true
			continue
		}
		if clientListHeader && len(parts) == 5 {
			in, _ := strconv.Atoi(parts[2])
			out, _ := strconv.Atoi(parts[3])
			c := clientInfo{
				CommonName:    parts[0],
				BytesReceived: in,
				BytesSent:     out,
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
				c := clientInfo{
					CommonName:    parts[1],
					BytesReceived: in,
					BytesSent:     out,
				}
				clients = append(clients, c)
			}
		}
	}
	return clients
}

func (o *OpenVPNStatus) collectUsers(mx map[string]int64, clients []clientInfo) error {
	//now := time.Now().Unix()
	var name string

	for _, user := range clients {
		if !o.perUserMatcher.MatchString(user.CommonName) {
			continue
		}
		if !o.collectedUsers[user.CommonName] {
			o.collectedUsers[user.CommonName] = true
			if err := o.addUserCharts(user.CommonName); err != nil {
				o.Warning(err)
			}
		}
		mx[name+"_bytes_received"] = int64(user.BytesReceived)
		mx[name+"_bytes_sent"] = int64(user.BytesSent)
		//mx[name+"_connection_time"] = now - user.ConnectedSince
	}
	return nil
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
