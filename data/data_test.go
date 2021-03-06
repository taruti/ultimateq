package data

import (
	"github.com/aarondl/ultimateq/irc"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) } //Hook into testing package
type s struct{}

var _ = Suite(&s{})

var server = "irc.server.net"
var users = []string{"nick1!user1@host1", "nick2!user2@host2"}
var nicks = []string{"nick1", "nick2"}
var channels = []string{"#CHAN1", "#CHAN2"}

var self = Self{
	User: CreateUser("me!my@host.com"),
}

func (s *s) TestStore(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(st, NotNil)
	c.Check(err, IsNil)
	c.Check(st.Self.ChannelModes, NotNil)

	// Should die on creating kinds
	fakeCaps := &irc.ProtoCaps{}
	fakeCaps.ParseISupport(&irc.IrcMessage{Args: []string{
		"NICK", "CHANTYPES=#&", "PREFIX=(ov)@+",
	}})
	st, err = CreateStore(fakeCaps)
	c.Check(st, IsNil)
	c.Check(err, NotNil)

	// Should die on creating user modes
	fakeCaps = &irc.ProtoCaps{}
	fakeCaps.ParseISupport(&irc.IrcMessage{Args: []string{
		"NICK", "CHANTYPES=#&", "CHANMODES=a,b,c,d",
	}})
	st, err = CreateStore(fakeCaps)
	c.Check(st, IsNil)
	c.Check(err, NotNil)

	// Should die on creating ChannelFinder
	fakeCaps = &irc.ProtoCaps{}
	fakeCaps.ParseISupport(&irc.IrcMessage{Args: []string{
		"NICK", "CHANTYPES=H", "PREFIX=(ov)@+", "CHANMODES=a,b,c,d",
	}})
	st, err = CreateStore(fakeCaps)
	c.Check(st, IsNil)
	c.Check(err, NotNil)
}

func (s *s) TestStore_UpdateProtoCaps(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)

	fakeCaps := &irc.ProtoCaps{}
	fakeCaps.ParseISupport(&irc.IrcMessage{Args: []string{
		"NICK", "CHANTYPES=!", "PREFIX=(q)@", "CHANMODES=,,,q",
	}})
	fakeCaps.ParseMyInfo(&irc.IrcMessage{Args: []string{
		"irc.test.net", "test-12", "q", "abc",
	}})

	c.Assert(st.selfkinds.kinds['q'], Equals, 0)
	c.Assert(st.kinds.kinds['q'], Equals, 0)
	c.Assert(st.umodes.GetModeBit('q'), Equals, 0)
	c.Assert(st.cfinder.IsChannel("!"), Equals, false)
	st.Protocaps(fakeCaps)
	c.Assert(st.selfkinds.kinds['q'], Not(Equals), 0)
	c.Assert(st.kinds.kinds['q'], Not(Equals), 0)
	c.Assert(st.umodes.GetModeBit('q'), Not(Equals), 0)
	c.Assert(st.cfinder.IsChannel("!"), Equals, true)
}

func (s *s) TestStore_GetUser(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	c.Check(st.GetUser(users[0]), IsNil)
	c.Check(st.GetUser(users[1]), IsNil)
	st.addUser(users[0])
	c.Check(st.GetUser(users[0]), NotNil)
	c.Check(st.GetUser(users[1]), IsNil)
	st.addUser(users[1])
	c.Check(st.GetUser(users[0]), NotNil)
	c.Check(st.GetUser(users[1]), NotNil)

	st, err = CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	oldHost := "nick!user@host.com"
	newHost := "nick!user@host.net"
	st.addUser(oldHost)
	c.Check(st.GetUser(oldHost).GetFullhost(), Equals, oldHost)
	c.Check(st.GetUser(newHost).GetFullhost(), Not(Equals), newHost)
	st.addUser(newHost)
	c.Check(st.GetUser(oldHost).GetFullhost(), Not(Equals), oldHost)
	c.Check(st.GetUser(newHost).GetFullhost(), Equals, newHost)
}

