# VTM bot
This is a bot to be used in rolling dice for VTM

After downloading the files run `go build` to create the executable

To run the bot use `./DiscrodBot -t BOT_TOKEN`

! or / can be used for legacy commands

## Commands
/roll - prompts user for the following - [dice pool] [dc] [action]{optional}

/reroll - rerolls lowest three dice

/coin - flips a coin - no beginning character required

/wyk - for storyteller use

/saveall - save all users to database

## Legacy Commands - no longer in use 
/N,DC - N is number of dice to roll, DC is DC of the check - Rolls N dice, does proper math vs DC to get number of sucesses)

/N,DC,Reason - similar to rolling dice but outputs the reason for the roll as well


## Commands in Progress
/testing - debug command, force pushes data to the datbase currently.

/Schedule - used to schedule future sessions between players


