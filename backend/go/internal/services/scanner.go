package services

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"indexarr/internal/config"
	"indexarr/internal/models"
	"indexarr/internal/repository"
)

// Scanner handles media library scanning
type Scanner struct {
	db          *sql.DB
	config      *config.Config
	extractor   *Extractor
	tmdb        *TMDBClient
	tv          *TVClient
	broadcaster *Broadcaster
	running     bool
	stopChan    chan struct{}
	mu          sync.Mutex
}

// NewScanner creates a new scanner service
func NewScanner(db *sql.DB, cfg *config.Config, broadcaster *Broadcaster) *Scanner {
	return &Scanner{
		db:          db,
		config:      cfg,
		extractor:   NewExtractor(cfg.MediainfoPath, cfg.ScanTimeout),
		tmdb:        NewTMDBClient(cfg.TMDBAPIKey),
		tv:          NewTVClient(cfg.TMDBAPIKey), // Uses TMDB for TV shows
		broadcaster: broadcaster,
		stopChan:    make(chan struct{}),
	}
}

// IsRunning returns whether a scan is currently in progress
func (s *Scanner) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// Stop signals the scanner to stop
func (s *Scanner) Stop() {
	s.mu.Lock()
	if s.running {
		close(s.stopChan)
	}
	s.mu.Unlock()
}

// Scan performs a full library scan
func (s *Scanner) Scan() (*models.ScanResult, error) {
	return s.ScanPaths(s.config.MediaLibraryPaths)
}

// ScanPaths performs a scan on specified paths (used for manual scans via API)
func (s *Scanner) ScanPaths(paths []string) (*models.ScanResult, error) {
	log.Println("Starting scan")
	start := time.Now()

	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil, fmt.Errorf("scan already in progress")
	}
	s.running = true
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()

	result := &models.ScanResult{
		Errors: []string{},
	}

	// Update scan status to running
	status := &models.ScanStatus{
		Status:     "running",
		StartedAt:  time.Now().Format(time.RFC3339),
		FilesFound: 0,
	}
	if err := repository.UpdateScanStatus(s.db, status); err != nil {
		log.Printf("Failed to update scan status: %v", err)
	}

	// Collect all media files
	var files []string
	for _, libPath := range paths {
		if libPath == "" {
			continue
		}

		log.Printf("Scanning library path: %s", libPath)

		if _, err := os.Stat(libPath); os.IsNotExist(err) {
			result.Errors = append(result.Errors, fmt.Sprintf("Path does not exist: %s", libPath))
			continue
		}

		err := filepath.WalkDir(libPath, func(path string, d fs.DirEntry, err error) error {
			// Check for stop signal
			select {
			case <-s.stopChan:
				return fmt.Errorf("scan stopped by user")
			default:
			}

			if err != nil {
				log.Printf("Error accessing path %s: %v", path, err)
				return nil // Continue walking
			}

			if d.IsDir() {
				name := d.Name()

				// Skip hidden directories
				if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "@") {
					return fs.SkipDir
				}

				// Skip extra media folders
				for _, extraFolder := range s.config.SkipFolders {
					if strings.EqualFold(name, extraFolder) {
						return fs.SkipDir
					}
				}

				return nil
			}

			// Check if it's a video file
			if IsVideoFile(path) {
				files = append(files, path)
			}

			return nil
		})

		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error walking %s: %v", libPath, err))
		}
	}

	result.FilesFound = len(files)
	log.Printf("Found %d media files", result.FilesFound)

	// Update status with files found
	status.FilesFound = result.FilesFound
	repository.UpdateScanStatus(s.db, status)

	// Broadcast scan start to WebSocket clients
	if s.broadcaster != nil {
		s.broadcaster.BroadcastScanStart(result.FilesFound, status.StartedAt)
	}

	// Process each file sequentially
	for i, filePath := range files {
		// Check for stop signal
		select {
		case <-s.stopChan:
			status.Status = "stopped"
			status.CompletedAt = time.Now().Format(time.RFC3339)
			status.ErrorMessage = "Scan stopped by user"
			repository.UpdateScanStatus(s.db, status)
			// Broadcast stopped event to WebSocket clients
			if s.broadcaster != nil {
				s.broadcaster.BroadcastScanStopped()
			}
			return result, fmt.Errorf("scan stopped by user")
		default:
		}

		if err := s.processFile(filePath, result); err != nil {
			log.Printf("Error processing %s: %v", filePath, err)
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", filepath.Base(filePath), err))
		}

		result.FilesProcessed++

		// Update progress periodically
		if i%10 == 0 || i == len(files)-1 {
			status.FilesProcessed = result.FilesProcessed
			repository.UpdateScanStatus(s.db, status)
			// Broadcast progress to WebSocket clients
			if s.broadcaster != nil {
				s.broadcaster.BroadcastScanProgress(result.FilesProcessed, result.FilesFound)
			}
		}
	}

	// Update status to completed
	status.Status = "completed"
	status.CompletedAt = time.Now().Format(time.RFC3339)
	status.FilesProcessed = result.FilesProcessed
	if len(result.Errors) > 0 {
		status.ErrorMessage = fmt.Sprintf("%d errors during scan", len(result.Errors))
	}
	repository.UpdateScanStatus(s.db, status)

	// Broadcast completion to WebSocket clients
	if s.broadcaster != nil {
		s.broadcaster.BroadcastScanComplete(result.FilesProcessed, result.MoviesAdded, result.EpisodesAdded)
	}

	duration := time.Since(start)
	log.Printf("Scan completed in %v - %d files processed, %d movies added, %d episodes added, %d errors",
		duration.Round(time.Second), result.FilesProcessed, result.MoviesAdded, result.EpisodesAdded, len(result.Errors))

	if len(result.Errors) > 0 {
		// Log the first 100 errors for visibility
		for i, err := range result.Errors {
			if i >= 100 {
				log.Printf("  - ... %d more lines", len(result.Errors)-100)
				break
			}
			log.Printf("  - %s", err)
		}
	}

	return result, nil
}

