package tinderbee

import (
	"github.com/mnzt/tinder"
	"github.com/muesli/beehive/bees"
)

// TriggerUpdateEvent triggers an update event
func (mod *TinderBee) TriggerUpdateEvent(update *tinder.UpdatesResponse) {
	// Using following from update as event

	// type UpdatesResponse struct {
	//     Matches []struct {
	//         ID                string `json:"_ID"`
	//         CommonFriendCount int    `json:"common_friend_count"`
	//         CommonLikeCount   int    `json:"common_like_count"`
	//         MessageCount      int    `json:"message_count"`
	//         Person struct {
	//             ID       string `json:"_ID"`
	//             Bio      string `json:"bio"`
	//             Birth    string `json:"birth_date"`
	//             Gender   int    `json:"gender"`
	//             Name     string `json:"name"`
	//             PingTime string `json:"ping_time"`
	//         }   `json:"person"`
	//     }   `json:"matches"`
	// }

	for _, v := range update.Matches {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "update",
			Options: []bees.Placeholder{
				{
					Name:  "ID",
					Type:  "string",
					Value: v.ID,
				},
				{
					Name:  "common_friend_count",
					Type:  "int",
					Value: v.CommonFriendCount,
				},
				{
					Name:  "common_like_count",
					Type:  "int",
					Value: v.CommonLikeCount,
				},
				{
					Name:  "message_count",
					Type:  "int",
					Value: v.MessageCount,
				},
				{
					Name:  "person_ID",
					Type:  "string",
					Value: v.Person.ID,
				},
				{
					Name:  "person_bio",
					Type:  "string",
					Value: v.Person.Bio,
				},
				{
					Name:  "person_birth",
					Type:  "string",
					Value: v.Person.Birth,
				},
				{
					Name:  "person_gender",
					Type:  "int",
					Value: v.Person.Gender,
				},
				{
					Name:  "person_name",
					Type:  "string",
					Value: v.Person.Name,
				},
				{
					Name:  "person_ping_time",
					Type:  "string",
					Value: v.Person.PingTime,
				},
			},
		}
		mod.evchan <- ev
	}
}

// TriggerUserEvents triggers all events for the given user
func (mod *TinderBee) TriggerUserEvents(user *tinder.UserResponse) {
	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "user",
		Options: []bees.Placeholder{
			{
				Name:  "status",
				Type:  "int",
				Value: user.Status,
			},
			{
				Name:  "connection_count",
				Type:  "int",
				Value: user.Results.ConnectionCount,
			},
			{
				Name:  "common_like_count",
				Type:  "int",
				Value: user.Results.CommonLikeCount,
			},
			{
				Name:  "common_friend_count",
				Type:  "int",
				Value: user.Results.CommonFriendCount,
			},
			{
				Name:  "common_likes",
				Type:  "[]string",
				Value: user.Results.CommonLikes,
			},
			{
				Name:  "common_interests",
				Type:  "[]string",
				Value: user.Results.CommonInterests,
			},
			{
				Name:  "uncommon_interests",
				Type:  "[]string",
				Value: user.Results.UncommonInterests,
			},
			{
				Name:  "common_friends",
				Type:  "[]string",
				Value: user.Results.CommonFriends,
			},
			{
				Name:  "ID",
				Type:  "string",
				Value: user.Results.ID,
			},
			{
				Name:  "bio",
				Type:  "string",
				Value: user.Results.Bio,
			},
			{
				Name:  "birthdate",
				Type:  "string",
				Value: user.Results.BirthDate.String(),
			},
			{
				Name:  "gender",
				Type:  "int",
				Value: user.Results.Gender,
			},
			{
				Name:  "name",
				Type:  "string",
				Value: user.Results.Name,
			},
			{
				Name:  "ping_time",
				Type:  "string",
				Value: user.Results.PingTime.String(),
			},
			{
				Name:  "birthdate_info",
				Type:  "string",
				Value: user.Results.BirthDateInfo,
			},
			{
				Name:  "distance_miles",
				Type:  "int",
				Value: user.Results.DistanceMi,
			},
		},
	}
	mod.evchan <- ev
}

// TriggerRecommendationsEvent triggers a recommendation event
func (mod *TinderBee) TriggerRecommendationsEvent(recs *tinder.RecommendationsResponse) {
	for _, v := range recs.Results {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "recommendation",
			Options: []bees.Placeholder{
				{
					Name:  "ID",
					Type:  "string",
					Value: v.ID,
				},
				{
					Name:  "bio",
					Type:  "string",
					Value: v.Bio,
				},
				{
					Name:  "birth",
					Type:  "string",
					Value: v.Birth.String(),
				},
				{
					Name:  "gender",
					Type:  "int",
					Value: v.Gender,
				},
				{
					Name:  "name",
					Type:  "string",
					Value: v.Name,
				},
				{
					Name:  "distance_miles",
					Type:  "int",
					Value: v.DistanceInMiles,
				},
				{
					Name:  "common_like_count",
					Type:  "int",
					Value: v.CommonLikeCount,
				},
				{
					Name:  "common_friend_count",
					Type:  "int",
					Value: v.CommonFriendCount,
				},
				{
					Name:  "ping_time",
					Type:  "string",
					Value: v.PingTime,
				},
			},
		}
		mod.evchan <- ev
	}
}
