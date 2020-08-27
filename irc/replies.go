package irc

type reply int

const (
	ERR_NEEDMOREPARAMS   reply = 461
	ERR_USERSDISABLED    reply = 446
	ERR_ALREADYREGISTRED reply = 462
	ERR_NONICKNAMEGIVEN  reply = 431
	ERR_NOORIGIN         reply = 409
	ERR_NOSUCHCHANNEL    reply = 403
	ERR_NORECIPIENT      reply = 411
	ERR_NOTEXTTOSEND     reply = 412
	ERR_NOSUCHNICK       reply = 401

	RPL_WELCOME  reply = 1
	RPL_YOURHOST reply = 2
	RPL_CREATED  reply = 3
	RPL_MYINFO   reply = 4

	RPL_MOTDSTART reply = 375
	RPL_MOTD      reply = 372
	RPL_ENDOFMOTD reply = 376
)
