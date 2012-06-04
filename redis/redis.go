package redis

import (
	"encoding/json"
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
		ConnectionCount: 1,
	}
}

var nextID int

type Connection struct {
	net.Conn
	id int
}

type Client struct {
	pool   chan Connection
	config Config
}

func New(config Config) *Client {
	this := new(Client)
	this.config = config

	this.pool = make(chan Connection, config.ConnectionCount)
	for i := 0; i < config.ConnectionCount; i++ {
		this.pool <- this.newConnection()
	}

	//	err := this.Test()
	//	checkForError(err)
	return this
}

func Load(configfile io.Reader) *Client {
	config := DefaultConfiguration()
	dec := json.NewDecoder(configfile)
	err := dec.Decode(&config)
	checkForError(err)

	return New(config)
}

func (this *Client) newConnection() Connection {
	conn, err := net.Dial(this.config.NetType, this.config.NetAddress)
	checkForError(err)

	if this.config.Password != "" {
		//		this.UsePassword(config.Password).WithConn(conn)
	}
	if this.config.DBid != 1 {
		//		this.UseDatabase(config.DBid).WithConn(conn)		
	}
	c := Connection{conn, nextID}
	nextID++
	return c
}

func (this *Client) useConnection(callback func(Connection)) {
	conn := <-this.pool
	defer func() {
		this.pool <- conn
	}()

	callback(conn)
}

func (this *Client) useNewConnection(callback func(Connection)) {
	conn := this.newConnection()
	defer func() {
		conn.Close()
	}()

	callback(conn)
}

func checkForError(err error) {
	if err != nil {
		panic(err)
	}
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
