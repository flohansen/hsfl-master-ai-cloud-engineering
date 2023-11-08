package handler

import(
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// ReverseProxyConfig represents the configuration for the Reverse Proxy.
type ReverseProxyConfig struct {
	AuthServiceURL       string
	BulletinBoardServiceURL string
	FeedServiceURL         string
}

// ReverseProxyHandler is a struct representing the Reverse Proxy.
type ReverseProxyHandler struct {
	AuthService       *httputil.ReverseProxy
	BulletinBoardService *httputil.ReverseProxy
	FeedService         *httputil.ReverseProxy
}

func NewReverseProxy(config ReverseProxyConfig) *ReverseProxyHandler {
	return &ReverseProxyHandler{
		AuthService:       newSingleHostReverseProxy(config.AuthServiceURL),
		BulletinBoardService: newSingleHostReverseProxy(config.BulletinBoardServiceURL),
		FeedService:          newSingleHostReverseProxy(config.FeedServiceURL),
	}
}
func newSingleHostReverseProxy(targetURL string) *httputil.ReverseProxy {
	target, _ := url.Parse(targetURL)
	return httputil.NewSingleHostReverseProxy(target)
}

func CleanPath(path string) string {
	
	// Split the path by "/"
    segments := strings.Split(path, "/")

    // Create a slice to store the non-empty segments
    cleanSegments := make([]string, 0, len(segments))

    // Iterate through the segments and add non-empty segments to the cleanSegments slice
    for _, segment := range segments {
        if segment != "" {
            cleanSegments = append(cleanSegments, segment)
        }
    }

    // If all segments were empty or the path was empty, return an empty string
    if len(cleanSegments) == 0 {
        return ""
    }

    // Join the clean segments with "/" to create the cleaned path
    cleanPath := strings.Join(cleanSegments, "/")

    return cleanPath
}
func (rp *ReverseProxyHandler) determineTarget(path string) string {
    if(path == ""){
		return "unknown"
	}
	cleanPath := CleanPath(path)
	segments := strings.Split(cleanPath, "/")

    // Determine the target based on the first segment
    if len(segments) > 0 {
        return segments[0]
    }

    // Default to a target of "unknown" if the path is empty
    return "unknown"
}

func (rp *ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Determine the target service based on the request path
	target := rp.determineTarget(r.URL.Path)

	// Route the request to the appropriate service
	switch target {
	case "auth":
		rp.AuthService.ServeHTTP(w, r)
	case "bulletinboard":
		rp.BulletinBoardService.ServeHTTP(w, r)
	case "feed":
		rp.FeedService.ServeHTTP(w, r)
	default:
		http.Error(w, "Service not found", http.StatusNotFound)
	}
}