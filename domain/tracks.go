//go:generate mockgen -source tracks.go -destination tracks_mock.go -package domain

package domain

import "fmt"

// Track represents a track on an external service.
type Track struct {
	Artist string
	Name   string
}

func (t Track) String() string {
	return fmt.Sprintf("%s - %s", t.Artist, t.Name)
}
