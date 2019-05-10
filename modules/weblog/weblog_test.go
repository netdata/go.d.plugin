package weblog

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebLog_Init(t *testing.T) {
	weblog := New()
	weblog.Config.Filter.Include = "~ .php$"
	weblog.Config.URLCategories = []RawCategory{
		{"foo", "= foo"},
		{"bar", "= bar"},
	}
	weblog.Config.UserCategories = []RawCategory{
		{"baz", "= baz"},
		{"foobar", "= foobar"},
	}
	ok := weblog.Init()

	require.True(t, ok)
	assert.True(t, weblog.filter.MatchString("/abc.php"))
	assert.False(t, weblog.filter.MatchString("/abc.html"))

	assert.Len(t, weblog.urlCategories, 2)
	assert.Equal(t, "foo", weblog.urlCategories[0].name)
	assert.True(t, weblog.urlCategories[0].Matcher.MatchString("foo"))
	assert.Equal(t, "bar", weblog.urlCategories[1].name)
	assert.True(t, weblog.urlCategories[1].Matcher.MatchString("bar"))

	assert.Len(t, weblog.userCategories, 2)
	assert.Equal(t, "baz", weblog.userCategories[0].name)
	assert.True(t, weblog.userCategories[0].Matcher.MatchString("baz"))
	assert.Equal(t, "foobar", weblog.userCategories[1].name)
	assert.True(t, weblog.userCategories[1].Matcher.MatchString("foobar"))
}

func TestWebLog_Collect(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer os.Remove(tmp.Name())

	done, wait := generateLog(tmp)
	time.Sleep(150 * time.Millisecond)

	weblog := New()
	weblog.Config.Path = tmp.Name()

	ok := weblog.Init()
	require.True(t, ok)
	ok = weblog.Check()
	require.True(t, ok)

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
