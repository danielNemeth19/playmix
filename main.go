package main

import (
	// "bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/abema/go-mp4"
)

const (
	ExtensionApplication = "http://www.videolan.org/vlc/playlist/0"
	Xmlns                = "http://xspf.org/ns/0/"
	XmlnsVlc             = "http://www.videolan.org/vlc/playlist/ns/0/"
)

var videoExtensions = []string{".mp4", ".mkv", ".avi", ".flv", ".mpeg"}

type Extension struct {
	XMLName     xml.Name `xml:"extension"`
	Application string   `xml:"application,attr"`
	Id          int      `xml:"vlc:id"`
}

type Track struct {
	XMLName  xml.Name  `xml:"track"`
	Location string    `xml:"location"`
	Ext      Extension `xml:"extension"`
}

type TrackList struct {
	XMLName xml.Name `xml:"trackList"`
	Tracks  []*Track `xml:"track"`
}

type PlayList struct {
	XMLName  xml.Name  `xml:"playlist"`
	Xmlns    string    `xml:"xmlns,attr"`
	XmlnsVlc string    `xml:"xmlns:vlc,attr"`
	Version  string    `xml:"version,attr"`
	Title    string    `xml:"title"`
	Tl       TrackList `xml:"trackList"`
}

type FolderContent struct {
	Name string
	Id   int
}

func isVideoFile(ext string) bool {
	for _, v := range videoExtensions {
		if v == ext {
			return true
		}
	}
	return false
}

func collectExtensions(p string, extensions *[]string, seen *map[string]bool) *[]string {
	fmt.Println("Checking folder: ", p)
	files, err := os.ReadDir(p)
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	for _, obj := range files {
		if obj.IsDir() {
			extensions = collectExtensions(p+"/"+obj.Name(), extensions, seen)
		} else {
			extension := filepath.Ext(obj.Name())
			if !(*seen)[extension] {
				fmt.Printf("Seen first: %s -- %s, %v\n", extension, obj.Name(), obj.IsDir())
				(*seen)[extension] = true
				*extensions = append(*extensions, extension)
			}
		}
	}
	return extensions
}

func getFolderContent(p string) ([]FolderContent, error) {
	extensions := []string{}
	seen := map[string]bool{}
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("Error raised: %w\n", err)
	}
	var content []FolderContent
	for i, file := range files {
		extension := filepath.Ext(file.Name())
		fmt.Printf("File %d: %s, extension: %s\n", i, file.Name(), extension)
		fmt.Printf("Verdict: %v\n", isVideoFile(extension))
		if !seen[extension] {
			seen[extension] = true
			extensions = append(extensions, extension)
		}
	}
	fmt.Printf("Extensions: %v\n", extensions)
	return content, nil
}

func buildTrackList(p string) *TrackList {
	files, err := os.ReadDir(p)
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	trackList := &TrackList{}
	var tracks []*Track
	for i, v := range files {
		ext := &Extension{Application: ExtensionApplication, Id: i}
		location := "file://" + p + v.Name()
		track := &Track{Location: location, Ext: *ext}
		tracks = append(tracks, track)
	}
	trackList.Tracks = tracks
	return trackList
}

func getPath() (string, error) {
	path := os.Getenv("MEDIA_SOURCE")
	if path == "" {
		return "", fmt.Errorf("MEDIA_SOURCE environment variable not set")
	}
	return path, nil
}

func dumpConsole(s any) {
	out, _ := xml.MarshalIndent(s, " ", "  ")
	fmt.Println(xml.Header + string(out))
}

func writePlayList(s any) error {
	outFile, err := os.Create("temp.xspf")
	if err != nil {
		return fmt.Errorf("Error creating file: %w\n", err)
	}
	defer outFile.Close()

	_, err = outFile.WriteString(xml.Header)
	if err != nil {
		return fmt.Errorf("Error writing header: %w\n", err)
	}

	encoder := xml.NewEncoder(outFile)
	encoder.Indent("", "\t")
	err = encoder.Encode(&s)
	if err != nil {
		return fmt.Errorf("Error in encoding xml: %w\n", err)
	}
	return nil
}

func testMetaData(p string) {
	f, err := os.Open(p)
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	defer f.Close()
	_, err = mp4.ReadBoxStructure(f, func(h *mp4.ReadHandle) (interface{}, error) {
		if h.BoxInfo.Type.String() == "mvhd" {
			fmt.Println("depth", len(h.Path))

			// Box Type (e.g. "mdhd", "tfdt", "mdat")
			fmt.Println("type", h.BoxInfo.Type.String())

			// Box Size
			fmt.Println("size", h.BoxInfo.Size)
            box, _, err := h.ReadPayload()
            if err != nil {
                return nil, err
            }
			str, err := mp4.Stringify(box, h.BoxInfo.Context)
            if err != nil {
                fmt.Printf("Error raised: %s\n", err)
                return nil, err
            }
            fmt.Println(str)
		}

		// if h.BoxInfo.IsSupportedType() {
			// Payload
			// box, _, err := h.ReadPayload()
			// if err != nil {
				// return nil, err
			// }
			// str, err := mp4.Stringify(box, h.BoxInfo.Context)
			// if err != nil {
				// return nil, err
			// }
			// f, _ := os.Create("temp.txt")
			// buffer := bufio.NewWriter(f)
			// buffer.WriteString(str)

			// Expands children
			// return h.Expand()
		// }
		return nil, nil
	})
}

func main() {
	path, err := getPath()
	if err != nil {
		log.Fatalf("Error raised: %s\n", err)
	}
	// extensions := &[]string{}
	// seen := &map[string]bool{}
	// res := collectExtensions(path, extensions, seen)
	// fmt.Printf("Extensions: %v\n", *res)

	testMetaData(path)
	// content, err := getFolderContent(path)
	// if err != nil {
	// log.Fatalf("Error raised: %s\n", err)
	// }
	// fmt.Printf("Content: %v\n", content)
	// trackList := buildTrackList(path)

	// pl := &PlayList{Xmlns: Xmlns, XmlnsVlc: XmlnsVlc, Version: "1", Title: "Test playlist", Tl: *trackList}
	// writePlayList(pl)
}
