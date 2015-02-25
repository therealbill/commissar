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

COMMISSAR_GAMESERVERCOUNT
COMMISSAR_NAME                  
COMMISSAR_REDISCONNECTIONSTRING
COMMISSAR_REDISAUTHTOKEN      
COMMISSAR_JSONOUT            
COMMISSAR_READERCOUNT       
COMMISSAR_USERCOUNT        
COMMISSAR_MATCHESPERUSER  
COMMISSAR_TOTALMATCHES   
COMMISSAR_GAMESERVERCOUNT 
COMMISSAR_POOLSIZE       
COMMISSAR_PIPELINE      
COMMISSAR_GAMENAME     
