package broom

import (
    "bufio"
    "io"
    "unicode/utf8"
)
const BUFFERSIZE = 100

type Buffered struct {
    // I wonder that why bufio itself is not enough??
    scanner *bufio.Scanner
    buffer []rune
    pos int
}

func NewBuffered(r io.Reader) *Buffered{
    scanner := bufio.NewScanner(r)
    scanner.Split(bufio.ScanRunes)
    buffer := &Buffered{scanner:scanner, buffer:make([]rune,0), pos:0}
    buffer.PopulateBuffer(BUFFERSIZE)
    return buffer
}

func (buffer *Buffered)PopulateBuffer(upto int) int {
    i := 0
    for buffer.scanner.Scan() {
        r, _ := utf8.DecodeRuneInString(buffer.scanner.Text())
        buffer.buffer = append(buffer.buffer, r)
        if len(buffer.buffer) >= BUFFERSIZE || i >= upto {
            return i
        }
    }
    return i
}

func (buffer *Buffered)Peek() (r rune, eos bool) {
    if len(buffer.buffer) > 0 {
        return buffer.buffer[0], false
    } else {
        if buffer.PopulateBuffer(BUFFERSIZE) > 0 {
            return buffer.buffer[0], false
        }
        return 0, true
    }
    panic("never reach")
}

func (buffer *Buffered)Consume(n int) {
    if k := len(buffer.buffer) ; k < n {
        panic("not enough content to consume")
    }
    buffer.buffer = buffer.buffer[n:]
    buffer.pos += n
    buffer.PopulateBuffer(BUFFERSIZE)
}


func isDigit(r rune) bool {
    return '0' <= r && r <= '9'
}

func isLeftParen(r rune) bool {
    return '(' == r
}

func isRightParen(r rune) bool {
    return ')' == r
}

func isDot(r rune) bool {
    return '.' == r
}

func isDoubleQuote(r rune) bool {
    return '"' == r
}

func isBackSlash(r rune) bool {
    return '\\' == r
}

func isWhite(r rune) bool {
    return r == ' ' ||  r =='\t' || r == '\v' || r =='\r' || r == '\n'
}


type Reader struct {
    buffer *Buffered
    out chan Token
    state ReaderState
}

type ReaderState func(reader *Reader) ReaderState

func NewReader(buffer *Buffered) *Reader {
    return &Reader{buffer: buffer, out: make(chan Token, 2), state: TopLevel}
}

func (reader *Reader) Read () Token{
    for { //reader.state != nil {
        select {
        case tkn := <-reader.out:
            return tkn
        default:
            reader.state = reader.state (reader)
        }
    }
    panic("nil state")
}

func (reader *Reader) Emit (token Token){
    reader.out <- token
}

func (reader *Reader) ZapWhite() (rune, error) {
    r, eos := reader.buffer.Peek()
    if eos {
        return 0, io.EOF
    }
    for isWhite(r) {
        reader.buffer.Consume(1)
        r, eos = reader.buffer.Peek()
        if eos {
            return 0, io.EOF
        }
    }
    return r, nil
}

func TopLevel(reader *Reader) ReaderState {
    r, err := reader.ZapWhite()
    if err != nil {
        reader.Emit(Token{id:TOKEN_ENDOFINPUT})
        return nil
    }
    switch {
    case isLeftParen(r) :
        reader.Emit(reader.MakeLeftParen())
        return TopLevel
    case isRightParen(r):
        reader.Emit(reader.MakeRightParen())
        return TopLevel
    case isDoubleQuote(r):
        reader.Emit(reader.tryString())
        return TopLevel
    default:
        reader.Emit(reader.tryChunk())
        return TopLevel
    }
    panic("Should not reach here")
}

func (reader *Reader) tryString () Token {
    xs := make([]rune,0)
    pos := reader.buffer.pos
    reader.buffer.Consume(1) // skip "
    r, eos := reader.buffer.Peek()
    escaped := false
    for !eos && !(!escaped && isDoubleQuote(r)) {
        if !escaped && isBackSlash(r) {
            escaped = true
        } else {
            xs = append(xs, r)
            escaped = false
        }
        reader.buffer.Consume(1)
        r, eos = reader.buffer.Peek()
    }
    if eos {
        panic("string must be closed")
    }
    if isDoubleQuote(r) {
        reader.buffer.Consume(1) // skip "
    }
    return Token{id: TOKEN_STRING, v:string(xs), pos:pos}
}

func (reader *Reader) tryChunk () Token {
    xs := make([]rune,0)
    pos := reader.buffer.pos
    r, eos := reader.buffer.Peek()
    for !eos && !isWhite(r) && !isLeftParen(r) && !isRightParen(r) {
        xs = append(xs, r)
        reader.buffer.Consume(1)
        r, eos = reader.buffer.Peek()
    }
    return Token{id: TOKEN_CHUNK, v:string(xs), pos:pos}
}

type Token struct {
    id int
    pos int
    v string
}

const (
    TOKEN_ENDOFINPUT= iota
    TOKEN_CHUNK //Number, Symbol
    TOKEN_STRING
    TOKEN_LEFTPAREN
    TOKEN_RIGHTPAREN
    TOKEN_DOT
)

func (r *Reader)MakeLeftParen() Token {
    t := Token{pos:r.buffer.pos, id:TOKEN_LEFTPAREN}
    r.buffer.Consume(1)
    return t
}

func (r *Reader)MakeRightParen() Token {
    t := Token{pos:r.buffer.pos, id:TOKEN_RIGHTPAREN}
    r.buffer.Consume(1)
    return t
}


