package irc

import (
	"bytes"
	"fmt"
	. "launchpad.net/gocheck"
	"strings"
	"testing"
)

func Test(t *testing.T) { TestingT(t) } //Hook into testing package

type s struct{}

var _ = Suite(&s{})

func (s *s) TestIrcMessage_Test(c *C) {
	args := []string{"#chan1", "#chan2"}
	msg := IrcMessage{
		Args: []string{strings.Join(args, ",")},
	}
	for i, v := range msg.Split(0) {
		c.Check(args[i], Equals, v)
	}
}

func (s *s) TestMsgTypes_Privmsg(c *C) {
	args := []string{"#chan", "msg arg"}
	pmsg := &Message{&IrcMessage{
		Name:   PRIVMSG,
		Args:   args,
		Sender: "user@host.com",
	}}

	c.Check(pmsg.Target(), Equals, args[0])
	c.Check(pmsg.Message(), Equals, args[1])
}

func (s *s) TestMsgTypes_Notice(c *C) {
	args := []string{"#chan", "msg arg"}
	notice := &Message{&IrcMessage{
		Name:   NOTICE,
		Args:   args,
		Sender: "user@host.com",
	}}

	c.Check(notice.Target(), Equals, args[0])
	c.Check(notice.Message(), Equals, args[1])
}

type fakeHelper struct {
	*Helper
}

func (f *fakeHelper) GetKey() string {
	return ""
}

func (s *s) TestHelper_Send(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	format := "PRIVMSG %v :%v"
	target := "#chan"
	msg := "msg"
	h.Send(format, target, msg)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprint(format, target, msg))
}

func (s *s) TestHelper_Sendln(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	header := "PRIVMSG"
	target := "#chan"
	msg := "msg"
	h.Sendln(header, target, msg)
	expect := fmt.Sprintln(header, target, msg)
	c.Check(string(buf.Bytes()), Equals, expect[:len(expect)-1])
}

func (s *s) TestHelper_Sendf(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	format := "PRIVMSG %v :%v"
	target := "#chan"
	msg := "msg"
	h.Sendf(format, target, msg)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf(format, target, msg))
}

func (s *s) TestHelper_Privmsg(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	ch := "#chan"
	s1, s2 := "string1", "string2"
	h.Privmsg(ch, s1, s2)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v %v :%v",
		PRIVMSG, ch, fmt.Sprint(s1, s2)))
}

func (s *s) TestHelper_Privmsgln(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	ch := "#chan"
	s1, s2 := "string1", "string2"
	expect := fmt.Sprintln(s1, s2)
	h.Privmsgln(ch, s1, s2)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v %v :%v",
		PRIVMSG, ch, expect[:len(expect)-1]))
}

func (s *s) TestHelper_Privmsgf(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	ch := "#chan"
	format := "%v - %v"
	s1, s2 := "string1", "string2"
	h.Privmsgf(ch, format, s1, s2)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v %v :%v",
		PRIVMSG, ch, fmt.Sprintf(format, s1, s2)))
}

func (s *s) TestHelper_Notice(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	ch := "#chan"
	s1, s2 := "string1", "string2"
	h.Notice(ch, s1, s2)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v %v :%v",
		NOTICE, ch, fmt.Sprint(s1, s2)))
}

func (s *s) TestHelper_Noticeln(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	ch := "#chan"
	s1, s2 := "string1", "string2"
	expect := fmt.Sprintln(s1, s2)
	h.Noticeln(ch, s1, s2)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v %v :%v",
		NOTICE, ch, expect[:len(expect)-1]))
}

func (s *s) TestHelper_Noticef(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	ch := "#chan"
	format := "%v - %v"
	s1, s2 := "string1", "string2"
	h.Noticef(ch, format, s1, s2)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v %v :%v",
		NOTICE, ch, fmt.Sprintf(format, s1, s2)))
}

func (s *s) TestHelper_Join(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	ch := "#chan"
	h.Join()
	c.Check(buf.Len(), Equals, 0)
	h.Join(ch)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v :%v", JOIN, ch))

	buf = bytes.Buffer{}
	h.Writer = &buf
	h.Join(ch, ch)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v :%v,%v", JOIN, ch, ch))
}

func (s *s) TestHelper_Part(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	ch := "#chan"
	h.Part()
	c.Check(buf.Len(), Equals, 0)
	h.Part(ch)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v :%v", PART, ch))

	buf = bytes.Buffer{}
	h.Writer = &buf
	h.Part(ch, ch)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v :%v,%v", PART, ch, ch))
}

func (s *s) TestHelper_Quit(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	msg := "quitting"
	h.Quit(msg)
	c.Check(string(buf.Bytes()), Equals, fmt.Sprintf("%v :%v", QUIT, msg))
}

func (s *s) TestHelper_splitSend(c *C) {
	buf := bytes.Buffer{}
	h := &Helper{&buf}
	header := "PRIVMSG #chan :"
	s0 := "message"
	h.splitSend([]byte(header), []byte(s0))
	c.Check(buf.Len(), Equals, len(header)+len(s0))

	buf = bytes.Buffer{}
	h = &Helper{&buf}
	header = "PRIVMSG #chan :"
	s1 := strings.Repeat("a", 510)
	s2 := strings.Repeat("b", 510)
	s3 := strings.Repeat("c", 200)
	err := h.splitSend([]byte(header), []byte(s1+s2+s3))
	c.Check(err, IsNil)
	c.Check(buf.Len(), Equals, len(header)*3+len(s1)+len(s2)+len(s3))
}
