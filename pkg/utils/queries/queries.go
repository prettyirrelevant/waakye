package queries

import (
	"github.com/prettyirrelevant/waakye/pkg/utils/types"
)

// TrackToSearchQuery takes a track and transforms it into a search query.
func TrackToSearchQuery(track types.Track) string {
	searchQuery := track.Title + " by"

	for index, artiste := range track.Artists {
		if len(track.Artists) == 1 {
			searchQuery += " " + artiste
		} else if len(track.Artists) > 1 && len(track.Artists)-1 == index {
			searchQuery += " and " + artiste
		} else {
			searchQuery += " " + artiste
			if index < len(track.Artists)-2 {
				searchQuery += ","
			}
		}
	}
	return searchQuery
}
