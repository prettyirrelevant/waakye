package convert

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"

	"github.com/prettyirrelevant/kilishi/streaming_platforms/aggregator"
	"github.com/prettyirrelevant/shaki/cmd/services"
)

var (
	ErrInvalidPlaylistURL               = errors.New("link provided does not match any of the supported streaming platforms")
	ErrPlaylistURLRequired              = errors.New("please provide a link to the playlist")
	ErrPlaylistSourceAndDestinationSame = errors.New("source and destination cannot be the same")
)

var ConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a playlist",
	Run: func(cmd *cobra.Command, args []string) {
		url := getPlaylistURLInput()
		source := getStreamingPlatformInput("Source")
		destination := getStreamingPlatformInput("Destination")
		if source == destination {
			panic(ErrPlaylistSourceAndDestinationSame)
		}

		// get the playlist
		s := spinner.New(spinner.CharSets[11], 10*time.Millisecond)
		s.Suffix = fmt.Sprintf(" Fetching playlist %s on %s...\n", url, source)
		s.Color("green", "bold") //nolint:errcheck // error check is needless here
		s.Start()

		playlist, err := services.GetPlaylist(url, source)
		s.Stop()
		if err != nil {
			panic(fmt.Errorf("an error occurred while fetching the playlist: %s", err.Error()))
		}

		// now search for each tracks in the playlist
		var wg sync.WaitGroup
		var result sync.Map

		s.Suffix = fmt.Sprintf(" Searching for %d tracks on %s...\n", len(playlist.Data.Tracks), destination)
		s.Restart()
		for index, track := range playlist.Data.Tracks {
			wg.Add(1)
			go func(i int, entry services.TrackResponse) {
				defer wg.Done()
				resp, _err := services.FindTrack(entry.Title, destination, entry.Artists)
				if _err != nil {
					result.Store(i, TrackResult{Success: false})
				} else {
					result.Store(i, TrackResult{Success: true, Result: resp.Data})
				}
			}(index, track)
		}
		wg.Wait()
		s.Stop()

		// show the summary of the playlist search.
		var successfulSearches []services.TrackResponse
		var failedSearchesIndex []int
		for i := 0; i < len(playlist.Data.Tracks); i++ {
			value, _ := result.Load(i)
			result, _ := value.(TrackResult)

			if result.Success {
				successfulSearches = append(successfulSearches, result.Result)
			} else {
				failedSearchesIndex = append(failedSearchesIndex, i)
			}
		}

		fmt.Printf("Found %d tracks from a total of %d.\n", len(successfulSearches), len(playlist.Data.Tracks))
		if len(playlist.Data.Tracks)-len(successfulSearches) > 0 {
			fmt.Println("Details of tracks not found below:")
			for _, v := range failedSearchesIndex {
				fmt.Printf("Track %d -> Title: %s Artists: %+v\n", v, playlist.Data.Tracks[v].Title, playlist.Data.Tracks[v].Artists)
			}
		}

		// show a prompt asking if user wants to proceed and create the playlist
		getConfirmationInput("Do you want to continue")
		s.Suffix = fmt.Sprintf(" Creating playlist on %s with %d tracks...\n", destination, len(successfulSearches))
		s.Restart()
		createPlaylistResp, err := services.CreatePlaylist(playlist.Data.Title, playlist.Data.Description, destination, successfulSearches)
		s.Stop()
		if err != nil {
			panic(err)
		}

		err = clipboard.Init()
		if err != nil {
			panic(err)
		}

		clipboard.Write(clipboard.FmtText, []byte(createPlaylistResp.Data))
		fmt.Printf("Playlist created successfully and copied to clipboard! Enjoy!!!\nURL: %s\n", createPlaylistResp.Data)
	},
}

type TrackResult struct {
	Success bool
	Result  services.TrackResponse
}

func getPlaylistURLInput() string {
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return ErrPlaylistURLRequired
		}

		if ok, _ := regexp.MatchString(`^https:\/\/(open\.spotify\.com|deezer\.com|music\.youtube\.com)`, input); !ok {
			return ErrInvalidPlaylistURL
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Link to the playlist",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
	}

	return result
}

func getConfirmationInput(label string) string {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
	}

	return result
}

func getStreamingPlatformInput(label string) string {
	prompt := promptui.Select{
		Label:    label,
		Items:    allMusicStreamingPlatforms(),
		HideHelp: true,
	}

	_, result, err := prompt.Run()
	if err != nil {
		os.Exit(1)
	}

	return result
}

func allMusicStreamingPlatforms() []string {
	var platforms []string
	for i := range aggregator.AllMusicStreamingPlatforms {
		platforms = append(platforms, string(i))
	}

	return platforms
}
