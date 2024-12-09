package goblet

import "fmt"

func UrlFor(endpoint string) string {
	return fmt.Sprintf("/%s", endpoint)
}
