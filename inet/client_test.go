package inet

import (
	"bytes"
	"github.com/aarondl/ultimateq/mocks"
	"io"
	. "launchpad.net/gocheck"
	"sync"
	"testing"
	"time"
)

func Test(t *testing.T) { TestingT(t) } //Hook into testing package
type s struct{}

var _ = Suite(&s{})

func init() {
	Log = func(string, ...interface{}) {}
}

func (s *s) TestCreateIrcClient(c *C) {
	conn := mocks.CreateConn()
	client := CreateIrcClient(conn, "name")
	c.Check(client.isShutdown, Equals, false)
	c.Check(client.conn, Equals, conn)
	c.Check(client.siphonchan, NotNil)
	c.Check(client.pumpchan, NotNil)
	c.Check(client.pumpservice, NotNil)
	c.Check(client.name, Equals, "name")
	c.Check(client.lastwrite.Before(time.Now()), Equals, true)
}

func (s *s) TestIrcClient_ImplementsReadWriteCloser(c *C) {
	client := CreateIrcClient(nil, "")
	c.Check(client, FitsTypeOf, io.ReadWriteCloser(client))
}

func (s *s) TestIrcClient_SpawnWorkers(c *C) {
	conn := mocks.CreateConn()

	client := CreateIrcClient(conn, "")
	c.Check(client.IsClosed(), Equals, false)
	client.SpawnWorkers(true, true)
	c.Check(client.IsClosed(), Equals, false)
	conn.Send([]byte{}, 0, io.EOF)
	client.Close()
	conn.WaitForDeath()
}

func (s *s) TestIrcClient_Pump(c *C) {
	test1 := []byte("PONG :arg1 arg2\r\n")
	test2 := []byte("NOTICE :arg1\r\n")
	split := 2

	conn := mocks.CreateConn()
	client := CreateIrcClient(conn, "")

	fakelast := time.Now().Truncate(5 * time.Hour)
	client.SpawnWorkers(true, false)
	ch := <-client.pumpservice
	ch <- []byte{} //Inconsequential, testcov error handling

	go func() {
		client.Write(test1)
		client.Write(test2)
		client.Close()
	}()

	c.Check(bytes.Compare(conn.Receive(split, nil), test1), Equals, 0)
	c.Check(bytes.Compare(conn.Receive(len(test1[split:]), nil),
		test1[split:]), Equals, 0)

	c.Check(bytes.Compare(conn.Receive(len(test2), nil), test2), Equals, 0)
	conn.WaitForDeath()
	conn.ResetDeath()

	//Shameful test coverage
	client = CreateIrcClient(conn, "")
	client.SpawnWorkers(true, false)
	client.Write(test1)
	c.Check(bytes.Compare(conn.Receive(0, io.EOF), test1), Equals, 0)
	client.Close()
	conn.WaitForDeath()
	c.Check(client.lastwrite.Equal(fakelast), Equals, false)
}

func (s *s) TestIrcClient_PumpTimeouts(c *C) {
	test1 := []byte("PRIVMSG :arg1 arg2\r\n")

	conn := mocks.CreateConn()
	client := CreateIrcClientFloodProtect(conn, "", 5, 5, 5, time.Millisecond)
	client.SpawnWorkers(true, false)

	go func() {
		for i := 0; i < 10; i++ {
			_, err := client.Write(test1)
			c.Check(err, IsNil)
		}
	}()

	for i := 0; i < 9; i++ {
		c.Check(bytes.Compare(conn.Receive(len(test1), nil), test1), Equals, 0)
	}
	c.Check(bytes.Compare(conn.Receive(len(test1), io.EOF), test1), Equals, 0)

	client.Close()
	conn.WaitForDeath()
}

func (s *s) TestIrcClient_Siphon(c *C) {
	test1 := []byte("PRIVMSG :msg\r\n")
	test2 := []byte("NOTICE :msg\r\n")
	test3 := []byte("PRIV")

	buf := append(append(append([]byte{}, test1...), test2...), test3...)

	conn := mocks.CreateConn()
	client := CreateIrcClient(conn, "")
	ch := client.ReadChannel()
	client.SpawnWorkers(false, true)

	go func() {
		conn.Send(buf, len(buf), io.EOF)
	}()

	msg := <-ch
	c.Check(bytes.Compare(test1[:len(test1)-2], msg), Equals, 0)
	msg = <-ch
	c.Check(bytes.Compare(test2[:len(test2)-2], msg), Equals, 0)
	_, ok := <-ch
	c.Check(ok, Equals, false)

	client.Close()
	conn.WaitForDeath()
	conn.ResetDeath()

	client = CreateIrcClient(conn, "")
	client.SpawnWorkers(false, true)
	go func() { conn.Send(buf, len(buf), nil) }()
	client.Close()
	conn.WaitForDeath()
}

func (s *s) TestIrcClient_ExtractMessages(c *C) {
	test1 := []byte("irc message 1\r\n")
	test2 := []byte("irc message 2\r\n")
	test3 := []byte("irc mess")
	buf := append(append(append([]byte{}, test1...), test2...), test3...)

	waiter := sync.WaitGroup{}
	waiter.Add(1)

	client := CreateIrcClient(nil, "")
	ret := 0

	go func() {
		var abort bool
		ret, abort = client.extractMessages(buf)
		c.Check(ret, Equals, len(test3))
		c.Check(abort, Equals, false)
		c.Check(bytes.Compare(buf[:ret], test3), Equals, 0)
		waiter.Done()
	}()
	msg1 := <-client.siphonchan
	c.Check(bytes.Compare(msg1, test1[:len(test1)-2]), Equals, 0)
	msg2 := <-client.siphonchan
	c.Check(bytes.Compare(msg2, test2[:len(test2)-2]), Equals, 0)
	waiter.Wait()

	buf = append(buf[:ret], []byte{'\r', '\n'}...)
	waiter.Add(1)
	go func() {
		var abort bool
		ret, abort := client.extractMessages(buf)
		c.Check(ret, Equals, 0)
		c.Check(abort, Equals, false)
		waiter.Done()
	}()
	msg3 := <-client.siphonchan
	c.Check(bytes.Compare(msg3, test3), Equals, 0)
	waiter.Wait()

	waiter.Add(1)
	client.killsiphon = make(chan int)
	go func() {
		_, abort := client.extractMessages(test1)
		c.Check(abort, Equals, true)
		waiter.Done()
	}()
	client.killsiphon <- 0
	waiter.Wait()
}

