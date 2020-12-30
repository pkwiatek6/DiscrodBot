# VTM bot
This is a bot to be used in rolling dice for VTM

After downloading the files run `go build` to create the executable

To run the bot use `./DiscrodBot -t BOT_TOKEN`

! or / can be used for commands

## Commands
/N,DC - N is number of dice to roll, DC is DC of the check - Rolls N dice, does proper math vs DC to get number of sucesses)

/N,DC,Reason - similar to rolling dice but outputs the reason for the roll as well

/r - rerolls lowest three dice

/reroll - rerolls lowest three dice

Flip a coin - flips a coin - no beginning character required

## Commands in Progress
/testing - debug command, force pushes data to the datbase currently.

/Schedule - used to schedule future sessions between players


