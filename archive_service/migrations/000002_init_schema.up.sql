CREATE TABLE episodes (
    id UUID PRIMARY KEY,
    chart_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uniq_chart_day
ON episodes (chart_id, DATE(created_at AT TIME ZONE 'UTC'));

CREATE INDEX idx_episodes_chart
ON episodes (chart_id);

CREATE TABLE tracks (
    id UUID PRIMARY KEY,
    artist TEXT NOT NULL,
    title TEXT NOT NULL,
    normalized_key TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_tracks_title_artist
ON tracks (artist, title);

CREATE TABLE track_episode (
    id UUID PRIMARY KEY,
    episode_id UUID NOT NULL,
    track_id UUID NOT NULL,

    current_position INT NOT NULL,

    previous_position INT NOT NULL,
    highest_position INT NOT NULL,
    episodes_count INT NOT NULL,
    times_at_peak_position INT NOT NULL,
    
    CONSTRAINT fk_episode
        FOREIGN KEY (episode_id)
        REFERENCES episodes(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_track
        FOREIGN KEY (track_id)
        REFERENCES tracks(id)
        ON DELETE CASCADE

);

CREATE INDEX idx_tracks_in_episode
ON track_episode (episode_id);