func (s *s) TestIrcClient_Close(c *C) {
	conn := mocks.CreateConn()

	client := CreateIrcClient(conn, "")

	go func() {
		err := client.Close()
		c.Check(err, IsNil)
	}()
	conn.WaitForDeath()
	c.Check(client.IsClosed(), Equals, true)

	err := client.Close() // Double closing should do nothing
	c.Check(err, IsNil)
	c.Check(client.IsClosed(), Equals, true)
}

func (s *s) TestIrcClient_ReadMessage(c *C) {
	client := CreateIrcClient(nil, "")
	read := []byte("PRIVMSG #chan :msg")
	go func() {
		client.siphonchan <- read
		close(client.siphonchan)
	}()
	msg, ok := client.ReadMessage()
	c.Check(ok, Equals, true)
	c.Check(bytes.Compare(msg, read), Equals, 0)
	msg, ok = client.ReadMessage()
	c.Check(ok, Equals, false)
}

func (s *s) TestIrcClient_Read(c *C) {
	client := CreateIrcClient(nil, "")
	read := []byte("PRIVMSG #chan :msg")
	go func() {
		client.siphonchan <- read
		close(client.siphonchan)
	}()
	buf := make([]byte, len(read))
	breakat := 2

	n, err := client.Read(buf[:breakat])
	c.Check(err, IsNil)
	c.Check(n, Equals, breakat)
	c.Check(bytes.Compare(buf[:breakat], read[:breakat]), Equals, 0)

	n, err = client.Read(buf[breakat:])
	c.Check(err, IsNil)
	c.Check(n, Equals, len(read)-breakat)
	c.Check(bytes.Compare(buf, read), Equals, 0)

	n, err = client.Read(buf)
	c.Check(n, Equals, 0)
	c.Check(err, Equals, io.EOF)
}

func (s *s) TestIrcClient_Write(c *C) {
	test1 := []byte("PRIVMSG #chan :msg\r\n")
	test2 := []byte("PRIVMSG #chan :msg2")

	client := CreateIrcClient(nil, "")
	ch := make(chan []byte)
	go func() {
		arg := append(test1, test2...)
		client.Write(nil) //Should be Consequenceless test cov
		n, err := client.Write(arg)
		c.Check(err, IsNil)
		c.Check(n, Equals, len(arg))
	}()
	client.pumpservice <- ch
	c.Check(bytes.Compare(<-ch, test1), Equals, 0)
	client.pumpservice <- ch
	c.Check(bytes.Compare(<-ch, append(test2, []byte{13, 10}...)), Equals, 0)

	close(client.pumpservice)
	n, err := client.Write(test1)
	c.Check(n, Equals, 0)
	c.Check(err, Equals, io.EOF)
	n, err = client.Write(test2) // Test abortion of no \r\n
	c.Check(n, Equals, 0)
	c.Check(err, Equals, io.EOF)
}

func (s *s) TestIrcClient_calcSleepTime(c *C) {
	var sleep time.Duration
	client := CreateIrcClient(nil, "")

	sleep = client.calcSleepTime(time.Now())
	c.Check(sleep, Equals, time.Duration(0))

	client = CreateIrcClientFloodProtect(nil, "", 5, 5, 5, time.Second)
	client.lastwrite = time.Now()
	for i := 1; i <= 5; i++ {
		sleep = client.calcSleepTime(time.Now())
		c.Check(sleep, Equals, time.Duration(0))
	}

	sleep = client.calcSleepTime(time.Now())
	c.Check(sleep, Not(Equals), time.Duration(0))

	// Check no-sleep and negative cases
	client = CreateIrcClientFloodProtect(nil, "", 5, 5, 5, time.Second)
	client.lastwrite = time.Now()
	sleep = client.calcSleepTime(time.Time{})
	c.Check(sleep, Equals, time.Duration(0))
	sleep = client.calcSleepTime(time.Now().Add(5 * time.Hour))
	c.Check(sleep, Equals, time.Duration(0))

}

func (s *s) TestfindChunks(c *C) {
	test1 := []byte("PRIVMSG #chan :msg\r\n")
	test2 := []byte("NOTICE #chan :msg2\r\n")
	test3 := []byte("PRIV")

	args := append(append(test1, test2...), test3...)
	expected := [][]byte{test1, test2, test3}
	start, remaining, abort := findChunks(args, func(result []byte) bool {
		c.Check(bytes.Compare(result, expected[0]), Equals, 0)
		expected = expected[1:]
		return false
	})
	c.Check(abort, Equals, false)
	c.Check(bytes.Compare(args[start:], test3), Equals, 0)

	start, remaining, abort = findChunks(test1, func(result []byte) bool {
		c.Check(bytes.Compare(test1, result), Equals, 0)
		return false
	})
	c.Check(start, Equals, 0)
	c.Check(abort, Equals, false)
	c.Check(remaining, Equals, false)

	_, _, abort = findChunks(args, func(result []byte) bool {
		return true
	})
	c.Check(abort, Equals, true)
}