// ScanMovie scans a single movie (used for manual refresh via API)
func (s *Scanner) ScanMovie(movieID int64) (*models.ScanResult, error) {
	movie, err := repository.GetMovieByID(s.db, movieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	result, err := s.ScanPaths([]string{movie.FilePath})
	if err != nil {
		return nil, err
	}

	// Remove movie if it was deleted from disk
	if result.FilesProcessed == 0 {
		log.Printf("Movie file not found during refresh, deleting movie: %s", movie.FilePath)
		if err := repository.DeleteMovie(s.db, movieID); err != nil {
			log.Printf("Failed to delete movie: %v", err)
		}
	} else {
		// Extract media info again to update any changes (e.g. new audio tracks)
		mediaInfo, fileSize, duration, err := s.extractor.Extract(movie.FilePath)
		if err != nil {
			log.Printf("Mediainfo extraction failed during refresh for %s: %v", movie.Title, err)
		} else {
			movie.MediaInfo = mediaInfo
			movie.FileSize = fileSize
			movie.Duration = duration / 60 // Convert seconds to minutes
		}

		// Fetch metadata from TMDB again to update any changes
		if err := s.tmdb.EnrichMovie(movie); err != nil {
			log.Printf("TMDB enrichment failed during refresh for %s: %v", movie.Title, err)
		} else {
			// Update movie with new metadata
			if err := repository.UpdateMovie(s.db, movie); err != nil {
				log.Printf("Failed to update movie during refresh: %v", err)
			}
			log.Printf("Movie refreshed: %s (%d)", movie.Title, movie.Year)
			result.MoviesAdded = 1 // Count as "added" for refresh purposes
			result.FilesProcessed = 1
			result.FilesFound = 1
			result.Errors = []string{}
			result.EpisodesAdded = 0
			result.MoviesAdded = 1
		}
	}

	return result, nil
}

// ScanSeries scans a single series (used for manual refresh via API)
func (s *Scanner) ScanSeries(seriesID int64) (*models.ScanResult, error) {
	log.Printf("Starting series refresh for ID: %d", seriesID)
	start := time.Now()

	result := &models.ScanResult{
		Errors: []string{},
	}

	// Step 1: Fetch series from database
	series, err := repository.GetSeriesByID(s.db, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch series: %w", err)
	}
	if series == nil {
		return nil, fmt.Errorf("series not found with ID: %d", seriesID)
	}

	log.Printf("Found series: %s", series.Title)

	// Step 2: Fetch all episodes to determine folders to scan
	episodes, err := repository.GetAllEpisodesForSeries(s.db, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch episodes: %w", err)
	}

	// Extract unique folder paths from existing episodes
	folderPaths := s.findSeriesFolderPaths(episodes)
	log.Printf("Will scan %d folder(s) for series: %v", len(folderPaths), folderPaths)

	// Step 3: Scan folder paths to detect new episodes
	scanResult, err := s.ScanPaths(folderPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to scan series folders: %w", err)
	}

	result.FilesFound = scanResult.FilesFound

	// Step 4 & 5: Check for missing episodes and delete them from database
	episodesToDelete := []int64{}
	for _, episode := range episodes {
		// Check if file still exists
		if _, err := os.Stat(episode.FilePath); os.IsNotExist(err) {
			log.Printf("Episode file missing: %s (S%02dE%02d), marking for removal", episode.FilePath, episode.SeasonNum, episode.EpisodeNum)
			episodesToDelete = append(episodesToDelete, episode.ID)
		}
	}

	// Delete missing episodes
	for _, episodeID := range episodesToDelete {
		if err := repository.DeleteEpisode(s.db, episodeID); err != nil {
			errMsg := fmt.Sprintf("Failed to remove missing episode %d: %v", episodeID, err)
			log.Printf("%s", errMsg)
			result.Errors = append(result.Errors, errMsg)
		}
	}

	// Delete seasons that have no episodes left
	if err := repository.DeleteEmptySeasons(s.db, seriesID); err != nil {
		log.Printf("Failed to delete empty seasons: %v", err)
	}

	// Step 6: Check if series folder is completely missing
	if scanResult.FilesFound == 0 && len(episodesToDelete) == len(episodes) {
		log.Printf("All episodes missing from disk, deleting series: %s", series.Title)
		if err := repository.DeleteSeries(s.db, seriesID); err != nil {
			errMsg := fmt.Sprintf("Failed to delete series: %v", err)
			log.Printf("%s", errMsg)
			result.Errors = append(result.Errors, errMsg)
		}
		return result, nil
	}

	// Step 7: Extract media info for each episode again to catch any changes
	for _, episode := range episodes {
		mediaInfo, fileSize, duration, err := s.extractor.Extract(episode.FilePath)
		if err != nil {
			log.Printf("Mediainfo extraction failed for %s: %v", episode.FilePath, err)
			// Continue with minimal info
			mediaInfo = &models.MediaInfo{
				VideoTracks:    []models.VideoTrack{},
				AudioTracks:    []models.AudioTrack{},
				SubtitleTracks: []models.SubtitleTrack{},
			}
		}

		episode.MediaInfo = mediaInfo
		episode.FileSize = fileSize
		episode.Duration = duration

		// Update episode in database
		if err := repository.UpdateEpisode(s.db, &episode); err != nil {
			log.Printf("Failed to update episode during refresh: %v", err)
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to update episode S%02dE%02d: %v", episode.SeasonNum, episode.EpisodeNum, err))
		}
	}

	// Re-enrich series metadata from TVDB
	if err := s.tv.EnrichSeries(series); err != nil {
		log.Printf("TVDB enrichment failed during refresh for %s: %v", series.Title, err)
	}

	// Update series in database
	if err := repository.UpdateSeries(s.db, series); err != nil {
		log.Printf("Failed to update series during refresh: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to update series: %v", err))
	}

	// Recalculate series counts
	if err := repository.UpdateSeriesCounts(s.db, seriesID); err != nil {
		log.Printf("Failed to update series counts: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to update series counts: %v", err))
	}

	result.FilesProcessed = scanResult.FilesProcessed
	result.EpisodesAdded = scanResult.EpisodesAdded
	result.MoviesAdded = 0 // No movies in series refresh

	// Merge any errors from the scan
	result.Errors = append(result.Errors, scanResult.Errors...)

	duration := time.Since(start)
	log.Printf("Series refresh completed in %v - %d files processed, %d episodes added, %d episodes deleted, %d errors",
		duration.Round(time.Second), result.FilesProcessed, result.EpisodesAdded, len(episodesToDelete), len(result.Errors))

	return result, nil
}

