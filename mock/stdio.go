package mock

type Stdio struct {
	buffer chan []byte
}

func NewStdio() *Stdio {
	return &Stdio{
		buffer: make(chan []byte, 16),
	}
}

func (m *Stdio) Read(b []byte) (int, error) {
	data := <-m.buffer
	return copy(b, data), nil
}

func (m *Stdio) Write(b []byte) (int, error) {
	msg := append([]byte{}, b...)
	m.buffer <- msg
	return len(msg), nil
}

func (m *Stdio) Println(s string) (int, error) {
	msg := append([]byte{}, s...)
	msg = append(msg, '\n')
	m.buffer <- msg
	return len(msg), nil
}