func (s *s) TestStore_GetChannel(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	c.Check(st.GetChannel(channels[0]), IsNil)
	c.Check(st.GetChannel(channels[1]), IsNil)
	st.addChannel(channels[0])
	c.Check(st.GetChannel(channels[0]), NotNil)
	c.Check(st.GetChannel(channels[1]), IsNil)
	st.addChannel(channels[1])
	c.Check(st.GetChannel(channels[0]), NotNil)
	c.Check(st.GetChannel(channels[1]), NotNil)
}

func (s *s) TestStore_GetUsersChannelModes(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	st.addUser(users[0])
	c.Check(st.GetUsersChannelModes(users[0], channels[0]), IsNil)
	st.addChannel(channels[0])
	c.Check(st.GetUsersChannelModes(users[0], channels[0]), IsNil)

	st.addToChannel(users[0], channels[0])
	c.Check(st.GetUsersChannelModes(users[0], channels[0]), NotNil)
}

func (s *s) TestStore_GetNUsers(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	c.Check(st.GetNUsers(), Equals, 0)
	st.addUser(users[0])
	st.addUser(users[0]) // Test that adding a user twice does nothing.
	c.Check(st.GetNUsers(), Equals, 1)
	st.addUser(users[1])
	c.Check(st.GetNUsers(), Equals, 2)
}

func (s *s) TestStore_GetNChannels(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	c.Check(st.GetNChannels(), Equals, 0)
	st.addChannel(channels[0])
	st.addChannel(channels[0]) // Test that adding a channel twice does nothing.
	c.Check(st.GetNChannels(), Equals, 1)
	st.addChannel(channels[1])
	c.Check(st.GetNChannels(), Equals, 2)
}

func (s *s) TestStore_GetNUserChans(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	c.Check(st.GetNUserChans(users[0]), Equals, 0)
	c.Check(st.GetNUserChans(users[0]), Equals, 0)
	st.addChannel(channels[0])
	st.addChannel(channels[1])
	c.Check(st.GetNUserChans(users[0]), Equals, 0)
	c.Check(st.GetNUserChans(users[0]), Equals, 0)
	st.addUser(users[0])
	st.addUser(users[1])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[0], channels[0]) // Test no duplicate adds.
	st.addToChannel(users[0], channels[1])
	st.addToChannel(users[1], channels[0])
	c.Check(st.GetNUserChans(users[0]), Equals, 2)
	c.Check(st.GetNUserChans(users[1]), Equals, 1)
}

func (s *s) TestStore_GetNChanUsers(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	c.Check(st.GetNChanUsers(channels[0]), Equals, 0)
	c.Check(st.GetNChanUsers(channels[0]), Equals, 0)
	st.addChannel(channels[0])
	st.addChannel(channels[1])
	c.Check(st.GetNChanUsers(channels[0]), Equals, 0)
	c.Check(st.GetNChanUsers(channels[0]), Equals, 0)
	st.addUser(users[0])
	st.addUser(users[1])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[0], channels[1])
	st.addToChannel(users[1], channels[0])
	c.Check(st.GetNChanUsers(channels[0]), Equals, 2)
	c.Check(st.GetNChanUsers(channels[1]), Equals, 1)
}

func (s *s) TestStore_EachUser(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	st.addUser(users[0])
	st.addUser(users[1])
	i := 0
	st.EachUser(func(u *User) {
		c.Check(users[i], Equals, u.GetFullhost())
		i++
	})
	c.Check(i, Equals, 2)
}

func (s *s) TestStore_EachChannel(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	st.addChannel(channels[0])
	st.addChannel(channels[1])
	i := 0
	st.EachChannel(func(ch *Channel) {
		c.Check(channels[i], Equals, ch.String())
		i++
	})
	c.Check(i, Equals, 2)
}

func (s *s) TestStore_EachUserChan(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	st.addUser(users[0])
	st.addChannel(channels[0])
	st.addChannel(channels[1])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[0], channels[1])
	i := 0
	st.EachUserChan(users[0], func(uc *UserChannel) {
		c.Check(channels[i], Equals, uc.Channel.String())
		i++
	})
	c.Check(i, Equals, 2)
}

func (s *s) TestStore_EachChanUser(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	st.addUser(users[0])
	st.addUser(users[1])
	st.addChannel(channels[0])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[1], channels[0])
	i := 0
	st.EachChanUser(channels[0], func(cu *ChannelUser) {
		c.Check(users[i], Equals, cu.User.GetFullhost())
		i++
	})
	c.Check(i, Equals, 2)
}

