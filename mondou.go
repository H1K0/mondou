package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var ENV struct {
	namespace map[string]interface{}
	stack []interface{}
	stack_ptr int
}

// Evaluate code
func Eval(code string) (res interface{}, err error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}
	// parsing quotes
	quotes := []int{-1, -1}
	for i, char := range code {
		if len(quotes) % 2 == 0 && i != 0 && code[i-1:i+1] == "//" {
			code = code[:i-1]
			break
		}
		if char == '"' && (i == 0 || code[i-1] != '\\' || len(quotes) % 2 == 0) {
			quotes = append(quotes, i)
		}
	}
	if len(quotes) % 2 == 1 {
		err = errors.New("inconsistent quotes")
		return
	}
	init_stack_ptr := ENV.stack_ptr
	defer func() { ENV.stack = ENV.stack[:init_stack_ptr]; ENV.stack_ptr = init_stack_ptr }()
	code_new := ""
	var tmp interface{}
	for i := 2; i < len(quotes); i += 2 {
		tmp, err = strconv.Unquote(code[quotes[i]:quotes[i+1]+1])
		if err != nil {
			err = errors.New("invalid string")
			return
		}
		ENV.stack = append(ENV.stack, tmp)
		code_new += fmt.Sprintf("%s $%d ", code[quotes[i-1]+1:quotes[i]], ENV.stack_ptr)
		ENV.stack_ptr++
	}
	code = code_new + code[quotes[len(quotes)-1]+1:]
	// parsing brackets
	s := -1
	e := -1
	k := 0
	code_new = ""
	for i, char := range code {
		switch char {
		case '(':
			if k == 0 {
				s = i
			}
			k++
		case ')':
			k--
			if k == 0 {
				tmp, err = Eval(code[s+1:i])
				if err != nil {
					return
				}
				ENV.stack = append(ENV.stack, tmp)
				code_new += fmt.Sprintf("%s $%d ", code[e+1:s], ENV.stack_ptr)
				ENV.stack_ptr++
				e = i
			} else if k < 0 {
				break
			}
		}
	}
	if k != 0 {
		err = errors.New("inconsistent brackets")
		return
	}
	code = code_new + code[e+1:]
	// parsing literals and names
	var expr []interface{}
	literals := strings.Fields(code)
	for _, lit := range literals {
		if lit[0] == '$' {
			var stk_ptr int
			stk_ptr, err = strconv.Atoi(lit[1:])
			if err != nil {
				err = errors.New("unexpected error while parsing")
				return
			}
			expr = append(expr, ENV.stack[stk_ptr])
			continue
		}
		if val, err := strconv.Atoi(lit); err == nil {
			expr = append(expr, val)
			continue
		}
		if val, err := strconv.ParseFloat(lit, 64); err == nil {
			expr = append(expr, val)
			continue
		}
		if lit == "+" || lit == "-" || lit == "*" || lit == "/" {
			l := len(expr)
			if l < 2 {
				err = errors.New(fmt.Sprintf("not enough operands for operation %s", lit))
				return
			}
			switch lit {
			case "+":
				tmp, err = Add(expr[l-2], expr[l-1])
			case "-":
				tmp, err = Substract(expr[l-2], expr[l-1])
			case "*":
				tmp, err = Multiply(expr[l-2], expr[l-1])
			case "/":
				tmp, err = Divide(expr[l-2], expr[l-1])
			}
			if err != nil {
				return
			}
			expr[l-2] = tmp
			expr = expr[:l-1]
			continue
		}
		if val, ok := ENV.namespace[lit]; ok {
			expr = append(expr, val)
		} else {
			err = errors.New(fmt.Sprintf("unknown name '%s'", lit))
			return
		}
	}
	res, err = expr[0], nil
	return
}

// Execute code
func Exec(code string) (res string, err error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}
	switch code[0] {
	case ':':
		res, err = LoadFromFile(code[1:])
		return
	case '!':
		res, err = SetVar(code[1:])
		return
	case '<':
		err = Print(code[1:])
		return
	case '>':
		res, err = ReadVar(code[1:])
		return
	case '@':
		res, err = DefFunc(code[1:])
		return
	}
	var tmp interface{}
	if strings.HasPrefix(code, "typeof") {
		tmp, err = Eval(code[6:])
		if err == nil {
			res = fmt.Sprintf("%T", tmp)
		}
		return
	}
	tmp, err = Eval(code)
	res = fmt.Sprintf("%#v", tmp)
	return
}

func main() {
	var err error
	var input string
	var res string
	fmt.Println("(c) Masahiko AMANO a.k.a. H1K0, 2024-present\n[問答] Mondou REPL\n")
	ENV.namespace = make(map[string]interface{})
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("mondou$ ")
		input, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				break
			} else {
				panic(err)
			}
		}
		res, err = Exec(input)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		} else if res != "" {
			fmt.Println(res)
		}
	}
}
