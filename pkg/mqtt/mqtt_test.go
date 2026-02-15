package mqtt

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type fakeToken struct {
	err  error
	done chan struct{}
}

func newFakeToken(err error, complete bool) *fakeToken {
	t := &fakeToken{err: err, done: make(chan struct{})}
	if complete {
		close(t.done)
	}
	return t
}

func (t *fakeToken) Wait() bool {
	<-t.done
	return true
}

func (t *fakeToken) WaitTimeout(d time.Duration) bool {
	select {
	case <-t.done:
		return true
	case <-time.After(d):
		return false
	}
}

func (t *fakeToken) Done() <-chan struct{} {
	return t.done
}

func (t *fakeToken) Error() error {
	return t.err
}

type fakeClient struct {
	connected        bool
	publishToken     pahomqtt.Token
	disconnectCalled bool
}

func (c *fakeClient) IsConnected() bool       { return c.connected }
func (c *fakeClient) IsConnectionOpen() bool  { return c.connected }
func (c *fakeClient) Connect() pahomqtt.Token { return newFakeToken(nil, true) }
func (c *fakeClient) Disconnect(quiesce uint) { c.disconnectCalled = true }
func (c *fakeClient) Publish(topic string, qos byte, retained bool, payload interface{}) pahomqtt.Token {
	return c.publishToken
}
func (c *fakeClient) Subscribe(topic string, qos byte, callback pahomqtt.MessageHandler) pahomqtt.Token {
	return newFakeToken(nil, true)
}
func (c *fakeClient) SubscribeMultiple(filters map[string]byte, callback pahomqtt.MessageHandler) pahomqtt.Token {
	return newFakeToken(nil, true)
}
func (c *fakeClient) Unsubscribe(topics ...string) pahomqtt.Token {
	return newFakeToken(nil, true)
}
func (c *fakeClient) AddRoute(topic string, callback pahomqtt.MessageHandler) {}
func (c *fakeClient) OptionsReader() pahomqtt.ClientOptionsReader {
	return pahomqtt.ClientOptionsReader{}
}

func TestPublisherPublishErrors(t *testing.T) {
	ctx := context.Background()

	var publisher *Publisher
	if err := publisher.Publish(ctx, []byte("data")); err == nil {
		t.Fatal("expected error for nil publisher")
	}

	publisher = &Publisher{client: nil}
	if err := publisher.Publish(ctx, []byte("data")); err == nil {
		t.Fatal("expected error for nil client")
	}

	publisher = &Publisher{client: &fakeClient{connected: false}, logger: zap.NewNop()}
	if err := publisher.Publish(ctx, []byte("data")); err == nil {
		t.Fatal("expected error for disconnected client")
	}
}

func TestPublisherPublishTokenError(t *testing.T) {
	ctx := context.Background()
	fake := &fakeClient{connected: true, publishToken: newFakeToken(errors.New("publish failed"), true)}
	publisher := &Publisher{client: fake, topic: "test", qos: 0, logger: zap.NewNop()}

	err := publisher.Publish(ctx, []byte("data"))
	if err == nil || !strings.Contains(err.Error(), "mqtt publish") {
		t.Fatalf("expected publish error, got %v", err)
	}
}

func TestPublisherPublishContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fake := &fakeClient{connected: true, publishToken: newFakeToken(nil, false)}
	publisher := &Publisher{client: fake, topic: "test", qos: 0, logger: zap.NewNop()}

	err := publisher.Publish(ctx, []byte("data"))
	if err == nil || !strings.Contains(err.Error(), "context cancelled") {
		t.Fatalf("expected context cancel error, got %v", err)
	}
}

func TestPublisherClose(t *testing.T) {
	var publisher *Publisher
	if err := publisher.Close(); err != nil {
		t.Fatalf("Close nil publisher: %v", err)
	}

	fake := &fakeClient{connected: true, publishToken: newFakeToken(nil, true)}
	publisher = &Publisher{client: fake, topic: "test", qos: 0, logger: zap.NewNop()}
	if err := publisher.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if !fake.disconnectCalled {
		t.Fatal("expected Disconnect to be called")
	}
}
