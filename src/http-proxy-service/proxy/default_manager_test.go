package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewDefaultManager(t *testing.T) {
	config := &Config{
		ListenAddress: "localhost:8080",
		ProxyRoutes: []Route{
			{
				Name:    "TestRoute",
				Context: "/test",
				Target:  "http://example.com",
			},
		},
	}

	// Create a new manager
	manager := NewDefaultManager(config)

	if manager == nil {
		t.Error("Expected a non-nil manager, got nil")
	}

	// You can write more specific tests based on your requirements.
}

func TestDefaultManager_GetProxyRouter(t *testing.T) {
	config := &Config{
		ListenAddress: "localhost:8080",
		ProxyRoutes: []Route{
			{
				Name:    "TestRoute",
				Context: "/test",
				Target:  "http://example.com",
			},
		},
	}
	manager := NewDefaultManager(config)

	router := manager.GetProxyRouter()

	if router == nil {
		t.Error("Expected a non-nil router, got nil")
	}
}

func TestProxyServer(t *testing.T) {
	want1 := "I'm service 1"
	want2 := "I'm service 2"

	testService1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(want1))
	}))
	defer testService1.Close()
	testService2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(want2))
	}))
	defer testService2.Close()

	// Create a sample configuration for testing
	config := &Config{
		ListenAddress: "localhost:8080",
		ProxyRoutes: []Route{
			{
				Name:    "Route1",
				Context: "/context1",
				Target:  testService1.URL,
			},
			{
				Name:    "Route2",
				Context: "/context2",
				Target:  testService2.URL,
			},
		},
	}

	// Start a test HTTP server using the NewDefaultManager
	proxyManager := NewDefaultManager(config)
	testProxyServer := httptest.NewServer(proxyManager.GetProxyRouter())
	defer testProxyServer.Close()

	// Test an HTTP request to Route1
	resp, err := testProxyServer.Client().Get(testProxyServer.URL + "/context1/")
	if err != nil {
		t.Errorf("HTTP GET to /context1/test failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// Check if the response body matches the expected value for Route1
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	if string(body) != want1 {
		t.Errorf("Expected response body for Route1 to be %s, but got %s", want1, string(body))
	}

	// Test an HTTP request to Route2
	resp, err = testProxyServer.Client().Get(testProxyServer.URL + "/context2/")
	if err != nil {
		t.Errorf("HTTP GET to /context2/test failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// Check if the response body matches the expected value for Route2
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	if string(body) != want2 {
		t.Errorf("Expected response body for Route2 to be %s, but got %s", want2, string(body))
	}
}

/*func TestDefaultManagerNewHandler(t *testing.T) {
	type fields struct {
		proxies []*httputil.ReverseProxy
		routing *router.Router
	}
	type args struct {
		p *httputil.ReverseProxy
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(http.ResponseWriter, *http.Request)
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := defaultManager{
				proxies: tt.fields.proxies,
				routing: tt.fields.routing,
			}
			if got := dp.newHandler(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultManagerNewProxy(t *testing.T) {
	type fields struct {
		proxies []*httputil.ReverseProxy
		routing *router.Router
	}
	type args struct {
		targetUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *httputil.ReverseProxy
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := defaultManager{
				proxies: tt.fields.proxies,
				routing: tt.fields.routing,
			}
			got, err := dp.newProxy(tt.args.targetUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("newProxy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newProxy() got = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
