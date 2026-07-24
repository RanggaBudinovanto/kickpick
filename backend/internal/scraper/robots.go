package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/temoto/robotstxt"
)

// CheckRobotsAllowed fetches <scheme>://<host>/robots.txt for targetURL and
// verifies userAgent is allowed to crawl its path. Section 14 PRD: every
// adapter must respect robots.txt — this makes that an enforced pre-flight
// check rather than a one-time manual read during development (PENDING.md
// previously flagged robots.txt as checked once by hand, not automatically).
//
// Fails closed: any error fetching/parsing robots.txt is treated as "not
// allowed", since silently proceeding on a fetch failure could accidentally
// scrape a site that actually disallows it.
func CheckRobotsAllowed(targetURL, userAgent string) error {
	parsed, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsed.Scheme, parsed.Host)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(robotsURL)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", robotsURL, err)
	}
	defer resp.Body.Close()

	robots, err := robotstxt.FromResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", robotsURL, err)
	}

	if !robots.TestAgent(parsed.Path, userAgent) {
		return fmt.Errorf("robots.txt at %s disallows %q for user-agent %q", robotsURL, parsed.Path, userAgent)
	}

	return nil
}