func (s *s) TestStore_GetUsers(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	st.addUser(users[0])
	st.addUser(users[1])
	c.Check(len(st.GetUsers()), Equals, 2)
	for i, user := range st.GetUsers() {
		c.Check(users[i], Equals, user)
	}
}

func (s *s) TestStore_GetChannels(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	st.addChannel(channels[0])
	st.addChannel(channels[1])
	c.Check(len(st.GetChannels()), Equals, 2)
	for i, channel := range st.GetChannels() {
		c.Check(channels[i], Equals, channel)
	}
}

func (s *s) TestStore_GetUserChans(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	c.Check(st.GetUserChans(users[0]), IsNil)
	st.addUser(users[0])
	st.addChannel(channels[0])
	st.addChannel(channels[1])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[0], channels[1])
	c.Check(len(st.GetUserChans(users[0])), Equals, 2)
	for i, channel := range st.GetUserChans(users[0]) {
		c.Check(channels[i], Equals, channel)
	}
}

func (s *s) TestStore_GetChanUsers(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	c.Check(st.GetChanUsers(channels[0]), IsNil)
	st.addUser(users[0])
	st.addUser(users[1])
	st.addChannel(channels[0])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[1], channels[0])
	c.Check(len(st.GetChanUsers(channels[0])), Equals, 2)
	for i, user := range st.GetChanUsers(channels[0]) {
		c.Check(users[i], Equals, user)
	}
}

func (s *s) TestStore_IsOn(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	st.addChannel(channels[0])
	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	st.addUser(users[0])
	st.addToChannel(users[0], channels[0])
	c.Check(st.IsOn(users[0], channels[0]), Equals, true)
}

func (s *s) TestStore_UpdateNick(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.NICK,
		Sender: users[0],
		Args:   []string{nicks[1]},
	}

	st.addUser(users[0])
	st.addChannel(channels[0])
	st.addToChannel(users[0], channels[0])

	c.Check(st.GetUser(users[0]), NotNil)
	c.Check(st.GetUser(users[1]), IsNil)
	c.Check(st.IsOn(users[0], channels[0]), Equals, true)
	c.Check(st.IsOn(users[1], channels[0]), Equals, false)

	st.Update(m)

	c.Check(st.GetUser(users[0]), IsNil)
	c.Check(st.GetUser(users[1]), NotNil)
	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	c.Check(st.IsOn(users[1], channels[0]), Equals, true)

	m.Sender = users[0]
	m.Args = []string{"newnick"}
	st.Update(m)
	c.Check(st.GetUser("newnick"), NotNil)
	c.Check(st.GetUser(nicks[0]), IsNil)
}

func (s *s) TestStore_UpdateJoin(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.JOIN,
		Sender: users[0],
		Args:   []string{channels[0]},
	}

	st.addChannel(channels[0])
	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	st.Update(m)
	c.Check(st.IsOn(users[0], channels[0]), Equals, true)

	st, _ = CreateStore(irc.CreateProtoCaps())
	st.Self = self
	st.addChannel(channels[0])

	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	st.Update(m)
	c.Check(st.IsOn(users[0], channels[0]), Equals, true)
}

func (s *s) TestStore_UpdateJoinSelf(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.JOIN,
		Sender: string(self.mask),
		Args:   []string{channels[0]},
	}

	c.Check(st.GetChannel(channels[0]), IsNil)
	c.Check(st.IsOn(st.Self.GetNick(), channels[0]), Equals, false)
	st.Update(m)
	c.Check(st.GetChannel(channels[0]), NotNil)
	c.Check(st.IsOn(st.Self.GetNick(), channels[0]), Equals, true)
}

