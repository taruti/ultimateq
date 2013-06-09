/*
irc package defines types and classes to be used by most other packages in
the ultimateq system. It is small and comprised mostly of helper like types
and constants.
*/
package irc

import "strings"

// IRC Messages, these messages are 1-1 constant to string lookups for ease of
// use when registering handlers etc.
const (
	JOIN    = "JOIN"
	KICK    = "KICK"
	MODE    = "MODE"
	NICK    = "NICK"
	NOTICE  = "NOTICE"
	PART    = "PART"
	PING    = "PING"
	PONG    = "PONG"
	PRIVMSG = "PRIVMSG"
	QUIT    = "QUIT"
	TOPIC   = "TOPIC"
)

// IRC Reply and Error Messages. These are sent in reply to a previous message.
const (
	RPL_WELCOME         = "001"
	RPL_YOURHOST        = "002"
	RPL_CREATED         = "003"
	RPL_MYINFO          = "004"
	RPL_ISUPPORT        = "005"
	RPL_BOUNCE          = "005"
	RPL_USERHOST        = "302"
	RPL_ISON            = "303"
	RPL_AWAY            = "301"
	RPL_UNAWAY          = "305"
	RPL_NOWAWAY         = "306"
	RPL_WHOISUSER       = "311"
	RPL_WHOISSERVER     = "312"
	RPL_WHOISOPERATOR   = "313"
	RPL_WHOISIDLE       = "317"
	RPL_ENDOFWHOIS      = "318"
	RPL_WHOISCHANNELS   = "319"
	RPL_WHOWASUSER      = "314"
	RPL_ENDOFWHOWAS     = "369"
	RPL_LISTSTART       = "321"
	RPL_LIST            = "322"
	RPL_LISTEND         = "323"
	RPL_UNIQOPIS        = "325"
	RPL_CHANNELMODEIS   = "324"
	RPL_NOTOPIC         = "331"
	RPL_TOPIC           = "332"
	RPL_INVITING        = "341"
	RPL_SUMMONING       = "342"
	RPL_INVITELIST      = "346"
	RPL_ENDOFINVITELIST = "347"
	RPL_EXCEPTLIST      = "348"
	RPL_ENDOFEXCEPTLIST = "349"
	RPL_VERSION         = "351"
	RPL_WHOREPLY        = "352"
	RPL_ENDOFWHO        = "315"
	RPL_NAMREPLY        = "353"
	RPL_ENDOFNAMES      = "366"
	RPL_LINKS           = "364"
	RPL_ENDOFLINKS      = "365"
	RPL_BANLIST         = "367"
	RPL_ENDOFBANLIST    = "368"
	RPL_INFO            = "371"
	RPL_ENDOFINFO       = "374"
	RPL_MOTDSTART       = "375"
	RPL_MOTD            = "372"
	RPL_ENDOFMOTD       = "376"
	RPL_YOUREOPER       = "381"
	RPL_REHASHING       = "382"
	RPL_YOURESERVICE    = "383"
	RPL_TIME            = "391"
	RPL_USERSSTART      = "392"
	RPL_USERS           = "393"
	RPL_ENDOFUSERS      = "394"
	RPL_NOUSERS         = "395"
	RPL_TRACELINK       = "200"
	RPL_TRACECONNECTING = "201"
	RPL_TRACEHANDSHAKE  = "202"
	RPL_TRACEUNKNOWN    = "203"
	RPL_TRACEOPERATOR   = "204"
	RPL_TRACEUSER       = "205"
	RPL_TRACESERVER     = "206"
	RPL_TRACESERVICE    = "207"
	RPL_TRACENEWTYPE    = "208"
	RPL_TRACECLASS      = "209"
	RPL_TRACERECONNECT  = "210"
	RPL_TRACELOG        = "261"
	RPL_TRACEEND        = "262"
	RPL_STATSLINKINFO   = "211"
	RPL_STATSCOMMANDS   = "212"
	RPL_ENDOFSTATS      = "219"
	RPL_STATSUPTIME     = "242"
	RPL_STATSOLINE      = "243"
	RPL_UMODEIS         = "221"
	RPL_SERVLIST        = "234"
	RPL_SERVLISTEND     = "235"
	RPL_LUSERCLIENT     = "251"
	RPL_LUSEROP         = "252"
	RPL_LUSERUNKNOWN    = "253"
	RPL_LUSERCHANNELS   = "254"
	RPL_LUSERME         = "255"
	RPL_ADMINME         = "256"
	RPL_ADMINLOC1       = "257"
	RPL_ADMINLOC2       = "258"
	RPL_ADMINEMAIL      = "259"
	RPL_TRYAGAIN        = "263"

	ERR_NOSUCHNICK        = "401"
	ERR_NOSUCHSERVER      = "402"
	ERR_NOSUCHCHANNEL     = "403"
	ERR_CANNOTSENDTOCHAN  = "404"
	ERR_TOOMANYCHANNELS   = "405"
	ERR_WASNOSUCHNICK     = "406"
	ERR_TOOMANYTARGETS    = "407"
	ERR_NOSUCHSERVICE     = "408"
	ERR_NOORIGIN          = "409"
	ERR_NORECIPIENT       = "411"
	ERR_NOTEXTTOSEND      = "412"
	ERR_NOTOPLEVEL        = "413"
	ERR_WILDTOPLEVEL      = "414"
	ERR_BADMASK           = "415"
	ERR_UNKNOWNCOMMAND    = "421"
	ERR_NOMOTD            = "422"
	ERR_NOADMININFO       = "423"
	ERR_FILEERROR         = "424"
	ERR_NONICKNAMEGIVEN   = "431"
	ERR_ERRONEUSNICKNAME  = "432"
	ERR_NICKNAMEINUSE     = "433"
	ERR_NICKCOLLISION     = "436"
	ERR_UNAVAILRESOURCE   = "437"
	ERR_USERNOTINCHANNEL  = "441"
	ERR_NOTONCHANNEL      = "442"
	ERR_USERONCHANNEL     = "443"
	ERR_NOLOGIN           = "444"
	ERR_SUMMONDISABLED    = "445"
	ERR_USERSDISABLED     = "446"
	ERR_NOTREGISTERED     = "451"
	ERR_NEEDMOREPARAMS    = "461"
	ERR_ALREADYREGISTRED  = "462"
	ERR_NOPERMFORHOST     = "463"
	ERR_PASSWDMISMATCH    = "464"
	ERR_YOUREBANNEDCREEP  = "465"
	ERR_YOUWILLBEBANNED   = "466"
	ERR_KEYSET            = "467"
	ERR_CHANNELISFULL     = "471"
	ERR_UNKNOWNMODE       = "472"
	ERR_INVITEONLYCHAN    = "473"
	ERR_BANNEDFROMCHAN    = "474"
	ERR_BADCHANNELKEY     = "475"
	ERR_BADCHANMASK       = "476"
	ERR_NOCHANMODES       = "477"
	ERR_BANLISTFULL       = "478"
	ERR_NOPRIVILEGES      = "481"
	ERR_CHANOPRIVSNEEDED  = "482"
	ERR_CANTKILLSERVER    = "483"
	ERR_RESTRICTED        = "484"
	ERR_UNIQOPPRIVSNEEDED = "485"
	ERR_NOOPERHOST        = "491"
	ERR_UMODEUNKNOWNFLAG  = "501"
	ERR_USERSDONTMATCH    = "502"
)

