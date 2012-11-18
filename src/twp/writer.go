package twp

import (
    "bufio"
    "encoding/binary"
    "io"
)

type Writer struct {
    bufio.Writer
}

const twpMagic = "TWP3\n"

func NewWriter(wr io.Writer) (*Writer) {
    return &Writer{*bufio.NewWriter(wr)}
}

func (wr *Writer) InitWithProtocol(id int) (err error) {
    if err = wr.WriteMagic(); err != nil {
        return err
    }
    if err = wr.WriteProtocolId(id); err != nil {
        return err
    }
    if err = wr.Flush(); err != nil {
        return err
    }
    return nil
}

func (wr *Writer) WriteMagic() (err error) {
    _, err = wr.Write([]byte(twpMagic))
    return err
}

func (wr *Writer) WriteProtocolId(id int) (err error) {
    return wr.WriteShortInt(id)
}

func (wr *Writer) WriteShortInt(val int) (err error) {
    return binary.Write(wr, binary.BigEndian, []byte{13, byte(val)})
}
