[![Build
Status](https://travis-ci.org/TheRealBill/commissar.svg?branch=master)](https://travis-ci.org/TheRealBill/commissar)

# Commissar: A Redis Use-Case Performance testing framework

This is the beginning of a general use-case based Redis performance testing
suite.

The goal is not to obtain "the fastest/biggest" results, but to give real-world
expectations based on use cases. Fo example if you're implementing a
leaderboard you're concerned with the performance of game result updates, how
many views of the leaderboard per second/minute can be handled, etc..

By breaking Redis performance into use-cases we can not only isolate common
scenarios but we can develop, as a community, the top implementations of
leaderboards and how the performance profiles of each differ.

## Short Term - first use-case: leaderboards

Currently the tool simulates storing results from a game of Tic-Tac-Toe. It can
be configured with a given number of game servers, user count, matches per
user, and number of "observers". These are all done via environment variables.
This allows you to leverage tools such as Docker to alter each run while
maintaining a consistent binary.

These are the current environment variables it uses:`

COMMISSAR_NAME                  
COMMISSAR_REDISCONNECTIONSTRING
COMMISSAR_REDISAUTHTOKEN      
COMMISSAR_JSONOUT            
COMMISSAR_GAMESERVERCOUNT 
COMMISSAR_READERCOUNT       
COMMISSAR_USERCOUNT        
COMMISSAR_MATCHESPERUSER  
COMMISSAR_TOTALMATCHES   
COMMISSAR_POOLSIZE       
COMMISSAR_PIPELINE      
COMMISSAR_GAMENAME     


## Medium-term - Develop and Document Use Cases


The goal here is to develop a specific set of uses cases along with their
implementation from the Redis usage pattern. With common use cases defined and
implemented we will be able to develop tests which 
  a) Demonstrate the use case
  b) document how and why it works the way it does, and 
  c) show how it performs.

As we develop these we will also be able to implement these in various
languages with multiple client side libraries.

# Get Involved!

Fork the repo, add stuff, submit issues (a great place to make suggestions for
use cases). There is also a Google group for the more in-depth conversations:
commissar@googlegroups.com