// Pseudo Messages, these messages are not real messages defined by the irc
// protocol but the bot provides them to allow for additional messages to be
// handled such as connect or disconnects which the irc protocol has no protocol
// defined for.
const (
	RAW        = "RAW"
	CONNECT    = "CONNECT"
	DISCONNECT = "DISCONNECT"
)

// Sender is the sender of an event, and should allow replies on a writing
// interface as well as a way to identify itself.
type Sender interface {
	// Writes a string to an endpoint that makes sense for the given event.
	Writeln(string) error
	// Retrieves a key to retrieve where this event was generated from.
	GetKey() string
}

// IrcMessage contains all the information broken out of an irc message.
type IrcMessage struct {
	// Name of the message. Uppercase constant name or numeric.
	Name string
	// The server or user that sent the message, a fullhost if one was supplied.
	Sender string
	// The args split by space delimiting.
	Args []string
}

// Split splits string arguments. A convenience method to avoid having to call
// splits and import strings.
func (m *IrcMessage) Split(index int) []string {
	return strings.Split(m.Args[index], ",")
}

// Message type provides a view around an IrcMessage to access it's parts in a
// more convenient way.
type Message struct {
	// Raw is the underlying irc message.
	Raw *IrcMessage
}

// Target retrieves the channel or user this message was sent to.
func (p *Message) Target() string {
	return p.Raw.Args[0]
}

// Message retrieves the message sent to the user or channel.
func (p *Message) Message() string {
	return p.Raw.Args[1]
}
