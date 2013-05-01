package redis

import (
	"bytes"
	"errors"
	"io"

//	"bufio"
)

//a response either has a value, or a list of subresponses (which themselves usually have values, but occasionally subresponses)
type response struct {
	val          string
	subresponses []*response
}

const (
	isMultibulk = '*'
	isBulk      = '$'
	isInt       = ':'
	isStatus    = '+'
	isError     = '-'
	bufferSize  = 256
)

var (
	delimiter = []byte{'\r', '\n'}
)

type command interface {
	arguments() []string
	callback() func(*response) error
}

//Anything that can execute a command is an Executor
type Executor interface {
	Execute(command)
}

//Anything that can execute a command, and can deal with resulting errors is a SafeExecutor
type SafeExecutor interface {
	Executor
	errCallback(error, string)
}

func buildCommand(arguments []string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if err := buf.WriteByte(isMultibulk); err != nil {
		return nil, err
	}
	if _, err := buf.WriteString(itoa(len(arguments))); err != nil {
		return nil, err
	}
	if _, err := buf.Write(delimiter); err != nil {
		return nil, err
	}

	for _, arg := range arguments {
		if err := buf.WriteByte(isBulk); err != nil {
			return nil, err
		}
		if _, err := buf.WriteString(itoa(len(arg))); err != nil {
			return nil, err
		}
		if _, err := buf.Write(delimiter); err != nil {
			return nil, err
		}
		if _, err := buf.WriteString(arg); err != nil {
			return nil, err
		}
		if _, err := buf.Write(delimiter); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func getResponse(conn io.Reader) (*response, error) {
	var buffer [1]byte
	_, err := conn.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	switch buffer[0] {
	case isError:
		errString, err := getString(conn)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(errString)
	case isStatus, isInt:
		return getStringResponse(conn)
	case isBulk:
		return getBulk(conn)
	case isMultibulk:
		return getMultiBulk(conn)
	default:
		return nil, errors.New("Unknown Data Type:'" + string(buffer[0:1]) + "'")
	}
	return nil, errors.New("unexpected data")
}

func getString(conn io.Reader) (string, error) {
	var buffer [bufferSize]byte
	j := -len(delimiter)
	i := 0
	for {
		if j >= 0 && bytes.Equal(buffer[j:i], delimiter) {
			return string(buffer[:j]), nil
		}
		if i >= bufferSize {
			return "", errors.New("Short Buffer - " + string(buffer[:]))
		}
		conn.Read(buffer[i : i+1])
		i++
		j++
	}
	return string(buffer[:]), nil
}

func getStringResponse(conn io.Reader) (*response, error) {
	val, err := getString(conn)
	if err != nil {
		return nil, err
	}
	return &response{
		val: val,
	}, nil
}

func getBulk(conn io.Reader) (*response, error) {
	line, err := getString(conn)
	if err != nil {
		return nil, err
	}

	strlen, err := atoi(line)
	if err != nil {
		return nil, err
	}
	if strlen == -1 {
		return nil, nil
	}

	b := make([]byte, strlen+len(delimiter))
	i, err := conn.Read(b)
	if err != nil {
		//the read should be successful
		return nil, err
	}
	if i != strlen+len(delimiter) {
		//the read should go through every byte we have set out
		return nil, errors.New("underread")
	}
	if !bytes.Equal(b[strlen:], delimiter) {
		//the read should end with a crlf
		return nil, errors.New("Incorrect Redis bulk length")
	}

	return &response{
		val: string(b[:strlen]),
	}, nil
}

func getMultiBulk(conn io.Reader) (*response, error) {
	line, err := getString(conn)
	if err != nil {
		return nil, err
	}

	cResponses, err := atoi(string(line))
	if err != nil {
		return nil, err
	}
	if cResponses == -1 {
		return nil, nil
	}

	r := new(response)
	r.subresponses = make([]*response, cResponses)

	for iResponse := 0; iResponse < int(cResponses); iResponse++ {
		var err error
		r.subresponses[iResponse], err = getResponse(conn)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

//TODO: change the slice argument for each Command Function to become a variable length parameter

/*

BoolCommand - the command type used when a boolean response is expected

*/

type boolCommand struct {
	args   []string
	output chan<- bool
}

//BoolCommand executes the command specified by the arguments specified
//It returns the response Redis generates coerced into a boolean value
func BoolCommand(e Executor, args []string) <-chan bool {
	c := make(chan bool, 1)
	e.Execute(boolCommand{args, c})
	return c
}

func (this boolCommand) arguments() []string {
	return this.args
}

func (this boolCommand) callback() func(*response) error {
	return func(r *response) error {
		defer close(this.output)
		if r != nil {
			this.output <- r.val == "1"
		}
		return nil
	}
}

/*

IntCommand - the command type used when an int response is expected

*/

type intCommand struct {
	args   []string
	output chan<- int
}

//IntCommand executes the command specified by the arguments specified
//It returns the response Redis generates coerced into an integer value
func IntCommand(e Executor, args []string) <-chan int {
	c := make(chan int, 1)
	e.Execute(intCommand{args, c})
	return c
}

func (this intCommand) arguments() []string {
	return this.args
}

func (this intCommand) callback() func(*response) error {
	return func(r *response) error {
		defer close(this.output)
		if r != nil {
			res, err := atoi(r.val)
			if err != nil {
				return err
			}
			this.output <- res
		}
		return nil
	}
}

/*

FloatCommand - the command type used when a float response is expected

*/

type floatCommand struct {
	args   []string
	output chan<- float64
}

//FloatCommand executes the command specified by the arguments specified
//It returns the response Redis generates coerced into a float value
func FloatCommand(e Executor, args []string) <-chan float64 {
	c := make(chan float64, 1)
	e.Execute(floatCommand{args, c})
	return c
}

func (this floatCommand) arguments() []string {
	return this.args
}

func (this floatCommand) callback() func(*response) error {
	return func(r *response) error {
		defer close(this.output)
		if r != nil {
			f, err := atof(r.val)
			if err != nil {
				return err
			}
			this.output <- f
		}
		return nil
	}
}

/*

StringCommand - the command type used when a string response is expected

*/

type stringCommand struct {
	args   []string
	output chan<- string
}

//StringCommand executes the command specified by the arguments specified
//It returns the response Redis generates coerced into a string value
func StringCommand(e Executor, args []string) <-chan string {
	c := make(chan string, 1)
	e.Execute(stringCommand{args, c})
	return c
}

func (this stringCommand) arguments() []string {
	return this.args
}

func (this stringCommand) callback() func(*response) error {
	return func(r *response) error {
		defer close(this.output)

		if r != nil {
			this.output <- r.val
		}
		return nil
	}
}

/*

SliceCommand - the command type used when a []string response is expected

*/

type sliceCommand struct {
	args   []string
	output chan<- []string
}

//SliceCommand executes the command specified by the arguments specified
//It returns the response Redis generates coerced into a slice/array
func SliceCommand(e Executor, args []string) <-chan []string {
	c := make(chan []string, 1)
	e.Execute(sliceCommand{args, c})
	return c
}

func (this sliceCommand) arguments() []string {
	return this.args
}

func (this sliceCommand) callback() func(*response) error {
	return func(r *response) error {
		defer close(this.output)

		if r != nil {
			actualResponse := make([]string, len(r.subresponses))
			for i, line := range r.subresponses {
				if line != nil {
					actualResponse[i] = line.val
					//warning, if it gets in here, could causes unintended problems
				}
			}

			this.output <- actualResponse
		}

		return nil
	}
}

/*

MaybeSliceCommand - the command type used when a []string response would normally be expected, but there's a chance that some of the strings won't be there
used for SORT command, when GET could be used on a key that doesn't have something there.  
SliceCommand will return an empty string, but if you need to differentiate between nonexistant and empty key values, this is necessary

*/

type maybeSliceCommand struct {
	args   []string
	output chan<- []*string
}

//MaybeSliceCommand executes the command specified by the arguments specified
//It returns the response Redis generates coerced into a slice of pointers
func MaybeSliceCommand(e Executor, args []string) <-chan []*string {
	c := make(chan []*string, 1)
	e.Execute(maybeSliceCommand{args, c})
	return c
}

func (this maybeSliceCommand) arguments() []string {
	return this.args
}

func (this maybeSliceCommand) callback() func(*response) error {
	return func(r *response) error {
		defer close(this.output)
		if r != nil {
			actualResponse := make([]*string, len(r.subresponses))
			for i, line := range r.subresponses {
				if line != nil {
					copy := line.val
					actualResponse[i] = &copy
				}
			}

			this.output <- actualResponse
		}
		return nil
	}
}

/*

MapCommand - the command type used when a map[string]string response is expected

*/

type mapCommand struct {
	args   []string
	output chan<- map[string]string
}

//BoolCommand executes the command specified by the arguments specified
//It returns the response Redis generates coerced into a map
func MapCommand(e Executor, args []string) <-chan map[string]string {
	c := make(chan map[string]string, 1)
	e.Execute(mapCommand{args, c})
	return c
}

func (this mapCommand) arguments() []string {
	return this.args
}

func (this mapCommand) callback() func(*response) error {
	return func(r *response) error {
		defer close(this.output)
		if r != nil {
			m := make(map[string]string, len(r.subresponses)/2)
			for i := 0; i+1 < len(r.subresponses); i += 2 {
				if r.subresponses[i] != nil && r.subresponses[i+1] != nil {
					m[r.subresponses[i].val] = r.subresponses[i+1].val
				}
			}
			this.output <- m
		}
		return nil
	}
}

/*

NilCommand - the command type used when no response is expected

*/

type nothing struct {
}

type nilCommand struct {
	args   []string
	output chan<- nothing
}

//NilCommand executes the command specified by the arguments specified
//It does not return a usable value
func NilCommand(e Executor, args []string) <-chan nothing {
	c := make(chan nothing, 1)
	e.Execute(nilCommand{args, c})
	return c
}

func (this nilCommand) arguments() []string {
	return this.args
}

func (this nilCommand) callback() func(*response) error {
	return func(r *response) error {
		defer close(this.output)
		if r != nil {
			this.output <- nothing{}
		}
		return nil
	}
}
