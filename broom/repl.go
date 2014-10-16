package broom

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Load(file io.Reader, env Environment, verbose bool) error {
	buf := NewBuffered(file)
	reader := NewReader(buf)
	builder := NewSExprBuilder()
	prg, err := builder.Run(reader)
	if err != nil {
		return err
	}
	for _, expr := range prg.Items() {
		if verbose {
			fmt.Println("in:", expr)
		}
		got := Eval(env, expr)
		if verbose {
			fmt.Println("-->", got)
		}
	}
	return nil
}

func Repl(in io.Reader, env Environment) {
	env.Bind(sym("dump"), func(dynamic Environment, cdr List) interface{} {
		dynamic.Dump()
		return nil
	})

	reader := bufio.NewReader(in)

	var chunks []string
	chunks = nil

	fmt.Println("Hello!")
	fmt.Print("broom > ")

	for line, _, err := reader.ReadLine(); err != io.EOF; line, _, err = reader.ReadLine() {
		if err != nil {
			panic(err)
		}
		if len(line) == 0 {
			input := strings.Join(chunks, "")
			if len(input) == 0 {
				fmt.Println("no input...")
				fmt.Print("broom > ")
				continue
			}
			expr, err := try2Build(input)
			chunks = nil
			if err != nil {
				fmt.Println("Something wrong with input!")
				fmt.Println(err)
				fmt.Print("broom > ")
				continue
			}
			fmt.Println("input:", expr)

			got, err := try2Eval(env, expr)
			if err != nil {
				fmt.Println("Failed eval!")
				fmt.Println(err)
				fmt.Print("broom > ")
				continue
			}
			fmt.Println("-->", got)
			fmt.Print("broom > ")
		} else {
			fmt.Print("... ")
			if chunks == nil {
				chunks = make([]string, 1)
			}
			chunks = append(chunks, string(line))
		}
	}
	fmt.Println()
	fmt.Println("bye!")
}

type MyErr string

func (e MyErr) Error() string {
	return string(e)
}

func try2Build(c string) (expr interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			expr = nil
			if ee, ok := e.(EvalError); ok {
				err = MyErr(ee.Error())
			} else {
				err = MyErr(e.(string))
			}
		}
	}()
	return BuildSExpr(NewBuffered(strings.NewReader(c))), nil
}

func try2Eval(env Environment, expr interface{}) (result interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			result = nil
			if ee, ok := e.(EvalError); ok {
				err = MyErr(ee.Error())
			} else {
				err = MyErr(e.(string))
			}
		}
	}()
	return Eval(env, expr), nil
}
