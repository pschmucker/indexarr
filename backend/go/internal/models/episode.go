package models

type Episode struct {
	ID           int64      `json:"id"`
	SeriesID     int64      `json:"seriesId"`
	SeasonNum    int        `json:"seasonNum"`
	EpisodeNum   int        `json:"episodeNum"`
	Title        string     `json:"title"`
	Duration     int        `json:"duration"` // seconds
	Status       string     `json:"status"` // available, missing
	FileSize     int64      `json:"fileSize"` // bytes
	FilePath     string     `json:"filePath"`
	DateAdded    string     `json:"dateAdded"`
	MediaInfo    *MediaInfo `json:"mediaInfo"`
}
