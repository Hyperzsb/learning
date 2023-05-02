package writer

import "fmt"

type Writer struct {
	Id     int
	Buffer *int
}

func NewWriter(id int, buffer *int) Writer {
	return Writer{Id: id, Buffer: buffer}
}

func (w *Writer) Write(content int) {
	*w.Buffer = (w.Id+1)*1000 + content
	fmt.Printf("Writer %d writes %d\n", w.Id, *w.Buffer)
}
