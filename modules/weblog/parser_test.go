package weblog

//func Test_logParser_ParseLine(t *testing.T) {
//	logs, _ := os.Open("tests/common.log")
//	parser := NewCSVParser(ParserConfig{}, logs)
//	format := common
//
//	for i := 0; ; i++ {
//		fields, err := parser.Read()
//		if err == io.EOF {
//			assert.Equal(t, 50, i)
//			break
//		}
//		row, err := format.Parse(fields)
//		require.NoError(t, err)
//		assert.Equal(t, "", row.Host)
//		assert.Equal(t, -1.0, row.ReqTime)
//		assert.Nil(t, row.UpstreamRespTime)
//		assert.Equal(t, -1, row.ReqLength)
//		assert.Equal(t, "", row.Custom)
//		switch i {
//		case 0: // 64.242.88.10 - - [07/Mar/2004:16:05:49 -0800] "GET /twiki/bin/edit/Main/Double_bounce_sender?topicparent=Main.ConfigurationVariables HTTP/1.1" 401 12846
//			assert.Equal(t, "64.242.88.10", row.RemoteAddr)
//			assert.Equal(t, "GET /twiki/bin/edit/Main/Double_bounce_sender?topicparent=Main.ConfigurationVariables HTTP/1.1", row.Request)
//			assert.Equal(t, "GET", row.Method)
//			assert.Equal(t, "/twiki/bin/edit/Main/Double_bounce_sender?topicparent=Main.ConfigurationVariables", row.URI)
//			assert.Equal(t, "1.1", row.Version)
//			assert.Equal(t, 401, row.Status)
//			assert.Equal(t, 12846, row.BytesSent)
//		case 1: // 64.242.88.10 - - [07/Mar/2004:16:06:51 -0800] "GET /twiki/bin/rdiff/TWiki/NewUserTemplate?rev1=1.3&rev2=1.2 HTTP/1.1" 200 4523
//			assert.Equal(t, "64.242.88.10", row.RemoteAddr)
//			assert.Equal(t, "GET /twiki/bin/rdiff/TWiki/NewUserTemplate?rev1=1.3&rev2=1.2 HTTP/1.1", row.Request)
//			assert.Equal(t, "GET", row.Method)
//			assert.Equal(t, "/twiki/bin/rdiff/TWiki/NewUserTemplate?rev1=1.3&rev2=1.2", row.URI)
//			assert.Equal(t, "1.1", row.Version)
//			assert.Equal(t, 200, row.Status)
//			assert.Equal(t, 4523, row.BytesSent)
//		case 33: // lj1036.inktomisearch.com - - [07/Mar/2004:17:18:36 -0800] "GET /robots.txt HTTP/1.0" 200 68
//			assert.Equal(t, "lj1036.inktomisearch.com", row.RemoteAddr)
//			assert.Equal(t, "GET /robots.txt HTTP/1.0", row.Request)
//			assert.Equal(t, 200, row.Status)
//			assert.Equal(t, 68, row.BytesSent)
//		case 49: // 64.242.88.10 - - [07/Mar/2004:17:53:45 -0800] "GET /twiki/bin/search/Main/SearchResult?scope=text®ex=on&search=Office%20*Locations[^A-Za-z] HTTP/1.1" 200 7771
//			assert.Equal(t, "64.242.88.10", row.RemoteAddr)
//			assert.Equal(t, "GET /twiki/bin/search/Main/SearchResult?scope=text®ex=on&search=Office%20*Locations[^A-Za-z] HTTP/1.1", row.Request)
//			assert.Equal(t, "GET", row.Method)
//			assert.Equal(t, "/twiki/bin/search/Main/SearchResult?scope=text®ex=on&search=Office%20*Locations[^A-Za-z]", row.URI)
//			assert.Equal(t, "1.1", row.Version)
//			assert.Equal(t, 200, row.Status)
//			assert.Equal(t, 7771, row.BytesSent)
//		}
//	}
//}
//
//func BenchmarkCSVParser(b *testing.B) {
//	content, err := ioutil.ReadFile("testdata/big_combined.log")
//	if err != nil {
//		b.Fatalf("%+v", err)
//	}
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		buf := bytes.NewBuffer(content)
//		parser := NewCSVParser(ParserConfig{}, buf)
//		for err == nil {
//			_, err = parser.Read()
//		}
//	}
//}
//
//func BenchmarkCSVParser2(b *testing.B) {
//	content, err := ioutil.ReadFile("testdata/big_combined.log")
//	if err != nil {
//		b.Fatalf("%+v", err)
//	}
//
//	var lines [][]byte
//	buf := bytes.NewBuffer(content)
//	scanner := bufio.NewScanner(buf)
//	for scanner.Scan() {
//		lines = append(lines, scanner.Bytes())
//	}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		for _, line := range lines {
//			buf := bytes.NewBuffer(line)
//			parser := NewCSVParser(ParserConfig{}, buf)
//			_, err = parser.Read()
//		}
//	}
//}
