package track_episode

import "github.com/google/uuid"

type TrackEpisode struct {
	ID              uuid.UUID
	EpisodeID       uuid.UUID
	TrackID         uuid.UUID
	CurrentPosition int

	PreviousPosition    int
	HighestPosition     int
	TimesAtPeakPosition int
	EpisodesCount       int
}
