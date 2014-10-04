package broom

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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

func isUnQuote(r rune) bool {
	return ',' == r
}

func isSplicingQuote(r rune) bool {
	return '@' == r
}

func isGenSymQuote(r rune) bool {
	return '~' == r
}

func isWhite(r rune) bool {
	return r == ' ' || r == '\t' || r == '\v' || r == '\r' || r == '\n'
}

const (
	MODE_NONE = iota
	MODE_CR
	MODE_LF
	MODE_CRLF
)

func isEndOfLine(r1, r2 rune) (int, int) {
	if r1 == '\r' && r2 == '\n' {
		return 2, MODE_CRLF
	}
	if r1 == '\r' && r2 != '\n' {
		return 1, MODE_CR
	}
	if r1 == '\n' {
		return 1, MODE_LF
	}
	return 0, MODE_NONE
}

func isSemicolon(r rune) bool {
	return r == ';'
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
		if token, err := reader.tryString(); err != nil {
			panic(err)
		} else {
			reader.Emit(token)
		}
		return TopLevel
	case isVerticalBar(r):
		reader.Emit(reader.MakeVerticalVar())
		return TopLevel
	case isSharp(r):
		reader.Emit(reader.MakeSharp())
		return TopLevel
	case isQuote(r):
		reader.Emit(reader.MakeQuote())
		return TopLevel
	case isQuasiQuote(r):
		reader.Emit(reader.MakeQuasiQuote())
		return TopLevel
	case isUnQuote(r):
		reader.Emit(reader.MakeUnquote())
		return TopLevel
	case isSplicingQuote(r):
		reader.Emit(reader.MakeSplicingQuote())
		return TopLevel
	case isSemicolon(r):
		reader.Emit(reader.MakeSemicolon())
		return ZapToLineEnd
	default:
		reader.Emit(reader.tryChunk())
		return TopLevel
	}
	panic("Should not reach here")
}

func ZapToLineEnd(reader *Reader) ReaderState {
	r, eos := reader.buffer.Peek()
	if eos {
		reader.Emit(Token{id: TOKEN_ENDOFINPUT})
		return nil
	}
	q := r
	reader.buffer.Consume(1)
	r, eos = reader.buffer.Peek()
	if eos {
		reader.Emit(Token{id: TOKEN_ENDOFINPUT})
		return nil
	}
	for {
		_, mode := isEndOfLine(q, r)
		switch mode {
		case MODE_NONE:
			q = r
			reader.buffer.Consume(1)
			r, eos = reader.buffer.Peek()
			if eos {
				reader.Emit(Token{id: TOKEN_ENDOFINPUT})
				return nil
			}
		case MODE_CR:
			reader.Emit(Token{pos: reader.buffer.pos, id: TOKEN_ENDOFLINE})
			return TopLevel
		case MODE_LF:
			reader.Emit(Token{pos: reader.buffer.pos, id: TOKEN_ENDOFLINE})
			return TopLevel
		case MODE_CRLF:
			reader.Emit(Token{pos: reader.buffer.pos, id: TOKEN_ENDOFLINE})
			reader.buffer.Consume(1)
			return TopLevel
		}
	}
	panic("Never reach")
}

func (reader *Reader) tryString() (Token, error) {
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
		return Token{id: TOKEN_ERROR}, &BuilderError{"string must be closed"}
	}
	if isDoubleQuote(r) {
		reader.buffer.Consume(1) // skip "
	}
	return Token{id: TOKEN_STRING, v: string(xs), pos: pos}, nil
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
	if s == "true" {
		return Token{id: TOKEN_TRUE, v: s, pos: pos}
	}
	if s == "false" {
		return Token{id: TOKEN_FALSE, v: s, pos: pos}
	}

	if strings.Contains(s, ".") {
		var f float64
		if _, err := fmt.Sscanf(s, "%f", &f); err == nil {
			return Token{id: TOKEN_FLOAT, v: s, pos: pos}
		}
	}
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
	TOKEN_ENDOFLINE
	TOKEN_SYMBOL
	TOKEN_INT
	TOKEN_FLOAT
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
	TOKEN_UNQUOTE
	TOKEN_SPLICINGQUOTE
	TOKEN_GENSYMQUOTE
	TOKEN_DOT
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_ERROR
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

func (r *Reader) MakeUnquote() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_UNQUOTE}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeSplicingQuote() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_SPLICINGQUOTE}
	r.buffer.Consume(1)
	return t
}

func (r *Reader) MakeGenSymQuote() Token {
	t := Token{pos: r.buffer.pos, id: TOKEN_GENSYMQUOTE}
	r.buffer.Consume(1)
	return t
}


type tokenSeq struct {
	typ    int
	items  []interface{}
	deco *Decorator
}

func NewTokenSeq(typ int, deco *Decorator) *tokenSeq {
	seq := new(tokenSeq)
	seq.typ = typ

	seq.items = make([]interface{}, 0)
	seq.deco = deco
	return seq
}

type SeqError struct {
	msg string
}

func (err *SeqError) Error() string {
	return err.msg
}

func (t *tokenSeq) Error(msg string) *SeqError {
	return &SeqError{msg: msg}
}

func (t *tokenSeq) Items() []interface{} {
	return t.items
}

func (t *tokenSeq) Append(expr interface{}) {
	t.items = append(t.items, expr)
}

func (t *tokenSeq) CloseParen() (interface{}, error) {
	if t.typ != TOKEN_LEFT_PAREN {
		return nil, t.Error("PAREN does not match")
	}
	return t.deco.Apply(List(t.Items()...)), nil
}

