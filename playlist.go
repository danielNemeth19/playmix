package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"math"
	"math/rand"
	"path/filepath"

	"github.com/alfg/mp4"
)

// TODO: let's sanitize track title by cutting the vlc record prefix
type MediaItem struct {
	AbsPath  string
	Location string
	Dir      string
	Name     string
	Id       int
	Duration float64
}

// TODO: This dirName could be used writing a proper title
func (m *MediaItem) getRelativeDir(rootParts []string) {
	fileParts := getPathParts(m.AbsPath)

	if len(fileParts) == len(rootParts) {
		m.Dir = filepath.Base(filepath.Dir(m.AbsPath))
	} else {
		relativeParts := fileParts[len(rootParts)-1:]
		m.Dir = filepath.Join(relativeParts...)
	}
}

type Summarizer struct {
	dBucket       DurationBucket
	ratio         int
	totalDuration float64
	totalScanned  int
	totalSelected int
}

func (s Summarizer) getRealRatio() float64 {
	return float64(s.totalSelected) / float64(s.totalScanned) * 100
}

func (s Summarizer) getData(w io.Writer) {
	fmt.Fprintf(w, "Total scanned: %d\n", s.totalScanned)
	fmt.Fprintf(w, "Duration distribution:\n")
	s.dBucket.summarize(w)
	fmt.Fprintf(w, "Total duration is: %f sec -- (%f) minutes\n", s.totalDuration, s.totalDuration/60)
	fmt.Fprintf(w, "Total selected: %d -- required ratio: %d -- got: %.2f%%\n", s.totalSelected, s.ratio, s.getRealRatio())
}

type DurationBucket struct {
	Dur0_5     int
	Dur5_10    int
	Dur10_30   int
	Dur30_60   int
	Dur60_180  int
	Dur180_240 int
	DurOver240 int
}

func (d *DurationBucket) allocate(duration float64) {
	switch {
	case duration < 5:
		d.Dur0_5++
	case duration <= 10:
		d.Dur5_10++
	case duration <= 30:
		d.Dur10_30++
	case duration <= 60:
		d.Dur30_60++
	case duration <= 180:
		d.Dur60_180++
	case duration <= 240:
		d.Dur180_240++
	default:
		d.DurOver240++
	}
}

func (d *DurationBucket) summarize(w io.Writer) {
	fmt.Fprintf(w, "Bucket <5 seconds: %d\n", d.Dur0_5)
	fmt.Fprintf(w, "Bucket 5-10 seconds: %d\n", d.Dur5_10)
	fmt.Fprintf(w, "Bucket 10-30 seconds: %d\n", d.Dur10_30)
	fmt.Fprintf(w, "Bucket 30-60 seconds: %d\n", d.Dur30_60)
	fmt.Fprintf(w, "Bucket 60-180 seconds: %d\n", d.Dur60_180)
	fmt.Fprintf(w, "Bucket 180-240 seconds: %d\n", d.Dur180_240)
	fmt.Fprintf(w, "Bucket 240< seconds: %d\n", d.DurOver240)
}

func selector(ratio int) bool {
	n := rand.Intn(100)
	if n < ratio {
		return true
	}
	return false
}

func toSkip(name string, skip []string) bool {
	for _, f := range skip {
		if f == name {
			return true
		}
	}
	return false
}

func isIncluded(rootParts []string, path string, include []string) bool {
	if len(include) == 0 {
		return true
	}
	parts := getPathParts(path)
	for _, folder := range parts[len(rootParts):] {
		for _, f := range include {
			if folder == f {
				return true
			}
		}
	}
	return false
}

func dateFilter(d fs.DirEntry, params Params) bool {
	file, _ := d.Info()
	mTime := file.ModTime().UTC()
	if mTime.After(params.fdate) && mTime.Before(params.tdate) {
		return true
	}
	return false
}