// findSeriesFolderPaths extracts unique parent directory paths from a list of episodes
// This supports series that may be split across multiple folders
func (s *Scanner) findSeriesFolderPaths(episodes []models.Episode) []string {
	folderMap := make(map[string]bool)

	for _, episode := range episodes {
		dir := filepath.Dir(episode.FilePath)
		folderMap[dir] = true
	}

	folders := make([]string, 0, len(folderMap))
	for folder := range folderMap {
		folders = append(folders, folder)
	}

	// Find common parent folder which does not belong to any media library path to avoid scanning the entire library if series episodes are stored in different folders. This is a common edge case for users who have their TV shows organized in multiple folders (e.g. by genre, by quality, etc.) but still want to be able to refresh the entire series metadata with one click.
	commonParent := findCommonParentFolder(folders)
	if commonParent != "" && !s.isLibraryPath(commonParent) {
		log.Printf("Series episodes are in multiple folders, using common parent folder for scan: %s", commonParent)
		return []string{commonParent}
	}

	return folders
}

// findCommonParentFolder takes a list of folder paths and returns the common parent folder if it exists, or an empty string if there is no common parent
func findCommonParentFolder(folders []string) string {
	if len(folders) == 0 {
		return ""
	}

	commonParent := folders[0]
	for _, folder := range folders[1:] {
		for !strings.HasPrefix(folder, commonParent) {
			commonParent = filepath.Dir(commonParent)
			if commonParent == "." || commonParent == "/" {
				return ""
			}
		}
	}

	return commonParent
}

