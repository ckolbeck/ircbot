//Copyright 2010 Cory Kolbeck <ckolbeck@gmail.com>.
//So long as this notice remains in place, you are welcome 
//to do whatever you like to or with this code.  This code is 
//provided 'As-Is' with no warrenty expressed or implied. 
//If you like it, and we happen to meet, buy me a beer sometime

package main

import (
	irc "github.com/ckolbeck/ircbot"
	"fmt"
)

//Implement a very simple bot which will respond to anyone addressing it and join
//any channels it is invited to.
var bot *irc.Bot

func main() {
	var err error

	//Create a new bot
	bot, err = irc.NewBot(
		"goirc-bot",         //The bot's nick
		"",                  //The bot's nickserv password (blank for none)
		"www.github.com",    //The bot's domain
		"chat.freenode.org", //IRC server to connect to
		7070,                //Remote port to connect on
		true,                //Use ssl?
		'!',                 //Char used to address the bot
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	//Control how the bot will react to various IRC commands using the bot's
	//Actions map.  By default the only thing the bot will do is respond to PING
	//requests.  Functions must have the signature:
	// `func(bot *irc.Bot, msg *irc.Message) *irc.Message`
	bot.Actions["INVITE"] = join

	//PRIVMSG can be handled as above, or you can use the convenience method below.
	//It accepts two functions, the first will be called for messages directed
	//to the bot via its attention char, addressing it by name, or by sending it a 
	//private message.  The second will be called for all other messages
	bot.SetPrivmsgHandler(sayHi, ctcpEcho)

	//Attempt to join the given channel, here we're joining an unkeyed channel.
	bot.JoinChannel("#echo", "")

	//No further work to be done in main, block indefinitely
	select {}
}

//Whether the bot is addressed with its attention char, in a private message, or with example-bot:
//cmd will hold the meaningful part of the message
//msg holds the raw irc message broken into prefix, command, args, trailing, and possibly a CTCP command
//For PRIVMSGs, trailing holds the text of the message
func sayHi(cmd string, msg *irc.Message) string {
	//CTCP Version needs to be handled here, since it will generally be sent in a pm
	if msg.Ctcp == "VERSION" { 
		return "GoIRC example bot"
	}

	return fmt.Sprintf("Hi there, %s. You said: %s", msg.GetSender(), msg.Trailing)
}

//This will be called for any message not directed at the bot
//Here, we just listen for any CTCP actions (usually sent with '/me') and copy them
//Note: PMs will not have ctcp actions echoed, since this method will not be called
func ctcpEcho(cmd string, msg *irc.Message) string {
	//The convenience methods are limited to sending
	//simple text PRIVMSGs back to the source of the
	//incoming message.  For more complex behavior,
	//Send(*Message) can be used.
	if msg.Ctcp == "ACTION" {
		bot.Send(&irc.Message{
			Command:  "PRIVMSG",
			Args:     msg.Args[:1],
			Ctcp:     "ACTION",
			Trailing: msg.Trailing,
		})
	}

	return "" //Nothing else will be sent if an empty string is returned 
}

//Function to join any channels the bot is invited to.
func join(bot *irc.Bot, msg *irc.Message) *irc.Message {
	return &irc.Message{
		Command: "JOIN",
		Args:    []string{msg.Trailing},
	}
}