func (t *tokenSeq) CloseBracket() (interface{}, error) {
	if t.typ != TOKEN_LEFT_BRACKET {
		return nil, t.Error("PAREN does not match")
	}
	return t.deco.Apply(t.Items()), nil
}

func (t *tokenSeq) CloseBrace() (interface{}, error) {
	if t.typ != TOKEN_LEFT_BRACE {
		return nil, t.Error("PAREN does not match")
	}
	var key interface{}
	m := make(map[interface{}]interface{})
	for i, v := range t.items {
		if i%2 == 0 {
			key = v
		} else {
			m[key] = v
		}
	}
	return t.deco.Apply(m), nil
}


type Decorator struct {
	stack []Symbol
}

func NewDecorator() *Decorator {
	return &Decorator{stack:make([]Symbol, 0)}
}

func (d *Decorator) Push(s Symbol) {
	d.stack = append(d.stack, s)
}

func (d *Decorator) HasSomething() bool {
	return len(d.stack) > 0
}

func (d *Decorator) Pop() Symbol {
	idx := len(d.stack) - 1
	last := d.stack[idx]
	d.stack = d.stack[0:idx]
	return last
}

func (d *Decorator) Apply(expr interface{}) interface{} {
	for ; d.HasSomething() ; {
		s := d.Pop()
		expr = List(s, expr)
	}
	return expr
}


type SExprBuilder struct {
	stack         []*tokenSeq
	deco *Decorator
}

func NewSExprBuilder() *SExprBuilder {
	b := &SExprBuilder{stack: make([]*tokenSeq, 0)}
	b.ResetDeco()
	return b
}

func (b *SExprBuilder)ResetDeco() {
	b.deco = NewDecorator()
}

func (b *SExprBuilder) Len() int {
	return len(b.stack)
}

func (b *SExprBuilder) Top() *tokenSeq {
	return b.stack[b.Len() - 1]
}

func (b *SExprBuilder) Pop() *tokenSeq {
	last := b.Len() - 1
	seq := b.stack[last]
	b.stack = b.stack[0:last]
	return  seq
}


func (b *SExprBuilder) Add(expr interface{}) {
	b.Top().Append(b.deco.Apply(expr))
}

func (b *SExprBuilder) startSeq(typ int) {
	seq := NewTokenSeq(typ, b.deco)
	b.ResetDeco()
	b.stack = append(b.stack, seq)
}

func (b *SExprBuilder) endSeq() *tokenSeq {
	return b.Pop()
}

type BuilderError struct {
	msg string
}

func (err *BuilderError) Error() string {
	return err.msg
}

func (b *SExprBuilder) Error(msg string) *BuilderError {
	return &BuilderError{msg: msg}
}

func (b *SExprBuilder) MakeParenObject(seq *tokenSeq) (interface{}, error) {
	if b.deco.HasSomething() {
		return nil, b.Error("can't decorate ')'.")
	}
	return seq.CloseParen()
}

func (b *SExprBuilder) MakeBracketObject(seq *tokenSeq) (interface{}, error) {
	if b.deco.HasSomething() {
		return nil, b.Error("can't decorate ']'.")
	}
	return seq.CloseBracket()
}

func (b *SExprBuilder) MakeBraceObject(seq *tokenSeq) (interface{}, error) {
	if b.deco.HasSomething() {
		return nil, b.Error("can't decorate '}'.")
	}
	return seq.CloseBrace()
}

func (builder *SExprBuilder) Run(reader *Reader) (*tokenSeq, error) {
	builder.startSeq(-1)
	for tk := reader.Read(); tk.id != TOKEN_ENDOFINPUT; tk = reader.Read() {
		switch tk.id {
		case TOKEN_QUOTE:
			builder.deco.Push(sym("quote"))
		case TOKEN_QUASIQUOTE:
			builder.deco.Push(sym("qq"))
		case TOKEN_UNQUOTE:
			builder.deco.Push(sym("uq"))
		case TOKEN_SPLICINGQUOTE:
			builder.deco.Push(sym("sq"))
		case TOKEN_LEFT_PAREN:
			builder.startSeq(tk.id)
		case TOKEN_RIGHT_PAREN:
			obj, err := builder.MakeParenObject(builder.endSeq())
			if err != nil {
				return nil, err
			}
			builder.Add(obj)
		case TOKEN_LEFT_BRACKET:
			builder.startSeq(tk.id)
		case TOKEN_RIGHT_BRACKET:
			obj, err := builder.MakeBracketObject(builder.endSeq())
			if err != nil {
				return nil, err
			}
			builder.Add(obj)
		case TOKEN_LEFT_BRACE:
			builder.startSeq(tk.id)
		case TOKEN_RIGHT_BRACE:
			obj, err := builder.MakeBraceObject(builder.endSeq())
			if err != nil {
				return nil, err
			}
			builder.Add(obj)
		case TOKEN_INT:
			var n int
			fmt.Sscanf(tk.v, "%d", &n)
			builder.Add(n)
		case TOKEN_FLOAT:
			var f float64
			fmt.Sscanf(tk.v, "%f", &f)
			builder.Add(f)
		case TOKEN_TRUE:
			builder.Add(true)
		case TOKEN_FALSE:
			builder.Add(false)
		case TOKEN_SYMBOL:
			builder.Add(sym(tk.v))
		case TOKEN_STRING:
			builder.Add(tk.v)
		}
	}
	seq := builder.endSeq()
	if seq.typ != -1 {
		return nil, builder.Error("expected TopLevel")
	}
	return seq, nil
}

func BuildSExpr(buf *Buffered) interface{} {
	reader := NewReader(buf)
	builder := NewSExprBuilder()
	seq, _ := builder.Run(reader)
	return seq.items[0]
}
