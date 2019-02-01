package testhelpers

import (
	"bytes"
	"io"
	"sync"
)

type SyncedBuffer struct {
	b  bytes.Buffer
	mu sync.Mutex
}

func (b *SyncedBuffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Bytes()
}

func (b *SyncedBuffer) Cap() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Cap()
}

func (b *SyncedBuffer) Grow(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.b.Grow(n)
}

func (b *SyncedBuffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Len()
}

func (b *SyncedBuffer) Next(n int) []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Next(n)
}

func (b *SyncedBuffer) Read(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Read(p)
}

func (b *SyncedBuffer) ReadByte() (c byte, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadByte()
}

func (b *SyncedBuffer) ReadBytes(delim byte) (line []byte, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadBytes(delim)
}

func (b *SyncedBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadFrom(r)
}

func (b *SyncedBuffer) ReadRune() (r rune, size int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadRune()
}

func (b *SyncedBuffer) ReadString(delim byte) (line string, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.ReadString(delim)
}

func (b *SyncedBuffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.b.Reset()
}

func (b *SyncedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.String()
}

func (b *SyncedBuffer) Truncate(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.b.Truncate(n)
}

func (b *SyncedBuffer) UnreadByte() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.UnreadByte()
}

func (b *SyncedBuffer) UnreadRune() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.UnreadRune()
}

func (b *SyncedBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Write(p)
}

func (b *SyncedBuffer) WriteByte(c byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteByte(c)
}

func (b *SyncedBuffer) WriteRune(r rune) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteRune(r)
}

func (b *SyncedBuffer) WriteString(s string) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteString(s)
}

func (b *SyncedBuffer) WriteTo(w io.Writer) (n int64, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteTo(w)
}
