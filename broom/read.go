package broom

import (
	"bufio"
	"fmt"
	"io"
	"unicode/utf8"
)

const BUFFERSIZE = 100

type Buffered struct {
	// I wonder that why bufio itself is not enough??
	scanner *bufio.Scanner
	buffer  []rune
	pos     int
}

func NewBuffered(r io.Reader) *Buffered {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	buffer := &Buffered{scanner: scanner, buffer: make([]rune, 0), pos: 0}
	buffer.PopulateBuffer(BUFFERSIZE)
	return buffer
}

func (buffer *Buffered) PopulateBuffer(upto int) int {
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

func (buffer *Buffered) Peek() (r rune, eos bool) {
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

func (buffer *Buffered) Consume(n int) {
	if k := len(buffer.buffer); k < n {
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

func isLeftBracket(r rune) bool {
	return '[' == r
}

func isRightBracket(r rune) bool {
	return ']' == r
}

func isLeftBrace(r rune) bool {
	return '{' == r
}

func isRightBrace(r rune) bool {
	return '}' == r
}

func isParenBracketBrace(r rune) bool {
	return isLeftParen(r) || isRightParen(r) ||
		isLeftBracket(r) || isRightBracket(r) ||
		isLeftBrace(r) || isRightBrace(r)
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

func isVerticalBar(r rune) bool {
	return '|' == r
}

func isSharp(r rune) bool {
	return '#' == r
}

func isQuote(r rune) bool {
	return '\'' == r
}

func isQuasiQuote(r rune) bool {
	return '`' == r
}

func isWhite(r rune) bool {
	return r == ' ' || r == '\t' || r == '\v' || r == '\r' || r == '\n'
}

type Reader struct {
	buffer *Buffered
	out    chan Token
	state  ReaderState
}

type ReaderState func(reader *Reader) ReaderState

func NewReader(buffer *Buffered) *Reader {
	return &Reader{buffer: buffer, out: make(chan Token, 2), state: TopLevel}
}

func (reader *Reader) Read() Token {
	for { //reader.state != nil {
		select {
		case tkn := <-reader.out:
			return tkn
		default:
			reader.state = reader.state(reader)
		}
	}
	panic("nil state")
}

func (reader *Reader) Emit(token Token) {
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
		reader.Emit(Token{id: TOKEN_ENDOFINPUT})
		return nil
	}
	switch {
	case isLeftParen(r):
		reader.Emit(reader.MakeLeftParen())
		return TopLevel
	case isRightParen(r):
		reader.Emit(reader.MakeRightParen())
		return TopLevel
	case isLeftBracket(r):
		reader.Emit(reader.MakeLeftBracket())
		return TopLevel
	case isRightBracket(r):
		reader.Emit(reader.MakeRightBracket())
		return TopLevel
	case isLeftBrace(r):
		reader.Emit(reader.MakeLeftBrace())
		return TopLevel
	case isRightBrace(r):
		reader.Emit(reader.MakeRightBrace())
		return TopLevel
	case isDoubleQuote(r):
		reader.Emit(reader.tryString())
		return TopLevel
	case isVerticalBar(r):
		reader.Emit(reader.MakeVerticalVar())
		return TopLevel
	case isSharp(r):
		reader.Emit(reader.MakeSharp())
		return TopLevel
	default:
		reader.Emit(reader.tryChunk())
		return TopLevel
	}
	panic("Should not reach here")
}

func (reader *Reader) tryString() Token {
	xs := make([]rune, 0)
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
	return Token{id: TOKEN_STRING, v: string(xs), pos: pos}
}

func (reader *Reader) tryChunk() Token {
	xs := make([]rune, 0)
	pos := reader.buffer.pos
	r, eos := reader.buffer.Peek()
	for !eos && !isWhite(r) && !isParenBracketBrace(r) {
		xs = append(xs, r)
		reader.buffer.Consume(1)
		r, eos = reader.buffer.Peek()
	}
	s := string(xs)
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err == nil {
		return Token{id: TOKEN_INT, v: s, pos: pos}
	}
	return Token{id: TOKEN_SYMBOL, v: s, pos: pos}
}

type Token struct {
	id  int
	pos int
	v   string
}

const (
	TOKEN_ENDOFINPUT = iota
	TOKEN_SYMBOL
	TOKEN_INT
	TOKEN_STRING
	TOKEN_LEFT_PAREN
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACKET
	TOKEN_RIGHT_BRACKET
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_VERTICAL_BAR
	TOKEN_SHARP
	TOKEN_SEMICOLON
	TOKEN_COLON
	TOKEN_QUOTE
	TOKEN_QUASIQUOTE
	TOKEN_DOT
)

func (r *Reader) MakeLeftParen() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_LEFT_PAREN}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeRightParen() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_RIGHT_PAREN}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeLeftBracket() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_LEFT_BRACKET}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeRightBracket() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_RIGHT_BRACKET}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeLeftBrace() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_LEFT_BRACE}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeRightBrace() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_RIGHT_BRACE}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeVerticalVar() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_VERTICAL_BAR}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeSharp() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_SHARP}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeColon() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_COLON}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeSemicolon() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_SEMICOLON}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeQuote() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_QUOTE}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeQuasiQuote() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_QUASIQUOTE}
	r.buffer.Consume(1)
	return t
}

