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

	delay         time.Duration
	dropThreshold int

	isTerminating bool
	queue         *types.TimedQueue[[]byte]

	droppedCount   int
	processedCount int
}

func new() *slogth {
	return &slogth{
		queue: types.NewTimedQueue[[]byte](),
	}
}

type entry struct {
	emitAt time.Time
	log    []byte
}

func (s *slogth) run() error {
	var (
		input <-chan entry
		fail  <-chan error
		err   error

		tick      = time.NewTicker(min(s.delay/10, time.Second)) // emit at least every second
		terminate = make(chan os.Signal, 1)
	)

	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)

	if s.dropThreshold == 0 {
		input, fail = s.ingestBlocking()
	} else {
		input, fail = s.ingestNonBlocking()
	}

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

func (s *slogth) ingestBlocking() (<-chan entry, <-chan error) {
	input := make(chan entry, s.dropThreshold)
	fail := make(chan error, 1)

	go func() { // read from stdin
		scanner := bufio.NewScanner(s.input)
		for scanner.Scan() {
			input <- entry{
				emitAt: time.Now().Add(s.delay),
				log:    append([]byte{}, scanner.Bytes()...), // have to copy over
			}
			s.processedCount++
		}
		if !s.isTerminating {
			if err := scanner.Err(); err != nil {
				fail <- err
			}
		}
	}()

	return input, fail
}

func (s *slogth) ingestNonBlocking() (<-chan entry, <-chan error) {
	input := make(chan entry, s.dropThreshold)
	fail := make(chan error, 1)

	go func() { // read from stdin
		scanner := bufio.NewScanner(s.input)
		for scanner.Scan() {
			select {
			case input <- entry{
				emitAt: time.Now().Add(s.delay),
				log:    append([]byte{}, scanner.Bytes()...), // have to copy over
			}:
				s.processedCount++
			default:
				s.droppedCount++
			}
		}
		if !s.isTerminating {
			if err := scanner.Err(); err != nil {
				fail <- err
			}
		}
	}()

	return input, fail
}
