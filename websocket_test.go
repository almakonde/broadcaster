package broadcaster

import "testing"

func TestWSConnect(t *testing.T) {
	server, err := startServer(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	client, err := newWSClient(server)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Authenticate(nil)
	if err != nil {
		t.Fatal(err)
	}

	stats, err := server.Broadcaster.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Connections != 1 {
		t.Errorf("Unexpected connection count: %d", stats.Connections)
	}
}

func TestWSCanConnect(t *testing.T) {
	server, err := startServer(&Server{
		CanConnect: func(data map[string]string) bool {
			return false
		},
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	client, err := newWSClient(server)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Authenticate(nil)
	if err == nil {
		t.Fatal("Expected error!")
	}
	if err.Error() != "websocket: close 401 Unauthorized" {
		t.Fatal("Did not properly deny access")
	}

	stats, err := server.Broadcaster.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Connections != 0 {
		t.Errorf("Unexpected connection count: %d", stats.Connections)
	}
}

func TestWSRefusesUnauthedCommands(t *testing.T) {
	server, err := startServer(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	client, err := newWSClient(server)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Send("bla", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Receive()
	if err == nil {
		t.Fatal("Expected error!")
	}
	if err.Error() != "websocket: close 401 Auth expected" {
		t.Fatal("Did not properly deny access")
	}

	stats, err := server.Broadcaster.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Connections != 0 {
		t.Errorf("Unexpected connection count: %d", stats.Connections)
	}
}

func TestWSSubscribe(t *testing.T) {
	server, err := startServer(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	client, err := newWSClient(server)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Authenticate(nil)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Subscribe("test")
	if err != nil {
		t.Fatal(err)
	}

	stats, err := server.Broadcaster.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Connections != 1 {
		t.Errorf("Unexpected connection count: %d", stats.Connections)
	}
	if stats.localSubscriptions["test"] != 1 {
		t.Errorf("Unexpected subscription count: %d", stats.localSubscriptions["test"])
	}
}

func TestWSCanSubscribe(t *testing.T) {
	server, err := startServer(&Server{
		CanSubscribe: func(data map[string]string, channel string) bool {
			return false
		},
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	client, err := newWSClient(server)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Authenticate(nil)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Subscribe("test")
	if err == nil {
		t.Fatal("Expected error!")
	}
	if err.Error() != "websocket: close 403 Channel refused" {
		t.Fatal("Did not properly deny access")
	}

	stats, err := server.Broadcaster.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Connections != 0 {
		// Should disconnect bad client
		t.Errorf("Unexpected connection count: %d", stats.Connections)
	}
	if stats.localSubscriptions["test"] != 0 {
		t.Errorf("Unexpected subscription count: %d", stats.localSubscriptions["test"])
	}
}

func TestWSMessageTypes(t *testing.T) {
	server, err := startServer(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	client, err := newWSClient(server)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Authenticate(nil)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Send("bla", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Receive()
	if err == nil {
		t.Fatal("Expected error!")
	}
	if err.Error() != "websocket: close 400 Unexpected message" {
		t.Fatal("Did not properly refuse message type")
	}

	stats, err := server.Broadcaster.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Connections != 0 {
		t.Errorf("Unexpected connection count: %d", stats.Connections)
	}
	if stats.localSubscriptions["test"] != 0 {
		t.Errorf("Unexpected subscription count: %d", stats.localSubscriptions["test"])
	}
}