func (s *s) TestStore_UpdatePart(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)

	m := &irc.IrcMessage{
		Name:   irc.PART,
		Sender: users[0],
		Args:   []string{channels[0]},
	}

	st.addUser(users[0])
	st.addUser(users[1])

	// Test coverage, make sure adding to a channel that doesn't exist does
	// nothing.
	st.addToChannel(users[0], channels[0])
	c.Check(st.IsOn(users[0], channels[0]), Equals, false)

	st.addChannel(channels[0])
	st.addChannel(channels[1])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[1], channels[0])
	st.addToChannel(users[0], channels[1])

	c.Check(st.IsOn(users[0], channels[0]), Equals, true)
	c.Check(st.IsOn(users[1], channels[0]), Equals, true)
	c.Check(st.IsOn(users[0], channels[1]), Equals, true)
	c.Check(st.IsOn(users[1], channels[1]), Equals, false)

	st.Update(m)
	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	c.Check(st.IsOn(users[1], channels[0]), Equals, true)
	c.Check(st.IsOn(users[0], channels[1]), Equals, true)
	c.Check(st.IsOn(users[1], channels[1]), Equals, false)

	m.Sender = users[1]
	st.Update(m)
	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	c.Check(st.IsOn(users[1], channels[0]), Equals, false)
	c.Check(st.IsOn(users[0], channels[1]), Equals, true)
	c.Check(st.IsOn(users[1], channels[1]), Equals, false)

	m.Sender = users[0]
	m.Args[0] = channels[1]
	st.Update(m)

	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	c.Check(st.IsOn(users[1], channels[0]), Equals, false)
	c.Check(st.IsOn(users[0], channels[1]), Equals, false)
	c.Check(st.IsOn(users[1], channels[1]), Equals, false)
}

func (s *s) TestStore_UpdatePartSelf(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.PART,
		Sender: string(self.mask),
		Args:   []string{channels[0]},
	}

	st.addUser(users[0])
	st.addUser(self.GetFullhost())
	st.addChannel(channels[0])
	st.addChannel(channels[1])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[0], channels[1])
	st.addToChannel(self.GetNick(), channels[0])

	c.Check(st.IsOn(users[0], channels[0]), Equals, true)
	c.Check(st.IsOn(users[0], channels[1]), Equals, true)
	c.Check(st.IsOn(self.GetNick(), channels[0]), Equals, true)
	st.Update(m)
	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	c.Check(st.IsOn(users[0], channels[1]), Equals, true)
	c.Check(st.IsOn(self.GetNick(), channels[0]), Equals, false)
}

func (s *s) TestStore_UpdateQuit(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.QUIT,
		Sender: users[0],
		Args:   []string{"quit message"},
	}

	// Test Quitting when we don't know the user
	st.Update(m)
	c.Check(st.GetUser(users[0]), IsNil)

	st.addUser(users[0])
	st.addUser(users[1])
	st.addChannel(channels[0])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[1], channels[0])

	c.Check(st.IsOn(users[0], channels[0]), Equals, true)
	c.Check(st.GetUser(users[0]), NotNil)
	c.Check(st.IsOn(users[1], channels[0]), Equals, true)
	c.Check(st.GetUser(users[1]), NotNil)

	st.Update(m)

	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
	c.Check(st.GetUser(users[0]), IsNil)
	c.Check(st.IsOn(users[1], channels[0]), Equals, true)
	c.Check(st.GetUser(users[1]), NotNil)

	m.Sender = users[1]
	st.Update(m)

	c.Check(st.IsOn(users[1], channels[0]), Equals, false)
	c.Check(st.GetUser(users[1]), IsNil)
}

func (s *s) TestStore_UpdateKick(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.KICK,
		Sender: users[1],
		Args:   []string{channels[0], users[0]},
	}

	st.addUser(users[0])
	st.addUser(users[1])

	st.addChannel(channels[0])
	st.addToChannel(users[0], channels[0])

	c.Check(st.IsOn(users[0], channels[0]), Equals, true)
	st.Update(m)
	c.Check(st.IsOn(users[0], channels[0]), Equals, false)
}

func (s *s) TestStore_UpdateKickSelf(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.KICK,
		Sender: users[1],
		Args:   []string{channels[0], st.Self.GetNick()},
	}

	st.addUser(st.Self.GetFullhost())
	st.addChannel(channels[0])
	st.addToChannel(users[0], channels[0])

	c.Check(st.GetChannel(channels[0]), NotNil)
	st.Update(m)
	c.Check(st.GetChannel(channels[0]), IsNil)
}

