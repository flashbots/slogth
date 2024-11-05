package slogth

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flashbots/slogth/types"
)

type slogth struct {
	input  io.Reader
	output io.Writer

	delay time.Duration

	isTerminating bool
	queue         *types.TimedQueue[[]byte]
}

func new() *slogth {
	return &slogth{
		queue: types.NewTimedQueue[[]byte](),
	}
}

func (s *slogth) run() error {
	type entry struct {
		emitAt time.Time
		log    []byte
	}

	var (
		tick      = time.NewTicker(min(s.delay/10, time.Second)) // emit at least every second
		input     = make(chan entry)
		fail      = make(chan error, 1)
		terminate = make(chan os.Signal, 1)
		err       error
	)

	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)

	go func() { // read from stdin
		scanner := bufio.NewScanner(s.input)
		for scanner.Scan() {
			input <- entry{
				emitAt: time.Now().Add(s.delay),
				log:    append([]byte{}, scanner.Bytes()...), // have to copy over
			}
		}
		if !s.isTerminating {
			if err := scanner.Err(); err != nil {
				fail <- err
			}
		}
	}()

loop: // emit logs after delay
	for {
		select {
		case now := <-tick.C:
			s.queue.PopBefore(now, func(log []byte) {
				fmt.Fprintln(s.output, string(log))
			})
		case entry := <-input:
			s.queue.Push(entry.emitAt, entry.log)
		case <-terminate:
			s.isTerminating = true
			break loop
		case err = <-fail:
			break loop
		}
	}

	// flush the queue
	s.queue.Pop(func(timestamp time.Time, log []byte) {
		time.Sleep(time.Until(timestamp))
		fmt.Fprintln(s.output, string(log))
	})

	return err
}
