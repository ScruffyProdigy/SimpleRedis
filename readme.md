SimpleRedis
===========

SimpleRedis is an object-oriented Redis library for golang

There are a couple of differences between this library and others.

###1) It is object oriented.

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

###2) It uses channels

While Redis is blazing fast, it *still* has to use network I/O, and often times there will be things you can do while that is happening

"`s.Get()`" returns a channel, which, when Redis has returned information, will contain a string.  If you want the data immediately, you should use "`str := <-s.Get()`"

The reasons for doing this are:

* Helps to remind you that you can do things while waiting for Redis
* Some operations (e.g. anything sent within a transaction) don't return immediately, and the result can only be obtained by waiting
* Gives a natural interface for dealing with situations when Redis won't return anything (e.g. Popping from an empty List - "`str,ok := <-l.LeftPop()`")
* Makes it easier to control

Usage
-----

* Figure out how you plan on connecting to Redis, and get a Config object set up properly
* Use the Config to create a Client object
	* You will probably make this object global
	* if not, make sure any object that needs to define Redis Objects has access to it
* Create methods for your objects that return Redis Objects
	* defining a Redis Object is a very lightweight operation, you should not need to be worried about the overhead
	* these methods should probably be private

*Example:*

            func (u *User) base() Redis.Prefix {
				//namespacing everything from within the user to help prevent clashes
                return global.Redis.Prefix("User:"+u.id+":")
            }
    
            func (u *User) friends() Redis.IntSet {
                return u.base().IntSet("Friends")
            }

* Create methods that interact with these objects
	* these methods will probably be public

*Example:*
			//note: not using the channel arrows, because this is not a time-sensitive operation
			func (u *User) AddFriend(otherUser *User) {
				u.friends().Add(otherUser.id)
				otherUser.friends().Add(u.id)
			}
			
			func (u *User) Unfriend(otherUser *User) {
				u.friends().Remove(otherUser.id)
				otherUser.friends().Remove(u.id)
			}
			
For more usage details, please see: http://godoc.org/github.com/ScruffyProdigy/SimpleRedis/redis