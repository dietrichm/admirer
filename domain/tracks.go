//go:generate mockgen -source tracks.go -destination tracks_mock.go -package domain

package domain

// Track represents a track on an external service.
type Track struct {
	Artist string
	Name   string
}