func collectMediaContent(p string, fsys fs.FS, params Params) ([]MediaItem, Summarizer, error) {
	var items []MediaItem
	rootParts := getPathParts(p)
	summary := Summarizer{ // TODO: factors this out to be a parameter
		ratio:         params.ratio,
		totalScanned:  0,
		totalSelected: 0,
		totalDuration: 0,
		dBucket:       DurationBucket{},
	}
	idx := 0
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && toSkip(d.Name(), params.skipF) {
			return filepath.SkipDir
		}
		absPath := filepath.Join(p, path)
		if !d.IsDir() && isMediaFile(filepath.Ext(d.Name())) && isIncluded(rootParts, absPath, params.includeF) && dateFilter(d, params) {
			if selector(params.ratio) {
				duration, err := getDuration(fsys, path)
				if err != nil {
					return err
				}
				summary.dBucket.allocate(duration)
				if duration > float64(params.minDuration) && duration < float64(params.maxDuration) {
					location := getUrlEncodedPath(absPath)
					item := MediaItem{Id: idx, AbsPath: absPath, Location: location, Name: d.Name(), Duration: duration}
					item.getRelativeDir(rootParts)
					items = append(items, item)
					summary.totalDuration += duration
					summary.totalSelected++
				}
			}
			summary.totalScanned++
			idx++
			if summary.totalScanned%500 == 0 {
				fmt.Printf("Processed %d files\n", summary.totalScanned)
			}
		}
		return nil
	})
	return items, summary, err
}

type unbufferedReaderAt struct {
	R io.Reader
	S *io.SectionReader
}

func NewUnbufferedReaderAt(r io.Reader, size int64) io.ReaderAt {
	return &unbufferedReaderAt{
		R: r,
		S: io.NewSectionReader(r.(io.ReaderAt), 0, size),
	}
}

func (u *unbufferedReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	return u.S.ReadAt(p, off)
}

func getDuration(fsys fs.FS, p string) (duration float64, err error) {
	file, err := fsys.Open(p)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	readerAt := NewUnbufferedReaderAt(file, info.Size())
	mp4, err := mp4.OpenFromReader(readerAt, info.Size())
	if err != nil {
		return 0, err
	}
	if mp4.Moov == nil {
		return 0, fmt.Errorf("Moov box not found for %s. Is this mp4?", p)
	}
	rawDuration := float64(mp4.Moov.Mvhd.Duration)
	timeScale := float64(mp4.Moov.Mvhd.Timescale)

	duration = rawDuration / timeScale
	return duration, nil
}

func randomizePlaylist(playlist []MediaItem, stabilizer int) {
	if stabilizer > len(playlist) {
		rand.Shuffle(len(playlist), func(i, j int) {
			playlist[i], playlist[j] = playlist[j], playlist[i]
		})
	} else {
		rand.Shuffle(len(playlist), func(i, j int) {
			if i%stabilizer != 0 {
				playlist[i], playlist[j] = playlist[j], playlist[i]
			}
		})
	}
}

func writePlayList(s any, w io.Writer) error {
	xmlH := []byte(xml.Header)
	_, err := w.Write(xmlH)
	if err != nil {
		return fmt.Errorf("Error writing header: %w\n", err)
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "\t")
	err = encoder.Encode(&s)
	if err != nil {
		return fmt.Errorf("Error in encoding xml: %w\n", err)
	}
	return nil
}

func buildPlayList(content []MediaItem, options Options) *PlayList {
	playList := &PlayList{Xmlns: Xmlns, XmlnsVlc: XmlnsVlc, Version: "1"}
	trackList := &TrackList{}
	tracks := []*Track{}

	for i, media := range content {
		ext := &Extension{Application: ExtensionApplication, Id: i}
		if options.audio != "" {
			ext.Option = options.audio
		} else {
			if options.StartTime > 0 {
				// ext.Option = options.start_time
				fmt.Printf("deal with this later")
			}
		}
		track := &Track{Location: media.Location, Title: media.Name, Duration: math.Round(media.Duration), Ext: *ext}
		tracks = append(tracks, track)
	}
	trackList.Tracks = tracks
	playList.Tl = *trackList
	return playList
}
