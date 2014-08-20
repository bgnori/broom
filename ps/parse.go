package ps

import (
    "io"
)

func Parse(r io.Reader) []Value {
    tknz := NewTokenizer(r)
    ch := tknz.Tokenize()
    result := make([]Value, 0)
    for tk := <-ch ; tk.tkid != TOKEN_ENDOFINPUT; tk = <-ch {
        result = append(result, Expr(tk, ch))
    }
    return result
}

func Expr(tk *Token, ch chan *Token) Value {
    switch tk.tkid {
    case TOKEN_CHUNK:
        return tk.value
    case TOKEN_LEFTPAREN:
        xs := make([]Value, 0)
        for {
            tk = <-ch
            switch tk.tkid {
            case TOKEN_RIGHTPAREN:
                return MakeList(nil, xs...)
            case TOKEN_DOT:
                afterdot := <-ch
                if afterdot.tkid == TOKEN_DOT {
                    panic("got '.' after '.'")
                }
                cdr := Expr(afterdot, ch)
                r := MakeList(cdr, xs...)
                if (<-ch).tkid != TOKEN_RIGHTPAREN {
                    panic("expected ')' of ... '.' expr ')'")
                }
                return r
            case TOKEN_ENDOFINPUT:
                panic("unexpected end of steam, Parlen must be closed.")
            default:
                xs = append(xs, Expr(tk, ch))
            }
        }
    case TOKEN_RIGHTPAREN:
        panic("got ')' !")
    case TOKEN_DOT:
        panic("got '.' !")
    default:
        println(tk.tkid)
        println(tk.start, tk.end)
        panic("unknown token id")
    }
}

