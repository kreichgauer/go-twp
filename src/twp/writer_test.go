package twp

import (
    "bytes"
    "encoding/binary"
    "strings"
    "testing"
)

// FIXME
func EncodeString(t *testing.T) {
    var buf, buf2 bytes.Buffer
    en := NewEncoder(&buf)
    hw := "Hello, World!\n"
    if err := en.EncodeString(hw); err != nil {
        pFatal(err, t)
    }
    buf2.WriteByte(byte(ShortString + len(hw)))
    buf2.WriteString(hw)
    verify(buf.Bytes(), buf2.Bytes(), t)

    buf.Reset()
    buf2.Reset()
    hw = strings.Repeat(hw, 100)
    if err := en.EncodeString(hw); err != nil {
        pFatal(err, t)
    }
    buf2.WriteByte(LongString)
    binary.Write(&buf2, binary.BigEndian, int32(len(hw)))
    buf2.WriteString(hw)
    verify(buf.Bytes(), buf2.Bytes(), t)
}
