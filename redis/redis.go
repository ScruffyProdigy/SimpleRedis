/*
SimpleRedis is an object-oriented Redis library for golang

There are a couple of differences between this library and others.

1) It is object oriented.

In most libraries, You have something along the lines of:

        Redis.Set("Test_String","Hello World")
        str := Redis.Get("Test_String")
	
In this one, instead of calling the functions directly, you use:

        s := Redis.String("Test_String")
        <-s.Set("Hello World")
        str := <-s.Get()
	
This accomplishes a few things:

* By Default, the "Test_String" only gets defined in one place, so there are fewer chances for mistyping errors
* It becomes easier to look up which operations are usable for different types of data
* It more accurately models how one tends to think about the data, which is typically in terms of the Redis primitives rather than the functions
	
If you do need to call the functions directly, You can call any of the "Command" functions in command.go

2) It uses channels

While Redis is blazing fast, it *still* has to use network I/O, and often times there will be things you can do while that is happening

"s.Get()" returns a channel, which, when Redis has returned information, will contain a string.  If you want the data immediately, you should use "`str := <-s.Get()`"

The reasons for doing this are:

* Helps to remind you that you can do things while waiting for Redis
* Some operations (e.g. anything sent within a transaction) don't return immediately, and the result can only be obtained by waiting
* Gives a natural interface for dealing with situations when Redis won't return anything (e.g. Popping from an empty List - "str,ok := <-l.LeftPop()")
* Makes it easier to control

*** Usage ***

* Figure out how you plan on connecting to Redis, and get a Config object set up properly
* Use the Config to create a Client object
	* You will probably make this object global
	* if not, make sure any object that needs to define Redis Objects has access to it
* Create methods for your objects that return Redis Objects
	* defining a Redis Object is a very lightweight operation, you should not need to be worried about the overhead
	* these methods should probably be private

		func (u *User) base() Redis.Prefix {
		//namespacing everything from within the user to help prevent clashes
		    return global.Redis.Prefix("User:"+u.id+":")
		}

		func (u *User) friends() Redis.IntSet {
		    return u.base().IntSet("Friends")
		}

* Create methods that interact with these objects
	* these methods will probably be public

		//note: not using the channel arrows, because this is not a time-sensitive operation
		func (u *User) AddFriend(otherUser *User) {
			u.friends().Add(otherUser.id)
			otherUser.friends().Add(u.id)
		}
	
		func (u *User) Unfriend(otherUser *User) {
			u.friends().Remove(otherUser.id)
			otherUser.friends().Remove(u.id)
		}
*/
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

//DefaultConfiguration returns a config with the easiest method for communicating with Redis
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

//Load reads in information, and uses the JSON information it finds therein to find the communcation hookup details for Redis
//it then returns a Client based on the supplied information
func Load(configfile io.Reader) (*Client, error) {
	config := DefaultConfiguration()
	dec := json.NewDecoder(configfile)
	err := dec.Decode(&config)
	if err != nil {
		return nil, err
	}

	return New(config)
}

//Close frees up all connections previously allocated
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

//Since redis operates in a separate thread, it isn't always possible to return an error status easily
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

//Creates a basic key
func (this *Client) Key(key string) Key {
	return newKey(this, key)
}

//Creates a String object
func (this *Client) String(key string) String {
	return newString(this, key)
}

//Creates an Integer object
func (this *Client) Integer(key string) Integer {
	return newInteger(this, key)
}

//Creates a Float object
func (this *Client) Float(key string) Float {
	return newFloat(this, key)
}

//Creates a Bits object
func (this *Client) Bits(key string) Bits {
	return newBits(this, key)
}

//Creates a Hash object
func (this *Client) Hash(key string) Hash {
	return newHash(this, key)
}

//Creates a List object
func (this *Client) List(key string) List {
	return newList(this, key)
}

//Creates an IntList object
func (this *Client) IntList(key string) IntList {
	return newIntList(this, key)
}

//Creates a Set Object
func (this *Client) Set(key string) Set {
	return newSet(this, key)
}

//Creates an IntSet Object
func (this *Client) IntSet(key string) IntSet {
	return newIntSet(this, key)
}

//Creates a SortedSet Object
func (this *Client) SortedSet(key string) SortedSet {
	return newSortedSet(this, key)
}

//Creates a SortedIntSet Object
func (this *Client) SortedIntSet(key string) SortedIntSet {
	return newSortedIntSet(this, key)
}

//Creates a Mutex Object
func (this *Client) Mutex(key string) Mutex {
	return newMutex(this, key, 1)
}

//Creates a Semaphore Object
func (this *Client) Semaphore(key string, count int) Mutex {
	return newMutex(this, key, count)
}

//Creates a ReadWriteMutex Object
func (this *Client) ReadWriteMutex(key string, readers int) *ReadWriteMutex {
	return newRWMutex(this, key, readers)
}

//Creates a Channel Object
func (this *Client) Channel(key string) Channel {
	return newChannel(this, key)
}

//Creates a Prefix Object
func (this *Client) Prefix(key string) Prefix {
	return newPrefix(this, key)
}
