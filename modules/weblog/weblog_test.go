package weblog

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/netdata/go-orchestrator/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.SetSeverity(logger.DEBUG)
}

func TestWebLog_Init(t *testing.T) {
	weblog := New()
	weblog.Config.Filter.Includes = []string{"~ .php$"}
	weblog.Config.URLCategories = []rawCategory{
		{"foo", "= foo"},
		{"bar", "= bar"},
	}
	weblog.Config.UserCategories = []rawCategory{
		{"baz", "= baz"},
		{"foobar", "= foobar"},
	}
	ok := weblog.Init()

	require.True(t, ok)
	assert.True(t, weblog.filter.MatchString("/abc.php"))
	assert.False(t, weblog.filter.MatchString("/abc.html"))

	assert.Len(t, weblog.urlCats, 2)
	assert.Equal(t, "foo", weblog.urlCats[0].name)
	assert.True(t, weblog.urlCats[0].Matcher.MatchString("foo"))
	assert.Equal(t, "bar", weblog.urlCats[1].name)
	assert.True(t, weblog.urlCats[1].Matcher.MatchString("bar"))

	assert.Len(t, weblog.userCats, 2)
	assert.Equal(t, "baz", weblog.userCats[0].name)
	assert.True(t, weblog.userCats[0].Matcher.MatchString("baz"))
	assert.Equal(t, "foobar", weblog.userCats[1].name)
	assert.True(t, weblog.userCats[1].Matcher.MatchString("foobar"))
}

func TestWebLog_Collect(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer os.Remove(tmp.Name())

	done, wait := generateLog(tmp)
	time.Sleep(150 * time.Millisecond)

	weblog := New()
	weblog.URLCategories = []rawCategory{{"BOOKS", "* *"}}
	defer weblog.Cleanup()

	weblog.Config.Path = tmp.Name()

	ok := weblog.Init()
	require.True(t, ok)
	ok = weblog.Check()
	require.True(t, ok)

	time.Sleep(150 * time.Millisecond)

	m := weblog.Collect()
	l := make([]string, 0)
	for k := range m {
		l = append(l, k)
	}
	sort.Strings(l)
	for _, v := range l {
		fmt.Println(fmt.Sprintf("\"%s\": %d,", v, m[v]))
	}

	done <- 1
	<-wait
}

func generateLog(f *os.File) (done chan<- int, wait <-chan int) {
	doneChan := make(chan int, 1)
	waitChan := make(chan int, 1)
	go func() {
		for {
			select {
			case <-doneChan:
				f.Close()
				waitChan <- 1
				return
			case <-time.After(time.Millisecond * 100):
				fmt.Fprintln(f, `127.0.0.1 - - [28/Jan/2019:11:18:12 +0900] "GET /order/books HTTP/1.1" 301 6295 "https://www.test.com/order" "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko" - 12625`)
				fmt.Fprintln(f, `127.0.0.1 - - [28/Jan/2019:11:18:12 +0900] "GET /order/books HTTP/1.1" 200 6295 "https://www.test.com/order" "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko" - 12625`)
			}
		}
	}()
	return doneChan, waitChan
}
