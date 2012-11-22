package twp

import (
    "bytes"
    "fmt"
    "reflect"
)

type Value interface{}

type Message struct {
    Id uint8
}

type RawMessage struct {
    Message
    Fields []Value
}

func (msg *RawMessage) String() (string) {
    var b bytes.Buffer
    b.WriteString(fmt.Sprintf("Message %d\n", msg.Id))
    for _, val := range msg.Fields {
        b.WriteString(fmt.Sprintf("%s: %v\n", reflect.TypeOf(val).Name(), val))
    }
    return b.String()
}
