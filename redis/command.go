package redis

import (
	"bytes"
	"io"

//	"bufio"
)

//a response either has a value, or a list of subresponses (which themselves usually have values, but occasionally subresponses)
type response struct {
	val          string
	subresponses []*response
}

const (
	isMultibulk    = '*'
	isBulk         = '$'
	isInt          = ':'
	isStatus       = '+'
	isError        = '-'
	bufferSize     = 16 //should be enough to read any line that we don't know the size of
	longBufferSize = 200
)

var (
	delimiter = []byte{'\r', '\n'}
)

type command interface {
	arguments() []string
	callback() func(*response)
}

type Executor interface {
	Execute(command)
}

func (this Client) Execute(command command) {
	go this.useConnection(func(conn Connection) {
		conn.Execute(command)
	})
}

func (this Connection) input(command command) {
	_, err := this.Write(buildCommand(command.arguments()))
	checkForError(err)
}

func (this Connection) output(command command) {
	command.callback()(getResponse(this))
}

func (this Connection) Execute(command command) {
	this.input(command)
	this.output(command)
}

func buildCommand(arguments []string) []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteByte(isMultibulk)
	buf.WriteString(itoa(len(arguments)))
	buf.Write(delimiter)

	for _, arg := range arguments {
		buf.WriteByte(isBulk)
		buf.WriteString(itoa(len(arg)))
		buf.Write(delimiter)
		buf.WriteString(arg)
		buf.Write(delimiter)
	}

	return buf.Bytes()
}

func getResponse(conn io.Reader) *response {
	var buffer [1]byte
	_, err := conn.Read(buffer[:])
	if err != nil {
		panic(err)
	}
	switch buffer[0] {
	case isError:
		panic("Redis Error:" + getLongString(conn))
	case isStatus, isInt:
		return getStringResponse(conn)
	case isBulk:
		return getBulk(conn)
	case isMultibulk:
		return getMultiBulk(conn)
	}
	return nil
}

func getString(conn io.Reader) string {
	var buffer [bufferSize]byte
	j := -len(delimiter)
	i := 0
	for {
		if j >= 0 && bytes.Equal(buffer[j:i], delimiter) {
			return string(buffer[:j])
		}
		if i >= bufferSize {
			panic("short buffer")
		}
		conn.Read(buffer[i : i+1])
		i++
		j++
	}
	return string(buffer[:])
}

func getLongString(conn io.Reader) string {
	var buffer [longBufferSize]byte
	j := -len(delimiter)
	i := 0
	for {
		if j >= 0 && bytes.Equal(buffer[j:i], delimiter) {
			return string(buffer[:j])
		}
		if i >= bufferSize {
			panic("short buffer")
		}
		conn.Read(buffer[i : i+1])
		i++
		j++
	}
	return string(buffer[:])
}

func getStringResponse(conn io.Reader) *response {
	return &response{
		val: getString(conn),
	}
}

func getBulk(conn io.Reader) *response {
	line := getString(conn)
	strlen := atoi(line)
	if strlen == -1 {
		return nil
	}

	b := make([]byte, strlen+len(delimiter))
	i, err := conn.Read(b)
	if err != nil {
		//the read should be successful
		panic(err)
	}
	if i != strlen+len(delimiter) {
		//the read should go through every byte we have set out
		panic("underread")
	}
	if !bytes.Equal(b[strlen:], delimiter) {
		//the read should end with a crlf
		panic("Incorrect Redis bulk length")
	}

	return &response{
		val: string(b[:strlen]),
	}
}

func getMultiBulk(conn io.Reader) *response {
	line := getString(conn)
	cResponses := atoi(string(line))
	if cResponses == -1 {
		return nil
	}

	r := new(response)
	r.subresponses = make([]*response, cResponses)

	for iResponse := 0; iResponse < int(cResponses); iResponse++ {
		r.subresponses[iResponse] = getResponse(conn)
	}
	return r
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

func (this BoolCommand) callback() func(*response) {
	return func(r *response) {
		defer close(this.output)
		if r != nil {
			this.output <- r.val == "1"
		}
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

func (this IntCommand) callback() func(*response) {
	return func(r *response) {
		defer close(this.output)
		if r != nil {
			this.output <- atoi(r.val)
		}
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

func (this FloatCommand) callback() func(*response) {
	return func(r *response) {
		defer close(this.output)
		if r != nil {
			this.output <- atof(r.val)
		}
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

func (this StringCommand) callback() func(*response) {
	return func(r *response) {
		defer close(this.output)

		if r != nil {
			this.output <- r.val
		}
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

func (this SliceCommand) callback() func(*response) {
	return func(r *response) {
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

func (this MaybeSliceCommand) callback() func(*response) {
	return func(r *response) {
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

func (this MapCommand) callback() func(*response) {
	return func(r *response) {
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

func (this NilCommand) callback() func(*response) {
	return func(r *response) {
		this.output <- nothing{}
		close(this.output)
	}
}
