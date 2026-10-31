package network

import (
	"context"
	"testing"
	"time"
)

// The proof must agree between the two ends of one session and differ between
// sessions — that difference is what stops a man in the middle, who holds two
// separate sessions, from relaying a client's proof to the server.
func TestQUICProofIsBoundToTheSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	l, err := QUICListen("127.0.0.1:0", QUICSettings{})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()

	type pair struct{ client, server string }
	connect := func() pair {
		t.Helper()
		cc, err := QUICDial(ctx, l.Addr().String(), QUICSettings{})
		if err != nil {
			t.Fatalf("dial: %v", err)
		}
		defer cc.CloseWithError(0, "")
		sc, err := l.Accept(ctx)
		if err != nil {
			t.Fatalf("accept: %v", err)
		}
		defer sc.CloseWithError(0, "")
		c, err := QUICProof(cc, "secret", "client")
		if err != nil {
			t.Fatalf("client proof: %v", err)
		}
		s, err := QUICProof(sc, "secret", "client")
		if err != nil {
			t.Fatalf("server-side proof: %v", err)
		}
		return pair{c, s}
	}

	a, b := connect(), connect()
	if a.client != a.server || b.client != b.server {
		t.Fatal("the two ends of one session computed different proofs")
	}
	if a.client == b.client {
		t.Fatal("two sessions produced the same proof; it is not bound to the session")
	}
}
