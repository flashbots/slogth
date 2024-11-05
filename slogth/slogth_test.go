package slogth

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/flashbots/slogth/mock"
	"github.com/stretchr/testify/assert"
)

func TestSlogth(t *testing.T) {
	stdin := mock.NewStdio()
	stdout := mock.NewStdio()

	s := new()

	s.input = stdin
	s.output = stdout
	s.delay = time.Second

	go func() {
		if err := s.run(); err != nil {
			assert.NoError(t, err)
		}
	}()

	{
		start := time.Now()
		_, err := stdin.Println(start.String())
		assert.NoError(t, err)

		var b []byte
		_, err = stdout.Read(b)
		assert.NoError(t, err)
		assert.Greater(t, time.Since(start), s.delay)
	}

	{
		count := 1000
		go func() {
			for countdown := count; countdown > 0; countdown-- {
				ts := int(time.Now().Unix())
				str := fmt.Sprintf("%d/%d", ts, countdown)
				_, err := stdin.Println(str)
				assert.NoError(t, err)
				time.Sleep(1 * time.Millisecond)
			}
			ts := int(time.Now().Unix())
			str := fmt.Sprintf("%d/%d", ts, 64)
			_, err := stdin.Println(str)
			assert.NoError(t, err)
			time.Sleep(5 * time.Millisecond)
		}()

		for countdown := count; countdown > 0; countdown-- {
			b := make([]byte, 64)
			l, err := stdout.Read(b)
			assert.NoError(t, err)
			str := strings.TrimSpace(string(b[:l]))
			parts := strings.Split(str, "/")
			assert.Equal(t, 2, len(parts), "%d: unexpected message '%s'", countdown, str)
			_countdown, err := strconv.Atoi(parts[1])
			assert.NoError(t, err)
			assert.Equal(t, countdown, _countdown)
			_ts, err := strconv.Atoi(parts[0])
			assert.NoError(t, err)
			ts := time.Unix(int64(_ts), 0)
			delay := time.Since(ts)
			assert.Greater(t, delay, s.delay, "%s/%d: %s <= %s", str, countdown, delay, s.delay)
		}

		assert.Equal(t, 0, s.queue.Length())

		t.Logf("capacity at the end: %d", s.queue.Capacity())
	}
}