func (s *s) TestStore_UpdateMode(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)

	m := &irc.IrcMessage{
		Name:   irc.MODE,
		Sender: users[0],
		Args: []string{channels[0],
			"+ovmb-vn", nicks[0], nicks[0], "*!*mask", nicks[1],
		},
	}

	fail := st.GetUsersChannelModes(users[0], channels[0])
	c.Check(fail, IsNil)

	st.addChannel(channels[0])
	st.addUser(users[0])
	st.addUser(users[1])
	st.addToChannel(users[0], channels[0])
	st.addToChannel(users[1], channels[0])

	u1modes := st.GetUsersChannelModes(users[0], channels[0])
	u2modes := st.GetUsersChannelModes(users[1], channels[0])
	u2modes.SetMode('v')
	st.GetChannel(channels[0]).Set("n")

	c.Check(st.GetChannel(channels[0]).IsSet("n"), Equals, true)
	c.Check(st.GetChannel(channels[0]).IsSet("mb"), Equals, false)
	c.Check(u1modes.HasMode('o'), Equals, false)
	c.Check(u1modes.HasMode('v'), Equals, false)
	c.Check(u2modes.HasMode('v'), Equals, true)
	st.Update(m)
	c.Check(st.GetChannel(channels[0]).IsSet("n"), Equals, false)
	c.Check(st.GetChannel(channels[0]).IsSet("mb *!*mask"), Equals, true)
	c.Check(u1modes.HasMode('o'), Equals, true)
	c.Check(u1modes.HasMode('v'), Equals, true)
	c.Check(u2modes.HasMode('v'), Equals, false)
}

func (s *s) TestStore_UpdateModeSelf(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self.User = self.User
	c.Check(err, IsNil)

	m := &irc.IrcMessage{
		Name:   irc.MODE,
		Sender: self.GetFullhost(),
		Args:   []string{self.GetNick(), "+i-o"},
	}

	st.Self.Set("o")

	c.Check(st.Self.IsSet("i"), Equals, false)
	c.Check(st.Self.IsSet("o"), Equals, true)
	st.Update(m)
	c.Check(st.Self.IsSet("i"), Equals, true)
	c.Check(st.Self.IsSet("o"), Equals, false)
}

func (s *s) TestStore_UpdateTopic(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)

	m := &irc.IrcMessage{
		Name:   irc.TOPIC,
		Sender: users[1],
		Args:   []string{channels[0], "topic topic"},
	}

	st.addChannel(channels[0])

	c.Check(st.GetChannel(channels[0]).GetTopic(), Equals, "")
	st.Update(m)
	c.Check(st.GetChannel(channels[0]).GetTopic(), Equals, "topic topic")
}

func (s *s) TestStore_UpdateRplTopic(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)

	m := &irc.IrcMessage{
		Name:   irc.RPL_TOPIC,
		Sender: server,
		Args:   []string{self.GetNick(), channels[0], "topic topic"},
	}

	st.addChannel(channels[0])

	c.Check(st.GetChannel(channels[0]).GetTopic(), Equals, "")
	st.Update(m)
	c.Check(st.GetChannel(channels[0]).GetTopic(), Equals, "topic topic")
}

func (s *s) TestStore_UpdatePrivmsg(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.PRIVMSG,
		Sender: users[0],
		Args:   []string{channels[0]},
	}

	st.addChannel(channels[0])

	c.Check(st.GetUser(users[0]), IsNil)
	c.Check(st.GetUsersChannelModes(users[0], channels[0]), IsNil)
	st.Update(m)
	c.Check(st.GetUser(users[0]), NotNil)
	c.Check(st.GetUsersChannelModes(users[0], channels[0]), NotNil)

	m.Sender = server
	size := len(st.users)
	st.Update(m)
	c.Check(len(st.users), Equals, size)
}

func (s *s) TestStore_UpdateNotice(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.NOTICE,
		Sender: users[0],
		Args:   []string{channels[0]},
	}

	st.addChannel(channels[0])

	c.Check(st.GetUser(users[0]), IsNil)
	c.Check(st.GetUsersChannelModes(users[0], channels[0]), IsNil)
	st.Update(m)
	c.Check(st.GetUser(users[0]), NotNil)
	c.Check(st.GetUsersChannelModes(users[0], channels[0]), NotNil)

	m.Sender = server
	size := len(st.users)
	st.Update(m)
	c.Check(len(st.users), Equals, size)
}