// isLibraryPath checks if a given path is one of the configured media library paths
func (s *Scanner) isLibraryPath(path string) bool {
	for _, libPath := range s.config.MediaLibraryPaths {
		if strings.EqualFold(filepath.Clean(libPath), filepath.Clean(path)) {
			return true
		}
	}
	return false
}

// processFile handles a single media file
func (s *Scanner) processFile(filePath string, result *models.ScanResult) error {
	// Parse filename
	parsed := ParseFilename(filePath)

	// Extract media info
	mediaInfo, fileSize, duration, err := s.extractor.Extract(filePath)
	if err != nil {
		log.Printf("Mediainfo extraction failed for %s: %v", filePath, err)
		// Continue with minimal info
		mediaInfo = &models.MediaInfo{
			VideoTracks:    []models.VideoTrack{},
			AudioTracks:    []models.AudioTrack{},
			SubtitleTracks: []models.SubtitleTrack{},
		}
	}

	if parsed.IsSeries {
		return s.processEpisode(filePath, parsed, mediaInfo, fileSize, duration, result)
	}
	return s.processMovie(filePath, parsed, mediaInfo, fileSize, duration, result)
}

// processMovie handles a movie file
func (s *Scanner) processMovie(filePath string, parsed *ParsedFilename, mediaInfo *models.MediaInfo, fileSize int64, duration int, result *models.ScanResult) error {
	// Check if movie already exists by file path
	exists, err := repository.MovieExistsByFilePath(s.db, filePath)
	if err != nil {
		return fmt.Errorf("failed to check for existing movie: %w", err)
	}
	if exists {
		log.Printf("Movie already exists for file: %s", filePath)
		return nil
	}

	movie := &models.Movie{
		Title:     parsed.Title,
		Year:      parsed.Year,
		Duration:  duration / 60, // Convert seconds to minutes
		Status:    "available",
		FileSize:  fileSize,
		FilePath:  filePath,
		Container: GetContainer(filePath),
		DateAdded: time.Now().Format(time.RFC3339),
		MediaInfo: mediaInfo,
	}

	// Try to enrich with TMDB metadata
	if err := s.tmdb.EnrichMovie(movie); err != nil {
		log.Printf("TMDB enrichment failed for %s: %v", parsed.Title, err)
		// Continue without TMDB data
	}

	// Insert into database
	_, err = repository.InsertMovie(s.db, movie)
	if err != nil {
		return fmt.Errorf("failed to insert movie: %w", err)
	}

	result.MoviesAdded++
	log.Printf("Added movie: %s (%d)", movie.Title, movie.Year)
	return nil
}

