package tpc

import (
	"reflect"
	"testing"
	"time"
)

func TestCannotCloseBeforeConnecting(t *testing.T) {
	c := New()
	err := c.Close()
	if err != ErrAlreadyDisconneced {
		t.Fatal("wrong error value")
	}
}

func TestConnectAndReceiveInsertAck(t *testing.T) {
	const uid1 = 838959385874432
	var uids1 = []uint64{uid1}
	c := New()
	ch := make(chan []uint64)
	c.Start()
	c.InsertSubscriptions(uid1)
	c.OnAck(func(userIDs ...uint64) {
		ch <- userIDs
	})
	select {
	case userIDs := <-ch:
		if reflect.DeepEqual(userIDs, uids1) {
			c.Close()
		}
	case <-time.After(2 * time.Second):
		t.Fatal("didn't receive sub ack")
	}
	<-time.After(time.Second)
}
