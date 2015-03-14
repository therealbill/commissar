
# Use Case: Leaderboard

A leaderboard is best decribed as keeping track of who is where in the
rankings of some sort of game or competition. For example, assume a bevy
of 50 players playing Tic-Tac-Toe. After every game you want to track
the winner and loser, and provide a leaderboard which displays the
rankings.

There are several ways to do this in general, which means several ways
to do it in Redis. In the above example you wanted to track win/lose.
However, some games have points earned or lost at the end of a game as
well. This is another form of leaderboard.

Thus the Leaderboard section of Commissar will need to describe cased
where you update Redis with:
a) Win/Lose only
b) Last Score
c) Points gained
d) Combinations of the above

# Leaderboard Tunables

A leaderboard will need to have a variable number of users, variable number of
matches per user, variable numbers of game servers and readers, and
controllable number of iterations.

# Leaderboard Output

There are two key metrics for leaderboards:
  - How many game results per second can we process?
  - How many leaderboard views per second can we process?



# Implementation Variants

There are several ways to implement each of the above board types. Each
implementation will need to be detailed and implemented, as well as
categorized for valid comparisons.

# Board Type: Win/Lose Only

In this scenario we are only concerned with keeping track of wins and
losses.


# Board Type: Last Score

In this scenario the player gets a score at the end of the game. We rank
players based on the last score they receveid.

# Board Type: Points Gained

In this scenario the players earn, or lose, a number of points. Wins or
losses are irrelevant. Thus we modify the running total at the game's
conclusion and display ranking based on current score.


