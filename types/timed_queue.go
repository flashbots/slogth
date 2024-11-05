package types

import "time"

type TimedQueue[T any] struct {
	head  int // head is there the new element would be pushed (written to).
	tail  int // tail is where the next popped element would be read from.
	queue []element[T]
}

type element[T any] struct {
	timestamp time.Time
	value     T
}

func NewTimedQueue[T any]() *TimedQueue[T] {
	return &TimedQueue[T]{
		queue: make([]element[T], 64),
	}
}

func (q *TimedQueue[T]) Capacity() int {
	return cap(q.queue)
}

func (q *TimedQueue[T]) Length() int {
	switch {
	case q.head >= q.tail:
		return q.head - q.tail
	case q.head < q.tail:
		return len(q.queue) + q.head - q.tail
	default:
		return 0
	}
}

func (q *TimedQueue[T]) Push(timestamp time.Time, value T) {
	switch {
	case q.head >= q.tail:
		//  0   1   2   3
		//  t       h      : 2 elements
		size := q.head - q.tail
		if size == len(q.queue)-1 { // time to grow
			newQueue := make([]element[T], 2*len(q.queue))
			copy(newQueue, q.queue)
			q.queue = newQueue
		}
		q.queue[q.head] = element[T]{
			timestamp: timestamp,
			value:     value,
		}
		q.head++
		if q.head == len(q.queue) {
			q.head = 0
		}
	case q.head < q.tail:
		//  0   1   2   3
		//  h           t  : 1 element
		size := len(q.queue) + q.head - q.tail
		if size == len(q.queue)-1 { // time to grow
			newQueue := make([]element[T], 2*len(q.queue))
			copy(newQueue, q.queue[q.tail:])
			copy(newQueue[len(q.queue)-q.tail:], q.queue[:q.head])
			q.queue = newQueue
			q.tail = 0
			q.head = size
		}
		q.queue[q.head] = element[T]{
			timestamp: timestamp,
			value:     value,
		}
		q.head++
	}
}

func (q *TimedQueue[T]) PopBefore(timestamp time.Time, do func(T)) {
	for q.head != q.tail {
		switch {
		case q.head >= q.tail:
			e := q.queue[q.tail]
			if !e.timestamp.Before(timestamp) {
				return
			}
			q.tail++
			do(e.value)
		default:
			e := q.queue[q.tail]
			if !e.timestamp.Before(timestamp) {
				return
			}
			q.tail++
			if q.tail == len(q.queue) {
				q.tail = 0
			}
			do(e.value)
		}
	}
}

func (q *TimedQueue[T]) Pop(do func(timestamp time.Time, value T)) {
	for q.head != q.tail {
		switch {
		case q.head >= q.tail:
			e := q.queue[q.tail]
			q.tail++
			do(e.timestamp, e.value)
		default:
			e := q.queue[q.tail]
			q.tail++
			if q.tail == len(q.queue) {
				q.tail = 0
			}
			do(e.timestamp, e.value)
		}
	}
}
