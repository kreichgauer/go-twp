package twp

type Tag uint8

const (
    EndOfContent        = 0
    NoValue             = 1
    Struct              = 2
    Sequence            = 3
    MessageOrUnion      = 4
    MessageOrUnionEnd   = 11
    RegisteredExtension = 12
    ShortInteger        = 13
    LongInteger         = 14
    ShortBinary         = 15
    LongBinary          = 16
    ShortString         = 17
    LongString          = 127
    Reserved            = 128
    ReservedEnd         = 159
    ApplicationType     = 160
    ApplicationTypeEnd  = 255
)
