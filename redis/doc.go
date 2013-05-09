/*
SimpleRedis is an object-oriented Redis library for golang.

Object-Oriented

SimpleRedis prefers to use an Object-Oriented approach to using Redis
In most libraries, You would use something along the lines of:

        Redis.Set("Test_String","Hello World")
        str := Redis.Get("Test_String")
	
In SimpleRedis, instead of calling the functions directly, you use:

        s := Redis.String("Test_String")
        <-s.Set("Hello World")
        str := <-s.Get()
	
This accomplishes a few things:

a) By Default, the "Test_String" only gets defined in one place, so there are fewer chances for mistyping errors

b) It becomes easier to look up which operations are usable for different types of data

c) It more accurately models how one tends to think about the data, which is typically in terms of the Redis primitives rather than the functions
	
If you do need to call the functions directly, You can call any of the "Command" functions in command.go

Concurrency

While Redis is blazing fast, it *still* has to use network I/O, often times there will be things you can do while that is happening
	s.Get()
returns a channel, which, when Redis has returned information, will contain a string.  
If you want the data immediately, you should use 
	str := <-s.Get()
The reasons for doing this are:

a) Helps to remind you that you can do things while waiting for Redis

b) Some operations (e.g. anything sent within a transaction) don't return immediately, and the result can only be obtained by waiting

c) Gives a natural interface for dealing with situations when Redis won't return anything (e.g. Popping from an empty List - "str,ok := <-l.LeftPop()")

d) Makes it easier to control when you pause for Redis

Usage

1) Figure out how you plan on connecting to Redis, and get a Config object set up properly

2) Use the Config to create a Client object.  This Client object will likely end up global, and at the very minimum needs to be accesible by all of the object methods that define Redis Objects in the next step

3) Create private methods for your objects that return Redis Objects

		//we use a namespace to prevent any redis data from the user conflicting with any other namespace
		func (u *User) base() Redis.Prefix {
		    return global.Redis.Prefix("User:"+u.id+":")
		}

		//define the user's friends by a redis set containing their IDs
		func (u *User) friends() Redis.IntSet {
		    return u.base().IntSet("Friends")
		}

4) Create public methods that interact with these objects

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