func (s *s) TestStore_UpdateWelcome(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)
	m := &irc.IrcMessage{
		Name:   irc.RPL_WELCOME,
		Sender: server,
		Args:   []string{nicks[1], "Welcome to"},
	}

	st.Update(m)
	c.Check(st.Self.GetFullhost(), Equals, nicks[1])
	c.Check(st.users[nicks[1]].GetFullhost(), Equals, st.Self.GetFullhost())

	m = &irc.IrcMessage{
		Name:   irc.RPL_WELCOME,
		Sender: server,
		Args:   []string{nicks[1], "Welcome to " + users[1]},
	}

	st.Update(m)
	c.Check(st.Self.GetFullhost(), Equals, users[1])
	c.Check(st.users[nicks[1]].GetFullhost(), Equals, st.Self.GetFullhost())
}

func (s *s) TestStore_UpdateRplNamereply(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)

	m := &irc.IrcMessage{
		Name:   irc.RPL_NAMREPLY,
		Sender: server,
		Args: []string{
			self.GetNick(), "=", channels[0],
			"@" + nicks[0] + " +" + nicks[1] + " " + self.GetNick(),
		},
	}

	st.addChannel(channels[0])

	c.Check(st.GetUsersChannelModes(users[0], channels[0]), IsNil)
	c.Check(st.GetUsersChannelModes(users[1], channels[0]), IsNil)
	c.Check(st.GetUsersChannelModes(self.GetNick(), channels[0]), IsNil)
	st.Update(m)
	c.Check(
		st.GetUsersChannelModes(users[0], channels[0]).String(), Equals, "o")
	c.Check(
		st.GetUsersChannelModes(users[1], channels[0]).String(), Equals, "v")
	c.Check(st.GetUsersChannelModes(
		self.GetNick(), channels[0]).String(), Equals, "")
}

func (s *s) TestStore_RplWhoReply(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)

	m := &irc.IrcMessage{
		Name:   irc.RPL_WHOREPLY,
		Sender: server,
		Args: []string{
			self.GetNick(), channels[0], irc.Mask(users[0]).GetUsername(),
			irc.Mask(users[0]).GetHost(), "*.server.net", nicks[0], "Hx@d",
			"3 real name",
		},
	}

	st.addChannel(channels[0])

	c.Check(st.GetUser(users[0]), IsNil)
	c.Check(st.GetUsersChannelModes(users[0], channels[0]), IsNil)
	st.Update(m)
	c.Check(st.GetUser(users[0]), NotNil)
	c.Check(st.GetUser(users[0]).GetFullhost(), Equals, users[0])
	c.Check(st.GetUser(users[0]).GetRealname(), Equals, "real name")
	c.Check(
		st.GetUsersChannelModes(users[0], channels[0]).String(), Equals, "o")
}

func (s *s) TestStore_UpdateRplMode(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	c.Check(err, IsNil)

	m := &irc.IrcMessage{
		Name:   irc.RPL_CHANNELMODEIS,
		Sender: server,
		Args:   []string{self.GetNick(), channels[0], "+ntzl", "10"},
	}

	st.addChannel(channels[0])
	c.Check(st.GetChannel(channels[0]).IsSet("ntzl 10"), Equals, false)
	st.Update(m)
	c.Check(st.GetChannel(channels[0]).IsSet("ntzl 10"), Equals, true)
}

func (s *s) TestStore_UpdateRplBanlist(c *C) {
	st, err := CreateStore(irc.CreateProtoCaps())
	st.Self = self
	c.Check(err, IsNil)

	m := &irc.IrcMessage{
		Name:   irc.RPL_BANLIST,
		Sender: server,
		Args: []string{self.GetNick(), channels[0], nicks[0] + "!*@*", nicks[1],
			"1367197165"},
	}

	st.addChannel(channels[0])
	c.Check(st.GetChannel(channels[0]).HasBan(nicks[0]+"!*@*"), Equals, false)
	st.Update(m)
	c.Check(st.GetChannel(channels[0]).HasBan(nicks[0]+"!*@*"), Equals, true)
}
