package redeye

import (
	"testing"
	
)

func TestMessanger(t *testing.T) {
	broker := "tcp://localhost:1883"
	tpath  := "/test/path"
	m := NewMessanger(broker, tpath)
	if (m.Broker != broker) {
		t.Error("Expected broker: (", broker, ") got (", m.Broker, ")")
	}

	if (m.BasePath != tpath) {
		t.Error("Expected BasePath: (", tpath, ") got (", m.BasePath, ")")
	}

	if (m.Subscriptions != nil) {
		t.Errorf("Subscriptions expected (0) got (%d)", len(m.Subscriptions))		
	}

	q := m.Start()
	if q == nil {
		t.Error("Expected a channel for commands but nil returned")
	}
	
	m.Subscribe("/pinkeye/test")
	if (len(m.Subscriptions) != 1) {
		t.Errorf("Subscriptions expected (1) got (%d)", len(m.Subscriptions))		
	}
}
