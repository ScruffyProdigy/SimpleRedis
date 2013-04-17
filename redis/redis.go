package redis

import (
	"encoding/json"
	"errors"
	"io"
	"net"
)

type Config struct {
	NetType         string `json:"nettype"`
	NetAddress      string `json:"netaddr"`
	DBid            int    `json:"dbid"`
	Password        string `json:"password"`
	ConnectionCount int    `json:"conncount"`
}

func DefaultConfiguration() Config {
	return Config{
		NetType:         "tcp",
		NetAddress:      "127.0.0.1:6379",
		DBid:            0,
		Password:        "",
		ConnectionCount: 100,
	}
}

type errCallback func(error, string)

func (this errCallback) Call(e error, s string) {
	if this == nil {
		panic(errors.New(e.Error() + ":" + s))
	} else {
		this(e, s)
	}
}

type Client struct {
	nextID      int
	isClosed    bool
	pool        chan *Connection     // 	a semaphore of connections to draw from when multiple threads want to connect
	used        map[*Connection]bool //a set of all connections currently being used
	config      Config               //	connection details, so we know how to connect to redis
	errCallback errCallback          //	a callback function - since we operate in a separate goroutine, we can't return an error, instead we call this function sending it the error, and the command we tried to issue
}

func New(config Config) (*Client, error) {
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

	this.used = make(map[*Connection]bool)

	return this, nil
}

func Load(configfile io.Reader) (*Client, error) {
	config := DefaultConfiguration()
	dec := json.NewDecoder(configfile)
	err := dec.Decode(&config)
	if err != nil {
		return nil, err
	}

	return New(config)
}

func (this *Client) Close() {
	if this.isClosed {
		return
	}
	this.isClosed = true

	close(this.pool)
	for conn := range this.pool {
		conn.Close()
	}

	for conn, _ := range this.used {
		conn.Close()
	}
}

func (this Client) Execute(command command) {
	go this.useConnection(func(conn *Connection) {
		conn.Execute(command)
	})
}

func (this Client) ErrCallback(e error, s string) {
	this.errCallback.Call(e, s)
}

func (this *Client) SetErrorCallback(callback func(error, string)) {
	this.errCallback = errCallback(callback)
}

func (this *Client) newConnection() (*Connection, error) {
	conn, err := net.Dial(this.config.NetType, this.config.NetAddress)
	if err != nil {
		return nil, err
	}

	c := &Connection{conn, this.nextID, this}

	if this.config.Password != "" {
		<-NilCommand(c, []string{"AUTH", this.config.Password})
	}
	if this.config.DBid != 0 {
		<-NilCommand(c, []string{"SELECT", itoa(this.config.DBid)})
	}
	this.nextID++
	return c, nil
}

func (this *Client) useConnection(callback func(*Connection)) {
	if this.isClosed {
		return
	}

	conn := <-this.pool
	this.used[conn] = true
	defer func() {
		delete(this.used, conn)
		this.pool <- conn
	}()

	callback(conn)
}

func (this *Client) useNewConnection(callback func(*Connection)) {
	conn, err := this.newConnection()
	if err != nil {
		this.errCallback.Call(err, "new connection")
	}

	defer func() {
		conn.Close()
	}()

	callback(conn)
}

func (this *Client) Key(key string) Key {
	return newKey(this, key)
}

func (this *Client) String(key string) String {
	return newString(this, key)
}

func (this *Client) Integer(key string) Integer {
	return newInteger(this, key)
}

func (this *Client) Float(key string) Float {
	return newFloat(this, key)
}

func (this *Client) Bits(key string) Bits {
	return newBits(this, key)
}

func (this *Client) Hash(key string) Hash {
	return newHash(this, key)
}

func (this *Client) List(key string) List {
	return newList(this, key)
}

func (this *Client) IntList(key string) IntList {
	return newIntList(this, key)
}

func (this *Client) FloatList(key string) FloatList {
	return newFloatList(this, key)
}

func (this *Client) Set(key string) Set {
	return newSet(this, key)
}

func (this *Client) IntSet(key string) IntSet {
	return newIntSet(this, key)
}

func (this *Client) FloatSet(key string) FloatSet {
	return newFloatSet(this, key)
}

func (this *Client) SortedSet(key string) SortedSet {
	return newSortedSet(this, key)
}

func (this *Client) SortedIntSet(key string) SortedIntSet {
	return newSortedIntSet(this, key)
}

func (this *Client) SortedFloatSet(key string) SortedFloatSet {
	return newSortedFloatSet(this, key)
}

func (this *Client) Mutex(key string) Mutex {
	return newMutex(this, key, 1)
}

func (this *Client) Semaphore(key string, count int) Mutex {
	return newMutex(this, key, count)
}

func (this *Client) ReadWriteMutex(key string, readers int) *ReadWriteMutex {
	return newRWMutex(this, key, readers)
}

func (this *Client) Channel(key string) Channel {
	return newChannel(this, key)
}

func (this *Client) Prefix(key string) Prefix {
	return newPrefix(this, key)
}
