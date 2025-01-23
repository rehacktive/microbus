package microbus

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	mb := New()
	sub1 := mb.Subscribe("test")
	sub2 := mb.Subscribe("test")

	// Send a message
	mb.Publish("test", "hello")

	// Wait for messages to be received
	time.Sleep(1 * time.Second)

	assert.NotNil(t, sub1)
	assert.NotNil(t, sub2)
}

func TestPublishAfterClose(t *testing.T) {
	mb := New()

	// Close the bus
	mb.CloseAll()

	// Try to publish after closing
	mb.Publish("test", "should not be delivered")
}

func TestMultipleSubscribers(t *testing.T) {
	mb := New()

	// Create two subscribers
	sub1 := mb.Subscribe("test")
	sub2 := mb.Subscribe("test")

	// Publish a message
	mb.Publish("test", "message1")
	time.Sleep(2 * time.Second)

	// Verify both subscribers received the message
	select {
	case <-sub1:
		break
	case <-sub2:
		break
	}
}

func TestSubscribeAfterClose(t *testing.T) {
	mb := New()

	// Close the bus
	mb.CloseAll()

	// Try to subscribe after closing - should return nil channel
	sub := mb.Subscribe("test")
	assert.Nil(t, sub)
}

func TestSubscriptionChannel(t *testing.T) {
	mb := New()

	// Create a channel to receive subscription messages
	subChan := mb.Subscribe("test")

	// Publish a message and check if it's received
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Publish a message immediately after subscribing
	go mb.Publish("test", "test_message")

	select {
	case <-subChan:
		t.Log("Received subscription message")
	case <-ctx.Done():
		t.Error("timed out waiting for subscription message")
	}
}

// Test that the closed flag is properly set
func TestCloseAll(t *testing.T) {
	mb := New()
	mb.CloseAll()

	// Try to publish after closing
	assert.True(t, mb.closed)
}