type tokenSeq struct {
	typ   int
	items []interface{}
}

func (t *tokenSeq) Items() []interface{} {
	return t.items
}

type SExprBuilder struct {
	stack []*tokenSeq
}

func NewSExprBuilder() *SExprBuilder {
	b := &SExprBuilder{stack: make([]*tokenSeq, 0)}
	return b
}

func (b *SExprBuilder) Len() int {
	return len(b.stack)
}

func (b *SExprBuilder) push(expr interface{}) {
	top := b.Len() - 1
	seq := b.stack[top]
	seq.items = append(seq.items, expr)
	b.stack[top] = seq
}

func (b *SExprBuilder) startSeq(typ int) {
	seq := new(tokenSeq)
	seq.typ = typ
	seq.items = make([]interface{}, 0)
	b.stack = append(b.stack, seq)
}

func (b *SExprBuilder) endSeq() *tokenSeq {
	last := b.Len() - 1
	seq := b.stack[last]
	b.stack = b.stack[0:last]
	return seq
}

func (builder *SExprBuilder) Run(reader *Reader) *tokenSeq {
	builder.startSeq(-1)
	for tk := reader.Read(); tk.id != TOKEN_ENDOFINPUT; tk = reader.Read() {
		switch tk.id {
		case TOKEN_LEFT_PAREN:
			builder.startSeq(tk.id)
		case TOKEN_RIGHT_PAREN:
			seq := builder.endSeq()
			if seq.typ != TOKEN_LEFT_PAREN {
				panic("PAREN does not match")
			}
			builder.push(List(seq.items...))
		case TOKEN_LEFT_BRACKET:
			builder.startSeq(tk.id)
		case TOKEN_RIGHT_BRACKET:
			seq := builder.endSeq()
			if seq.typ != TOKEN_LEFT_BRACKET {
				panic("PAREN does not match")
			}
			fmt.Println("seq.items", seq.items)
			builder.push(seq.items)
		case TOKEN_LEFT_BRACE:
			builder.startSeq(tk.id)
		case TOKEN_RIGHT_BRACE:
			seq := builder.endSeq()
			if seq.typ != TOKEN_LEFT_BRACE {
				panic("PAREN does not match")
			}
			m := make(map[interface{}]interface{})
			var key interface{}
			for i, v := range seq.items {
				println("got", i, v)
				if i%2 == 0 {
					key = v
				} else {
					println("putting", key, v)
					m[key] = v
				}
			}
			builder.push(m)
		case TOKEN_INT:
			var n int
			fmt.Sscanf(tk.v, "%d", &n)
			builder.push(n)
		case TOKEN_SYMBOL:
			builder.push(sym(tk.v))
		case TOKEN_STRING:
			builder.push(tk.v)
		}
	}
	seq := builder.endSeq()
	if seq.typ != -1 {
		panic("expected TopLevel")
	}
	return seq
}

func BuildSExpr(buf *Buffered) interface{} {
	reader := NewReader(buf)
	builder := NewSExprBuilder()
	seq := builder.Run(reader)
	return seq.items[0]
}
