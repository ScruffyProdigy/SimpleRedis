package redis

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"time"
)

//The Config details how you plan to go about communicating with Redis
type Config struct {
	NetType         string `json:"nettype"`
	NetAddress      string `json:"netaddr"`
	DBid            int    `json:"dbid"`
	Password        string `json:"password"`
	ConnectionCount int    `json:"conncount"`
}

//DefaultConfiguration returns a config with the easiest method for communicating with Redis.
//All of the fields are public, so anything that needs to be changed for your setup can be done without affecting other fields
func DefaultConfiguration() Config {
	return Config{
		NetType:         "tcp",
		NetAddress:      "127.0.0.1:6379",
		DBid:            0,
		Password:        "",
		ConnectionCount: 100,
	}
}

type errCallbackFunc func(error, string)

func (this errCallbackFunc) Call(e error, s string) {
	if this == nil {
		panic(errors.New(e.Error() + ":" + s))
	} else {
		this(e, s)
	}
}

func getError(rec interface{}) error {
	if err, ok := rec.(error); ok {
		return err
	}
	if str, ok := rec.(string); ok {
		return errors.New(str)
	}
	return errors.New("Unknown Error:" /*+fmt.Sprintf(rec)*/)
}

// The Client is the base for all communication to and from Redis
type Client struct {
	nextID       int
	isClosed     bool
	pool         chan *Connection // 	a semaphore of connections to draw from when multiple threads want to connect
	config       Config           //	connection details, so we know how to connect to redis
	fErrCallback errCallbackFunc  //	a callback function - since we operate in a separate goroutine, we can't return an error, instead we call this function sending it the error, and the command we tried to issue
}

//New gives back a Client that communicates using the details specified in the supplied Config
func New(config Config) (r *Client, e error) {
	//user has not had a chance to set an error callback at this point
	//so we should exit gracefully if an error happens during load
	defer func() {
		rec := recover()
		if rec != nil {
			r = nil
			e = getError(rec)
		}
	}()

	this := new(Client)
	this.config = config

	this.pool = make(chan *Connection, config.ConnectionCount)
	for i := 0; i < config.ConnectionCount; i++ {
		conn, err := this.newConnection()
		if err != nil {
			return nil, err
		}

		this.pool <- conn
	}

	return this, nil
}

//Load reads in information, and uses the JSON information it finds therein to find the communcation hookup details for Redis.
//It then returns a Client based on the supplied information
func Load(configfile io.Reader) (*Client, error) {
	config := DefaultConfiguration()
	dec := json.NewDecoder(configfile)
	err := dec.Decode(&config)
	if err != nil {
		return nil, err
	}

	return New(config)
}

//Close frees up all connections previously allocated.
//BUG: If you have connections still in use, things can get messy
func (this *Client) Close() error {
	if this.isClosed {
		return errors.New("Redis is already closed!")
	}
	this.isClosed = true

	timeout := time.After(1 * time.Second)
	for numClosed := 0; numClosed < this.config.ConnectionCount; numClosed++ {
		select {
		case conn := <-this.pool:
			conn.Close()
		case <-timeout:
			this.errCallback(errors.New("Connections are still in use"), "Closing Redis")
			return errors.New("Could not close all connections")
		}
	}
	close(this.pool)

	return nil
}

//Execute allows commands to be executed directly through the Client without needing to specify a key
func (this Client) Execute(command command) {
	go this.useConnection(func(conn *Connection) {
		conn.Execute(command)
	})
}

func (this Client) errCallback(e error, s string) {
	this.fErrCallback.Call(e, s)
}

//Since redis operates in a separate thread, it isn't always possible to return an error status easily.
//SetErrorCallback allows you to react to an error when it happens
func (this *Client) SetErrorCallback(callback func(error, string)) {
	this.fErrCallback = errCallbackFunc(callback)
}

func (this *Client) newConnection() (*Connection, error) {
	conn, err := net.Dial(this.config.NetType, this.config.NetAddress)
	if err != nil {
		return nil, err
	}

	c := &Connection{conn, this.nextID, this}

	if this.config.Password != "" {
		<-NilCommand(c, "AUTH", this.config.Password)
	}
	if this.config.DBid != 0 {
		<-NilCommand(c, "SELECT", itoa(this.config.DBid))
	}
	this.nextID++
	return c, nil
}

func (this *Client) useConnection(callback func(*Connection)) {
	if this.isClosed {
		return
	}

	conn := <-this.pool
	defer func() {
		this.pool <- conn
	}()

	callback(conn)
}

func (this *Client) useNewConnection(callback func(*Connection)) {
	conn, err := this.newConnection()
	if err != nil {
		this.errCallback(err, "new connection")
	}

	defer func() {
		conn.Close()
	}()

	callback(conn)
}

//Creates a basic key.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) Key(key string) Key {
	return newKey(this, key)
}

//Creates a String object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) String(key string) String {
	return newString(this, key)
}

//Creates an Integer object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) Integer(key string) Integer {
	return newInteger(this, key)
}

//Creates a Float object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) Float(key string) Float {
	return newFloat(this, key)
}

//Creates a Bits object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) Bits(key string) Bits {
	return newBits(this, key)
}

//Creates a Hash object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) Hash(key string) Hash {
	return newHash(this, key)
}

//Creates a List object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) List(key string) List {
	return newList(this, key)
}

//Creates an IntList object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) IntList(key string) IntList {
	return newIntList(this, key)
}

//Creates a Set Object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) Set(key string) Set {
	return newSet(this, key)
}

//Creates an IntSet Object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) IntSet(key string) IntSet {
	return newIntSet(this, key)
}

//Creates a SortedSet Object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) SortedSet(key string) SortedSet {
	return newSortedSet(this, key)
}

//Creates a SortedIntSet Object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) SortedIntSet(key string) SortedIntSet {
	return newSortedIntSet(this, key)
}

//Creates a Mutex Object.
//(Warning - this is *not* a lightweight function - there is some network I/O involved in mutex initialization)
func (this *Client) Mutex(key string) Mutex {
	return newMutex(this, key, 1)
}

//Creates a Semaphore Object.
//(Warning - this is *not* a lightweight function - there is some network I/O involved in mutex initialization)
func (this *Client) Semaphore(key string, count int) Mutex {
	return newMutex(this, key, count)
}

//Creates a ReadWriteMutex Object.
//(Warning - this is *not* a lightweight function - there is some network I/O involved in mutex initialization)
func (this *Client) ReadWriteMutex(key string, readers int) *ReadWriteMutex {
	return newRWMutex(this, key, readers)
}

//Creates a Channel Object.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) Channel(key string) Channel {
	return newChannel(this, key)
}

//Creates a Prefix Object, which helps namespace other Redis Objects.
//(This is a lightweight function - does *not* involve network I/O)
func (this *Client) Prefix(key string) Prefix {
	return newPrefix(this, key)
}
