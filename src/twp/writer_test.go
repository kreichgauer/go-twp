package twp

import (
    "bytes"
    "encoding/binary"
    "strings"
    "testing"
)

func WriteString(t *testing.T) {
    var buf, buf2 bytes.Buffer
    wr := NewWriter(&buf)
    hw := "Hello, World!\n"
    if err := wr.WriteString(hw); err != nil {
        pFatal(err, t)
    }
    buf2.WriteByte(byte(ShortString + len(hw)))
    buf2.WriteString(hw)
    verify(buf.Bytes(), buf2.Bytes(), t)

    buf.Reset()
    buf2.Reset()
    hw = strings.Repeat(hw, 100)
    if err := wr.WriteString(hw); err != nil {
        pFatal(err, t)
    }
    buf2.WriteByte(LongString)
    binary.Write(&buf2, binary.BigEndian, int32(len(hw)))
    buf2.WriteString(hw)
    verify(buf.Bytes(), buf2.Bytes(), t)
}
