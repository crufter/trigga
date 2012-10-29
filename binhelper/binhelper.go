package binhelper

import(
	"io"
	"fmt"
	"encoding/binary"
)

type byteReader struct {
	b		[]byte
	pos		int
}

func (b *byteReader) ReadByte() (byte, error) {
	if b.pos >= len(b.b) {
		return 0, fmt.Errorf("No good.")
	}
	ret := b.b[b.pos]
	b.pos++
	return ret, nil
}

func bReader(b []byte) *byteReader {
	return &byteReader{
		b,
		0,
	}
}

func ReadMsg(c io.Reader) ([]byte, error) {
	l := make([]byte, 8)
	amt, err := c.Read(l)
	if err != nil {
		return nil, err
	}
	if amt != 8 {
		return nil, fmt.Errorf("Can't read data len. %v", amt)
	}
	datlen, err := binary.ReadVarint(bReader(l))
	if err != nil {
		return nil, err
	}
	buf := make([]byte, datlen)
	amt, err = c.Read(buf)
	if amt != int(datlen) {
		return nil, fmt.Errorf("Len mismatch.")
	}
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func WriteMsg(c io.Writer, msg []byte) error {
	var l int64
	l = int64(len(msg))
	lenstore := make([]byte, 8)
	binary.PutVarint(lenstore, l)
	alldat := []byte{}
	alldat = append(alldat, lenstore...)
	alldat = append(alldat, msg...)
	i, err := fmt.Fprintf(c, string(alldat))
	if err != nil {
		return err
	}
	if i != len(msg) + 8 {
		return fmt.Errorf("Error in writing: lens don't match.")
	}
	return nil
}