package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Load and execute code from file
func LoadFromFile(path string) (res string, err error) {
	comment := strings.Index(path, "//")
	if comment == -1 {
		comment = len(path)
	}
	path = strings.TrimSpace(path)
	file, err := os.Open(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("could not load file '%s'", path))
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	linecount := 1
	for {
		code, err1 := reader.ReadString('\n')
		res, err = Exec(code)
		if err != nil {
			err = errors.New(fmt.Sprintf("'%s'@%d: %s", path, linecount, err))
			return
		}
		if err1 == io.EOF {
			break
		}
		linecount++
	}
	return
}

// Set variable
func SetVar(def string) (res string, err error) {
	def = strings.TrimSpace(def)
	sep := strings.Index(def, " ")
	var_name := def[:sep]
	if strings.Contains(var_name, ":") ||
	   strings.Contains(var_name, "!") ||
	   strings.Contains(var_name, "@") ||
	   strings.Contains(var_name, "$") ||
	   strings.Contains(var_name, "<") ||
	   strings.Contains(var_name, ">") {
		err = errors.New(fmt.Sprintf("invalid variable name: '%s'", var_name))
		return
	}
	var_value, err := Eval(def[sep:])
	if err != nil {
		return
	}
	ENV.namespace[var_name] = var_value
	return
}

// Print value
func Print(code string) (err error) {
	res, err := Eval(code)
	if err != nil {
		return
	}
	fmt.Printf("%v\n", res)
	return
}

// Read variable from console
func ReadVar(def string) (res string, err error) {
	var_name := strings.TrimSpace(def)
	if _, ok := ENV.namespace[var_name]; !ok {
		err = errors.New(fmt.Sprintf("undefined variable: '%s'", var_name))
		return
	}
	reader := bufio.NewReader(os.Stdin)
	var_value, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return
	}
	ENV.namespace[var_name] = var_value[:len(var_value)-1]
	return
}

// Define function
func DefFunc(def string) (res string, err error) {
	err = errors.New("function definition is not supported yet")
	return
}

// Add values
func Add(lhs, rhs interface{}) (res interface{}, err error) {
	ltype := fmt.Sprintf("%T", lhs)
	rtype := fmt.Sprintf("%T", rhs)
	if ltype == "int" && rtype == "float64" {
		res = float64(lhs.(int)) + float64(rhs.(float64))
		return
	}
	if ltype == "float64" && rtype == "int" {
		res = float64(lhs.(float64)) + float64(rhs.(int))
		return
	}
	if ltype != rtype {
		err = errors.New(fmt.Sprintf("unsupported operand type(s) for +: '%s' and '%s'", ltype, rtype))
		return
	}
	if ltype == "int" {
		res = lhs.(int) + rhs.(int)
		return
	}
	if ltype == "string" {
		res = lhs.(string) + rhs.(string)
		return
	}
	return
}

// Substract values
func Substract(lhs, rhs interface{}) (res interface{}, err error) {
	ltype := fmt.Sprintf("%T", lhs)
	rtype := fmt.Sprintf("%T", rhs)
	if ltype == "int" && rtype == "float64" {
		res = float64(lhs.(int)) - float64(rhs.(float64))
		return
	}
	if ltype == "float64" && rtype == "int" {
		res = float64(lhs.(float64)) - float64(rhs.(int))
		return
	}
	if ltype != rtype {
		err = errors.New(fmt.Sprintf("unsupported operand type(s) for -: '%s' and '%s'", ltype, rtype))
		return
	}
	if ltype == "int" {
		res = lhs.(int) - rhs.(int)
		return
	}
	return
}

// Multiply values
func Multiply(lhs, rhs interface{}) (res interface{}, err error) {
	ltype := fmt.Sprintf("%T", lhs)
	rtype := fmt.Sprintf("%T", rhs)
	if ltype == "int" && rtype == "float64" {
		res = float64(lhs.(int)) * float64(rhs.(float64))
		return
	}
	if ltype == "float64" && rtype == "int" {
		res = float64(lhs.(float64)) * float64(rhs.(int))
		return
	}
	if ltype == "string" && rtype == "int" {
		res = strings.Repeat(lhs.(string), rhs.(int))
		return
	}
	if ltype != rtype {
		err = errors.New(fmt.Sprintf("unsupported operand type(s) for *: '%s' and '%s'", ltype, rtype))
		return
	}
	if ltype == "int" {
		res = lhs.(int) * rhs.(int)
		return
	}
	return
}

// Divide values
func Divide(lhs, rhs interface{}) (res interface{}, err error) {
	ltype := fmt.Sprintf("%T", lhs)
	rtype := fmt.Sprintf("%T", rhs)
	if (rtype == "float64" || rtype == "int") && rhs == 0 {
		err = errors.New("division by zero")
		return
	}
	if ltype == "int" && rtype == "float64" {
		res = float64(lhs.(int)) / float64(rhs.(float64))
		return
	}
	if ltype == "float64" && rtype == "int" {
		res = float64(lhs.(float64)) / float64(rhs.(int))
		return
	}
	if ltype != rtype {
		err = errors.New(fmt.Sprintf("unsupported operand type(s) for /: '%s' and '%s'", ltype, rtype))
		return
	}
	if ltype == "int" {
		res = lhs.(int) / rhs.(int)
		return
	}
	return
}
