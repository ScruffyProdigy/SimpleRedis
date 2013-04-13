package redis

import (
	"bytes"
	"errors"
	"io"
	"strings"

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
	bufferSize  = 200
)

var (
	delimiter = []byte{'\r', '\n'}
)

type command interface {
	arguments() []string
	callback() func(*response) error
}

type Executor interface {
	Execute(command) error
	ErrCallback(error, string)
}

func (this Client) Execute(command command) error {
	go this.useConnection(func(conn *Connection) {
		err := conn.Execute(command)
		if err != nil {
			// we are in a separate routine and cannot return the error
			// use the callback instead
			this.errCallback.Call(err, strings.Join(command.arguments(), " "))
		}
	})
	return nil
}

func (this Client) ErrCallback(e error, s string) {
	this.errCallback.Call(e, s)
}

func (this Connection) input(command command) error {
	comm, err := buildCommand(command.arguments())
	if err != nil {
		return err
	}

	_, err = this.Write(comm)
	return err
}

func (this Connection) output(command command) error {
	res, err := getResponse(this)
	if err != nil {
		return err
	}

	return command.callback()(res)
}

func (this Connection) Execute(command command) error {
	err := this.input(command)
	if err != nil {
		return err
	}

	err = this.output(command)
	if err != nil {
		return err
	}

	return nil
}

func buildCommand(arguments []string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	var err error

	if err = buf.WriteByte(isMultibulk); err != nil {
		return nil, err
	}
	if _, err = buf.WriteString(itoa(len(arguments))); err != nil {
		return nil, err
	}
	if _, err = buf.Write(delimiter); err != nil {
		return nil, err
	}

	for _, arg := range arguments {
		if err = buf.WriteByte(isBulk); err != nil {
			return nil, err
		}
		if _, err = buf.WriteString(itoa(len(arg))); err != nil {
			return nil, err
		}
		if _, err = buf.Write(delimiter); err != nil {
			return nil, err
		}
		if _, err = buf.WriteString(arg); err != nil {
			return nil, err
		}
		if _, err = buf.Write(delimiter); err != nil {
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

/*

BoolCommand - the command type used when a boolean response is expected

*/

type BoolCommand struct {
	args   []string
	output chan<- bool
}

func newBoolCommand(args []string) (command, <-chan bool) {
	c := make(chan bool, 1)
	return BoolCommand{args, c}, c
}

func (this BoolCommand) arguments() []string {
	return this.args
}

func (this BoolCommand) callback() func(*response) error {
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

type IntCommand struct {
	args   []string
	output chan<- int
}

func newIntCommand(args []string) (command, <-chan int) {
	c := make(chan int, 1)
	return IntCommand{args, c}, c
}

func (this IntCommand) arguments() []string {
	return this.args
}

func (this IntCommand) callback() func(*response) error {
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

type FloatCommand struct {
	args   []string
	output chan<- float64
}

func newFloatCommand(args []string) (command, <-chan float64) {
	c := make(chan float64, 1)
	return FloatCommand{args, c}, c
}

func (this FloatCommand) arguments() []string {
	return this.args
}

func (this FloatCommand) callback() func(*response) error {
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

type StringCommand struct {
	args   []string
	output chan<- string
}

func newStringCommand(args []string) (command, <-chan string) {
	c := make(chan string, 1)
	return StringCommand{args, c}, c
}

func (this StringCommand) arguments() []string {
	return this.args
}

func (this StringCommand) callback() func(*response) error {
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

type SliceCommand struct {
	args   []string
	output chan<- []string
}

func newSliceCommand(args []string) (command, <-chan []string) {
	c := make(chan []string, 1)
	return SliceCommand{args, c}, c
}

func (this SliceCommand) arguments() []string {
	return this.args
}

func (this SliceCommand) callback() func(*response) error {
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

type MaybeSliceCommand struct {
	args   []string
	output chan<- []*string
}

func newMaybeSliceCommand(args []string) (command, <-chan []*string) {
	c := make(chan []*string, 1)
	return MaybeSliceCommand{args, c}, c
}

func (this MaybeSliceCommand) arguments() []string {
	return this.args
}

func (this MaybeSliceCommand) callback() func(*response) error {
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

type MapCommand struct {
	args   []string
	output chan<- map[string]string
}

func newMapCommand(args []string) (command, <-chan map[string]string) {
	c := make(chan map[string]string, 1)
	return MapCommand{args, c}, c
}

func (this MapCommand) arguments() []string {
	return this.args
}

func (this MapCommand) callback() func(*response) error {
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

type NilCommand struct {
	args   []string
	output chan<- nothing
}

func newNilCommand(args []string) (command, <-chan nothing) {
	c := make(chan nothing, 1)
	return NilCommand{args, c}, c
}

func (this NilCommand) arguments() []string {
	return this.args
}

func (this NilCommand) callback() func(*response) error {
	return func(r *response) error {
		defer close(this.output)
		this.output <- nothing{}
		return nil
	}
}
