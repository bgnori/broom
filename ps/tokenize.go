package ps

import (
    "bufio"
    "fmt"
    "io"
    "unicode/utf8"
)

const BUFFERSIZE = 100

const (
    TOKEN_ENDOFINPUT= iota
    TOKEN_CHUNK= iota
    TOKEN_LEFTPAREN = iota
    TOKEN_RIGHTPAREN = iota
    TOKEN_DOT = iota
)


type Token struct {
    tkid int
    start int
    end int
    value Value
}


type Tokenizer struct {
    scanner *bufio.Scanner
    buffer []rune
    pos int
}

func NewTokenizer(r io.Reader) *Tokenizer{
    scanner := bufio.NewScanner(r)
    scanner.Split(bufio.ScanRunes)
    tknz := &Tokenizer{scanner:scanner, buffer:make([]rune,0), pos:0}
    tknz.PopulateBuffer(BUFFERSIZE)
    return tknz
}

func (tknz *Tokenizer)PopulateBuffer(upto int) int {
    i := 0
    for tknz.scanner.Scan() {
        r, _ := utf8.DecodeRuneInString(tknz.scanner.Text())
        tknz.buffer = append(tknz.buffer, r)
        if len(tknz.buffer) >= BUFFERSIZE || i >= upto {
            return i
        }
    }
    return i
}

func (tknz *Tokenizer)Peek() (r rune, eos bool) {
    if len(tknz.buffer) > 0 {
        return tknz.buffer[0], false
    } else {
        if tknz.PopulateBuffer(BUFFERSIZE) > 0 {
            return tknz.buffer[0], false
        }
        return 0, true
    }
}

func (tknz *Tokenizer)Consume(n int) {
    tknz.buffer = tknz.buffer[n:]
    tknz.pos += n
    tknz.PopulateBuffer(BUFFERSIZE)
}

func (tknz *Tokenizer)EndOfInput() *Token {
    return &Token{tkid: TOKEN_ENDOFINPUT, start:tknz.pos, end:tknz.pos}
}


func (tknz *Tokenizer)MakeLeftParen() *Token {
    t := &Token{start:tknz.pos, end:tknz.pos+1, tkid:TOKEN_LEFTPAREN}
    tknz.Consume(1)
    return t
}

func (tknz *Tokenizer)MakeRightParen() *Token {
    t := &Token{start:tknz.pos, end:tknz.pos+1, tkid:TOKEN_RIGHTPAREN}
    tknz.Consume(1)
    return t
}

func (tknz *Tokenizer)MakeDot() *Token {
    t := &Token{tkid: TOKEN_DOT, start:tknz.pos, end:tknz.pos+1}
    tknz.Consume(1)
    return t
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


func isWhite(r rune) bool {
    return r == ' ' ||  r =='\t' || r == '\v' || r =='\r' || r == '\n'
}


func parseInt(s string) (Value, error) {
    var n int
    _, err := fmt.Sscanf(s, "%d", &n)
    return Int(n), err
}

func parseName(s string) (Value, error) {
    return Name(s), nil
}

func NewChunk(start, end int, s string) *Token {
    v, err := parseInt(s)
    if err != nil {
        v, _ = parseName(s)
    }
    return &Token{start:start, end: end, tkid: TOKEN_CHUNK, value: v}
}

func (tknz *Tokenizer)MakeChunk() *Token {
    i := 0
    for ; i < len(tknz.buffer); i++{
        r := tknz.buffer[i]
        if isLeftParen(r) || isRightParen(r) || isWhite(r) {
            t := NewChunk(tknz.pos, tknz.pos+i, string(tknz.buffer[0:i]))
            tknz.Consume(i)
            return t
        }
    }
    t := NewChunk(tknz.pos, tknz.pos+len(tknz.buffer), string(tknz.buffer[:]))
    tknz.Consume(len(tknz.buffer))
    return t
}

func (tknz *Tokenizer)Tokenize() chan *Token {
    ch := make(chan *Token)
    go func(){
        defer close(ch)
        for {
            if token := tknz.OneToken(); token.tkid == TOKEN_ENDOFINPUT {
                ch<- token
                break
            } else {
                ch<- token
            }
        }
    }()
    return ch
}

func (tknz *Tokenizer)OneToken() *Token {
    r, eos := tknz.Peek()
    if eos {
        return tknz.EndOfInput()
    }
    for isWhite(r) {
        tknz.Consume(1)
        r, eos = tknz.Peek()
        if eos {
            return tknz.EndOfInput()
        }
    }
    switch {
    case isLeftParen(r) :
        return tknz.MakeLeftParen()
    case isRightParen(r):
        return tknz.MakeRightParen()
    case isDot(r):
        return tknz.MakeDot()
    default:
        return tknz.MakeChunk()
    }
    panic("Should not reach here")
}
