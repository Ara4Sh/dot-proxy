package dnsproxy

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestNewDNSProxy(t *testing.T) {
	expectedEndpoint := "1.1.1.1:853"
	proxy := NewDNSProxy(expectedEndpoint)
	if proxy.cloudFlareDotEndpoint != expectedEndpoint {
		t.Errorf("Expected CloudFlareDotEndpoint to be %s, got %s", expectedEndpoint, proxy.cloudFlareDotEndpoint)
	}
}

func TestPrepareQuery(t *testing.T) {
	proxy := NewDNSProxy("1.1.1.1:853")
	query := []byte{0x00, 0x01, 0x02}
	preparedQuery := proxy.prepareQuery(query)
	if len(preparedQuery) != len(query)+2 {
		t.Errorf("Expected prepared query length to be %d, got %d", len(query)+2, len(preparedQuery))
	}
	if !bytes.Equal(preparedQuery[2:], query) {
		t.Errorf("Expected prepared query to contain original query bytes, got %v", preparedQuery[2:])
	}
}

func TestListenUDPInvalidAddress(t *testing.T) {
	proxy := NewDNSProxy("1.1.1.1:853")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Using an invalid address
	err := proxy.ListenUDP(ctx, "invalid", 0)
	if err == nil {
		t.Error("Expected ListenUDP to fail due to invalid address, but it did not")
	}
}

func TestListenTCPInvalidAddress(t *testing.T) {
	proxy := NewDNSProxy("1.1.1.1:853")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := proxy.ListenTCP(ctx, "invalid", 0)
	if err == nil {
		t.Error("Expected ListenUDP to fail due to invalid address, but it did not")
	}
}
