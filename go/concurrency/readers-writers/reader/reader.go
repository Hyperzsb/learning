package reader

import "fmt"

type Reader struct {
	Id     int
	Buffer *int
}

func NewReader(id int, buffer *int) Reader {
	return Reader{Id: id, Buffer: buffer}
}

func (r Reader) Read() {
	fmt.Printf("Reader %d reads %d\n", r.Id, *r.Buffer)
}