// processEpisode handles a TV episode file
func (s *Scanner) processEpisode(filePath string, parsed *ParsedFilename, mediaInfo *models.MediaInfo, fileSize int64, duration int, result *models.ScanResult) error {
	// Check if series exists, create if not
	series, err := repository.GetSeriesByTitle(s.db, parsed.Title)
	if err != nil {
		return fmt.Errorf("failed to lookup series: %w", err)
	}

	var seriesID int64
	var seriesTMDBID int

	if series == nil {
		// Create new series
		newSeries := &models.Series{
			Title:     parsed.Title,
			Status:    "ongoing",
			DateAdded: time.Now().Format(time.RFC3339),
		}

		// Try to enrich with TMDB metadata
		if err := s.tv.EnrichSeries(newSeries); err != nil {
			log.Printf("TMDB TV enrichment failed for %s: %v", parsed.Title, err)
		}

		// Check if series with same TVDB ID already exists (prevents duplicates)
		if newSeries.TVDBId > 0 {
			existingSeries, err := repository.GetSeriesByTVDBId(s.db, newSeries.TVDBId)
			if err != nil {
				return fmt.Errorf("failed to lookup series by TVDB ID: %w", err)
			}
			if existingSeries != nil {
				// Series already exists, reuse it
				seriesID = existingSeries.ID
				seriesTMDBID = int(existingSeries.TVDBId)
				log.Printf("Found existing series: %s (TVDB ID: %d)", existingSeries.Title, newSeries.TVDBId)
				// Skip the InsertSeries step below
			} else {
				// New series, insert it
				seriesID, err = repository.InsertSeries(s.db, newSeries)
				if err != nil {
					return fmt.Errorf("failed to insert series: %w", err)
				}
				seriesTMDBID = int(newSeries.TVDBId)
				log.Printf("Added series: %s (TVDB ID: %d)", newSeries.Title, newSeries.TVDBId)
			}
		} else {
			// No TVDB ID, insert new series anyway
			seriesID, err = repository.InsertSeries(s.db, newSeries)
			if err != nil {
				return fmt.Errorf("failed to insert series: %w", err)
			}
			seriesTMDBID = int(newSeries.TVDBId)
			log.Printf("Added series: %s (no TVDB ID)", newSeries.Title)
		}
	} else {
		seriesID = series.ID
		seriesTMDBID = int(series.TVDBId)
	}

	// Create episode
	episode := &models.Episode{
		SeriesID:   seriesID,
		SeasonNum:  parsed.Season,
		EpisodeNum: parsed.Episode,
		Duration:   duration, // Already in seconds
		Status:     "available",
		FileSize:   fileSize,
		FilePath:   filePath,
		DateAdded:  time.Now().Format(time.RFC3339),
		MediaInfo:  mediaInfo,
	}

	// Try to enrich episode with TMDB metadata
	if err := s.tv.EnrichEpisode(episode, seriesTMDBID); err != nil {
		log.Printf("TMDB episode enrichment failed: %v", err)
	}

	// If no title from TMDB, create a default title
	if episode.Title == "" {
		episode.Title = fmt.Sprintf("Episode %d", parsed.Episode)
	}

	// Ensure season exists
	_, err = repository.GetOrCreateSeason(s.db, seriesID, parsed.Season)
	if err != nil {
		log.Printf("Failed to create season: %v", err)
	}

	// Check if episode already exists
	existingEpisode, err := repository.GetEpisodeBySeriesSeasonEpisode(s.db, seriesID, parsed.Season, parsed.Episode)
	if err != nil {
		return fmt.Errorf("failed to lookup episode: %w", err)
	}

	if existingEpisode != nil {
		// Episode already exists - update if file path or details changed
		existingEpisode.Title = episode.Title
		existingEpisode.Duration = episode.Duration
		existingEpisode.Status = "available"
		existingEpisode.FileSize = episode.FileSize
		existingEpisode.FilePath = filePath

		if err := repository.UpdateEpisode(s.db, existingEpisode); err != nil {
			return fmt.Errorf("failed to update episode: %w", err)
		}
		log.Printf("Updated episode: %s S%02dE%02d - %s", parsed.Title, parsed.Season, parsed.Episode, existingEpisode.Title)
	} else {
		// New episode - insert it
		_, err = repository.InsertEpisode(s.db, episode)
		if err != nil {
			return fmt.Errorf("failed to insert episode: %w", err)
		}
		log.Printf("Added episode: %s S%02dE%02d - %s", parsed.Title, parsed.Season, parsed.Episode, episode.Title)
		result.EpisodesAdded++
	}

	// Update series counts
	repository.UpdateSeriesCounts(s.db, seriesID)
	return nil
}

// GetStatus returns the current scan status
func (s *Scanner) GetStatus() (*models.ScanStatus, error) {
	return repository.GetScanStatus(s.db)
}
