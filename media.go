package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	neturl "net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StoryReelMention represent story reel mention
type StoryReelMention struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        int     `json:"z"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	IsPinned int     `json:"is_pinned"`
	IsHidden int     `json:"is_hidden"`
	User     User
}

// StoryCTA represent story cta
type StoryCTA struct {
	Links []struct {
		LinkType                                int         `json:"linkType"`
		WebURI                                  string      `json:"webUri"`
		AndroidClass                            string      `json:"androidClass"`
		Package                                 string      `json:"package"`
		DeeplinkURI                             string      `json:"deeplinkUri"`
		CallToActionTitle                       string      `json:"callToActionTitle"`
		RedirectURI                             interface{} `json:"redirectUri"`
		LeadGenFormID                           string      `json:"leadGenFormId"`
		IgUserID                                string      `json:"igUserId"`
		AppInstallObjectiveInvalidationBehavior interface{} `json:"appInstallObjectiveInvalidationBehavior"`
	} `json:"links"`
}

// Item represents media items
//
// All Item has Images or Videos objects which contains the url(s).
// You can use Download function to get the best quality Image or Video from Item.
type Item struct {
	media    Media
	Comments *Comments `json:"-"`

	TakenAt          int64   `json:"taken_at"`
	Pk               int64   `json:"pk"`
	ID               string  `json:"id"`
	CommentsDisabled bool    `json:"comments_disabled"`
	DeviceTimestamp  int64   `json:"device_timestamp"`
	MediaType        int     `json:"media_type"`
	Code             string  `json:"code"`
	ClientCacheKey   string  `json:"client_cache_key"`
	FilterType       int     `json:"filter_type"`
	CarouselParentID string  `json:"carousel_parent_id"`
	CarouselMedia    []Item  `json:"carousel_media,omitempty"`
	User             User    `json:"user"`
	CanViewerReshare bool    `json:"can_viewer_reshare"`
	Caption          Caption `json:"caption"`
	CaptionIsEdited  bool    `json:"caption_is_edited"`
	Likes            int     `json:"like_count"`
	HasLiked         bool    `json:"has_liked"`
	// Toplikers can be `string` or `[]string`.
	// Use TopLikers function instead of getting it directly.
	Toplikers                    interface{} `json:"top_likers"`
	Likers                       []User      `json:"likers"`
	CommentLikesEnabled          bool        `json:"comment_likes_enabled"`
	CommentThreadingEnabled      bool        `json:"comment_threading_enabled"`
	HasMoreComments              bool        `json:"has_more_comments"`
	MaxNumVisiblePreviewComments int         `json:"max_num_visible_preview_comments"`
	// Previewcomments can be `string` or `[]string` or `[]Comment`.
	// Use PreviewComments function instead of getting it directly.
	Previewcomments interface{} `json:"preview_comments,omitempty"`
	CommentCount    int         `json:"comment_count"`
	PhotoOfYou      bool        `json:"photo_of_you"`
	// Tags are tagged people in photo
	Tags struct {
		In []Tag `json:"in"`
	} `json:"usertags,omitempty"`
	FbUserTags           Tag    `json:"fb_user_tags"`
	CanViewerSave        bool   `json:"can_viewer_save"`
	OrganicTrackingToken string `json:"organic_tracking_token"`
	// Images contains URL images in different versions.
	// Version = quality.
	Images          Images   `json:"image_versions2,omitempty"`
	OriginalWidth   int      `json:"original_width,omitempty"`
	OriginalHeight  int      `json:"original_height,omitempty"`
	ImportedTakenAt int64    `json:"imported_taken_at,omitempty"`
	Location        Location `json:"location,omitempty"`
	Lat             float64  `json:"lat,omitempty"`
	Lng             float64  `json:"lng,omitempty"`

	// Videos
	Videos            []Video `json:"video_versions,omitempty"`
	HasAudio          bool    `json:"has_audio,omitempty"`
	VideoDuration     float64 `json:"video_duration,omitempty"`
	ViewCount         float64 `json:"view_count,omitempty"`
	IsDashEligible    int     `json:"is_dash_eligible,omitempty"`
	VideoDashManifest string  `json:"video_dash_manifest,omitempty"`
	NumberOfQualities int     `json:"number_of_qualities,omitempty"`

	// Only for stories
	StoryEvents              []interface{}      `json:"story_events"`
	StoryHashtags            []interface{}      `json:"story_hashtags"`
	StoryPolls               []interface{}      `json:"story_polls"`
	StoryFeedMedia           []interface{}      `json:"story_feed_media"`
	StorySoundOn             []interface{}      `json:"story_sound_on"`
	CreativeConfig           interface{}        `json:"creative_config"`
	StoryLocations           []interface{}      `json:"story_locations"`
	StorySliders             []interface{}      `json:"story_sliders"`
	StoryQuestions           []interface{}      `json:"story_questions"`
	StoryProductItems        []interface{}      `json:"story_product_items"`
	StoryCTA                 []StoryCTA         `json:"story_cta"`
	ReelMentions             []StoryReelMention `json:"reel_mentions"`
	SupportsReelReactions    bool               `json:"supports_reel_reactions"`
	ShowOneTapFbShareTooltip bool               `json:"show_one_tap_fb_share_tooltip"`
	HasSharedToFb            int64              `json:"has_shared_to_fb"`
	Mentions                 []Mentions
	Audience                 string `json:"audience,omitempty"`
	StoryMusicStickers       []struct {
		X              float64 `json:"x"`
		Y              float64 `json:"y"`
		Z              int     `json:"z"`
		Width          float64 `json:"width"`
		Height         float64 `json:"height"`
		Rotation       float64 `json:"rotation"`
		IsPinned       int     `json:"is_pinned"`
		IsHidden       int     `json:"is_hidden"`
		IsSticker      int     `json:"is_sticker"`
		MusicAssetInfo struct {
			ID                       string `json:"id"`
			Title                    string `json:"title"`
			Subtitle                 string `json:"subtitle"`
			DisplayArtist            string `json:"display_artist"`
			CoverArtworkURI          string `json:"cover_artwork_uri"`
			CoverArtworkThumbnailURI string `json:"cover_artwork_thumbnail_uri"`
			ProgressiveDownloadURL   string `json:"progressive_download_url"`
			HighlightStartTimesInMs  []int  `json:"highlight_start_times_in_ms"`
			IsExplicit               bool   `json:"is_explicit"`
			DashManifest             string `json:"dash_manifest"`
			HasLyrics                bool   `json:"has_lyrics"`
			AudioAssetID             string `json:"audio_asset_id"`
			IgArtist                 struct {
				Pk            int    `json:"pk"`
				Username      string `json:"username"`
				FullName      string `json:"full_name"`
				IsPrivate     bool   `json:"is_private"`
				ProfilePicURL string `json:"profile_pic_url"`
				ProfilePicID  string `json:"profile_pic_id"`
				IsVerified    bool   `json:"is_verified"`
			} `json:"ig_artist"`
			PlaceholderProfilePicURL string `json:"placeholder_profile_pic_url"`
			ShouldMuteAudio          bool   `json:"should_mute_audio"`
			ShouldMuteAudioReason    string `json:"should_mute_audio_reason"`
			OverlapDurationInMs      int    `json:"overlap_duration_in_ms"`
			AudioAssetStartTimeInMs  int    `json:"audio_asset_start_time_in_ms"`
		} `json:"music_asset_info"`
	} `json:"story_music_stickers,omitempty"`
}

// Comment pushes a text comment to media item.
//
// If parent media is a Story this function will send a private message
// replying the Instagram story.
func (item *Item) Comment(text string) error {
	var opt *reqOptions
	var err error
	insta := item.media.instagram()

	switch item.media.(type) {
	case *StoryMedia:
		to, err := prepareRecipients(item)
		if err != nil {
			return err
		}

		query := insta.prepareDataQuery(
			map[string]interface{}{
				"recipient_users": to,
				"action":          "send_item",
				"media_id":        item.ID,
				"client_context":  generateUUID(),
				"text":            text,
				"entry":           "reel",
				"reel_id":         item.User.ID,
			},
		)
		opt = &reqOptions{
			Connection: "keep-alive",
			Endpoint:   fmt.Sprintf("%s?media_type=%s", urlReplyStory, item.MediaToString()),
			Query:      query,
			IsPost:     true,
		}
	case *FeedMedia: // normal media
		var data string
		data, err = insta.prepareData(
			map[string]interface{}{
				"comment_text": text,
			},
		)
		opt = &reqOptions{
			Endpoint: fmt.Sprintf(urlCommentAdd, item.Pk),
			Query:    generateSignature(data),
			IsPost:   true,
		}
	}
	if err != nil {
		return err
	}

	// ignoring response
	_, err = insta.sendRequest(opt)
	return err
}

// MediaToString returns Item.MediaType as string.
func (item *Item) MediaToString() string {
	switch item.MediaType {
	case 1:
		return "photo"
	case 2:
		return "video"
	case 8:
		return "carousel"
	}
	return ""
}

func setToItem(item *Item, media Media) {
	item.media = media
	item.User.inst = media.instagram()
	item.Comments = newComments(item)
	for i := range item.CarouselMedia {
		item.CarouselMedia[i].User = item.User
		setToItem(&item.CarouselMedia[i], media)
	}
}

func getname(name string) string {
	nname := name
	i := 1
	for {
		ext := path.Ext(name)

		_, err := os.Stat(name)
		if err != nil {
			break
		}
		if ext != "" {
			nname = strings.Replace(nname, ext, "", -1)
		}
		name = fmt.Sprintf("%s.%d%s", nname, i, ext)
		i++
	}
	return name
}

func download(inst *Instagram, url, dst string) (string, error) {
	file, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer file.Close()

	resp, err := inst.c.Get(url)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	return dst, err
}

type bestMedia struct {
	w, h int
	url  string
}

// GetBest returns best quality image or video.
//
// Arguments can be []Video or []Candidate
func GetBest(obj interface{}) string {
	m := bestMedia{}

	switch t := obj.(type) {
	// getting best video
	case []Video:
		for _, video := range t {
			if m.w < video.Width && video.Height > m.h && video.URL != "" {
				m.w = video.Width
				m.h = video.Height
				m.url = video.URL
			}
		}
		// getting best image
	case []Candidate:
		for _, image := range t {
			if m.w < image.Width && image.Height > m.h && image.URL != "" {
				m.w = image.Width
				m.h = image.Height
				m.url = image.URL
			}
		}
	}
	return m.url
}

var rxpTags = regexp.MustCompile(`#\w+`)

// Hashtags returns caption hashtags.
//
// Item media parent must be FeedMedia.
//
// See example: examples/media/hashtags.go
func (item *Item) Hashtags() []Hashtag {
	tags := rxpTags.FindAllString(item.Caption.Text, -1)

	hsh := make([]Hashtag, len(tags))

	i := 0
	for _, tag := range tags {
		hsh[i].Name = tag[1:]
		i++
	}

	for _, comment := range item.PreviewComments() {
		tags := rxpTags.FindAllString(comment.Text, -1)

		for _, tag := range tags {
			hsh = append(hsh, Hashtag{Name: tag[1:]})
		}
	}

	return hsh
}

// Delete deletes your media item. StoryMedia or FeedMedia
//
// See example: examples/media/mediaDelete.go
func (item *Item) Delete() error {
	insta := item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaDelete, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// SyncLikers fetch new likers of a media
//
// This function updates Item.Likers value
func (item *Item) SyncLikers() error {
	resp := respLikers{}
	insta := item.media.instagram()
	body, err := insta.sendSimpleRequest(urlMediaLikers, item.ID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &resp)
	if err == nil {
		item.Likers = resp.Users
	}
	return err
}

// Unlike mark media item as unliked.
//
// See example: examples/media/unlike.go
func (item *Item) Unlike() error {
	insta := item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaUnlike, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Like mark media item as liked.
//
// See example: examples/media/like.go
func (item *Item) Like() error {
	insta := item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaLike, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Save saves media item.
//
// You can get saved media using Account.Saved()
func (item *Item) Save() error {
	insta := item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaSave, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Download downloads media item (video or image) with the best quality.
//
// Input parameters are folder and filename. If filename is "" will be saved with
// the default value name.
//
// If file exists it will be saved
// This function makes folder automatically
//
// This function returns an slice of location of downloaded items
// The returned values are the output path of images and videos.
//
// This function does not download CarouselMedia.
//
// See example: examples/media/itemDownload.go
func (item *Item) Download(folder, name string) (imgs, vds string, err error) {
	var u *neturl.URL
	var nname string
	imgFolder := path.Join(folder, "images")
	vidFolder := path.Join(folder, "videos")
	inst := item.media.instagram()

	os.MkdirAll(folder, 0777)
	os.MkdirAll(imgFolder, 0777)
	os.MkdirAll(vidFolder, 0777)

	vds = GetBest(item.Videos)
	if vds != "" {
		if name == "" {
			u, err = neturl.Parse(vds)
			if err != nil {
				return
			}

			nname = path.Join(vidFolder, path.Base(u.Path))
		} else {
			nname = path.Join(vidFolder, name)
		}
		nname = getname(nname)

		vds, err = download(inst, vds, nname)
		return "", vds, err
	}

	imgs = GetBest(item.Images.Versions)
	if imgs != "" {
		if name == "" {
			u, err = neturl.Parse(imgs)
			if err != nil {
				return
			}

			nname = path.Join(imgFolder, path.Base(u.Path))
		} else {
			nname = path.Join(imgFolder, name)
		}
		nname = getname(nname)

		imgs, err = download(inst, imgs, nname)
		return imgs, "", err
	}

	return imgs, vds, fmt.Errorf("cannot find any image or video")
}

// TopLikers returns string slice or single string (inside string slice)
// Depending on TopLikers parameter.
func (item *Item) TopLikers() []string {
	switch s := item.Toplikers.(type) {
	case string:
		return []string{s}
	case []string:
		return s
	}
	return nil
}

// PreviewComments returns string slice or single string (inside Comment slice)
// Depending on PreviewComments parameter.
// If PreviewComments are string or []string only the Text field will be filled.
func (item *Item) PreviewComments() []Comment {
	switch s := item.Previewcomments.(type) {
	case []interface{}:
		if len(s) == 0 {
			return nil
		}

		switch s[0].(type) {
		case interface{}:
			comments := make([]Comment, 0)
			for i := range s {
				if buf, err := json.Marshal(s[i]); err != nil {
					return nil
				} else {
					comment := &Comment{}

					if err = json.Unmarshal(buf, comment); err != nil {
						return nil
					} else {
						comments = append(comments, *comment)
					}
				}
			}
			return comments
		case string:
			comments := make([]Comment, 0)
			for i := range s {
				comments = append(comments, Comment{
					Text: s[i].(string),
				})
			}
			return comments
		}
	case string:
		comments := []Comment{
			{
				Text: s,
			},
		}
		return comments
	}
	return nil
}

// StoryIsCloseFriends returns a bool
// If the returned value is true the story was published only for close friends
func (item *Item) StoryIsCloseFriends() bool {
	return item.Audience == "besties"
}

//Media interface defines methods for both StoryMedia and FeedMedia.
type Media interface {
	// Next allows pagination
	Next(...interface{}) bool
	// Error returns error (in case it have been occurred)
	Error() error
	// ID returns media id
	ID() string
	// Delete removes media
	Delete() error

	instagram() *Instagram
}

//StoryMedia is the struct that handles the information from the methods to get info about Stories.
type StoryMedia struct {
	inst     *Instagram
	endpoint string
	uid      int64

	err error

	Pk              interface{} `json:"id"`
	LatestReelMedia int64       `json:"latest_reel_media"`
	ExpiringAt      float64     `json:"expiring_at"`
	HaveBeenSeen    float64     `json:"seen"`
	CanReply        bool        `json:"can_reply"`
	Title           string      `json:"title"`
	CanReshare      bool        `json:"can_reshare"`
	ReelType        string      `json:"reel_type"`
	User            User        `json:"user"`
	Items           []Item      `json:"items"`
	ReelMentions    []string    `json:"reel_mentions"`
	PrefetchCount   int         `json:"prefetch_count"`
	// this field can be int or bool
	HasBestiesMedia      interface{} `json:"has_besties_media"`
	StoryRankingToken    string      `json:"story_ranking_token"`
	Broadcasts           []Broadcast `json:"broadcasts"`
	FaceFilterNuxVersion int         `json:"face_filter_nux_version"`
	HasNewNuxStory       bool        `json:"has_new_nux_story"`
	Status               string      `json:"status"`
}

// Delete removes instragram story.
//
// See example: examples/media/deleteStories.go
func (media *StoryMedia) Delete() error {
	insta := media.inst
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": media.ID(),
		},
	)
	if err == nil {
		_, err = insta.sendRequest(
			&reqOptions{
				Endpoint: fmt.Sprintf(urlMediaDelete, media.ID()),
				Query:    generateSignature(data),
				IsPost:   true,
			},
		)
	}
	return err
}

// ID returns Story id
func (media *StoryMedia) ID() string {
	switch id := media.Pk.(type) {
	case int64:
		return strconv.FormatInt(id, 10)
	case string:
		return id
	}
	return ""
}

func (media *StoryMedia) instagram() *Instagram {
	return media.inst
}

func (media *StoryMedia) setValues() {
	for i := range media.Items {
		setToItem(&media.Items[i], media)
	}
}

// Error returns error happened any error
func (media StoryMedia) Error() error {
	return media.err
}

// Seen marks story as seen.
/*
func (media *StoryMedia) Seen() error {
	insta := media.inst
	data, err := insta.prepareData(
		map[string]interface{}{
			"container_module":   "feed_timeline",
			"live_vods_skipped":  "",
			"nuxes_skipped":      "",
			"nuxes":              "",
			"reels":              "", // TODO xd
			"live_vods":          "",
			"reel_media_skipped": "",
		},
	)
	if err == nil {
		_, err = insta.sendRequest(
			&reqOptions{
				Endpoint: urlMediaSeen, // reel=1&live_vod=0
				Query:    generateSignature(data),
				IsPost:   true,
				UseV2:    true,
			},
		)
	}
	return err
}
*/

type trayRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Sync function is used when Highlight must be sync.
// Highlight must be sync when User.Highlights does not return any object inside StoryMedia slice.
//
// This function does NOT update Stories items.
//
// This function updates StoryMedia.Items
func (media *StoryMedia) Sync() error {
	insta := media.inst
	query := []trayRequest{
		{"SUPPORTED_SDK_VERSIONS", "9.0,10.0,11.0,12.0,13.0,14.0,15.0,16.0,17.0,18.0,19.0,20.0,21.0,22.0,23.0,24.0"},
		{"FACE_TRACKER_VERSION", "10"},
		{"segmentation", "segmentation_enabled"},
		{"COMPRESSION", "ETC2_COMPRESSION"},
	}
	qjson, err := json.Marshal(query)
	if err != nil {
		return err
	}

	id := media.Pk.(string)
	data, err := insta.prepareData(
		map[string]interface{}{
			"user_ids":                   []string{id},
			"supported_capabilities_new": b2s(qjson),
		},
	)
	if err != nil {
		return err
	}

	body, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlReelMedia,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err == nil {
		resp := trayResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			m, ok := resp.Reels[id]
			if ok {
				media.Items = m.Items
				media.setValues()
				return nil
			}
			err = fmt.Errorf("cannot find %s structure in response", id)
		}
	}
	return err
}

// Next allows pagination after calling:
// User.Stories
//
//
// returns false when list reach the end
// if StoryMedia.Error() is ErrNoMore no problem have been occurred.
func (media *StoryMedia) Next(params ...interface{}) bool {
	if media.err != nil {
		return false
	}

	insta := media.inst
	endpoint := media.endpoint
	if media.uid != 0 {
		endpoint = fmt.Sprintf(endpoint, media.uid)
	}

	body, err := insta.sendSimpleRequest(endpoint)
	if err == nil {
		m := StoryMedia{}
		err = json.Unmarshal(body, &m)
		if err == nil {
			// TODO check NextID media
			*media = m
			media.inst = insta
			media.endpoint = endpoint
			media.err = ErrNoMore // TODO: See if stories has pagination
			media.setValues()
			return true
		}
	}
	media.err = err
	return false
}

type TimelineItem struct {
	// MediaOrAd MediaOrAd `json:"media_or_ad"`
	MediaOrAd Item `json:"media_or_ad"`
}
type SharingFrictionInfo struct {
	ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
	BloksAppURL               interface{} `json:"bloks_app_url"`
}
type CarouselMedia struct {
	ID                    string              `json:"id"`
	MediaType             int                 `json:"media_type"`
	OriginalWidth         int                 `json:"original_width"`
	OriginalHeight        int                 `json:"original_height"`
	Pk                    int64               `json:"pk"`
	CarouselParentID      string              `json:"carousel_parent_id"`
	CanSeeInsightsAsBrand bool                `json:"can_see_insights_as_brand"`
	SharingFrictionInfo   SharingFrictionInfo `json:"sharing_friction_info"`
}
type FriendshipStatus struct {
	Following       bool `json:"following"`
	OutgoingRequest bool `json:"outgoing_request"`
	IsMutingReel    bool `json:"is_muting_reel"`
	IsBestie        bool `json:"is_bestie"`
	IsRestricted    bool `json:"is_restricted"`
}
type PreviewComments struct {
	Pk               int64  `json:"pk"`
	UserID           int64  `json:"user_id"`
	Text             string `json:"text"`
	Type             int    `json:"type"`
	CreatedAt        int    `json:"created_at"`
	CreatedAtUtc     int    `json:"created_at_utc"`
	ContentType      string `json:"content_type"`
	Status           string `json:"status"`
	BitFlags         int    `json:"bit_flags"`
	DidReportAsSpam  bool   `json:"did_report_as_spam"`
	ShareEnabled     bool   `json:"share_enabled"`
	User             User   `json:"user,omitempty"`
	IsCovered        bool   `json:"is_covered"`
	MediaID          int64  `json:"media_id"`
	HasTranslation   bool   `json:"has_translation"`
	HasLikedComment  bool   `json:"has_liked_comment"`
	CommentLikeCount int    `json:"comment_like_count"`
	ParentCommentID  int64  `json:"parent_comment_id,omitempty"`
}
type MediaOrAd struct {
	TakenAt                        int64                 `json:"taken_at"`
	Pk                             int64                 `json:"pk"`
	ID                             string                `json:"id"`
	DeviceTimestamp                int64                 `json:"device_timestamp"`
	MediaType                      int                   `json:"media_type"`
	Code                           string                `json:"code"`
	ClientCacheKey                 string                `json:"client_cache_key"`
	FilterType                     int                   `json:"filter_type"`
	CarouselMediaCount             int                   `json:"carousel_media_count"`
	CarouselMedia                  []CarouselMedia       `json:"carousel_media"`
	CanSeeInsightsAsBrand          bool                  `json:"can_see_insights_as_brand"`
	ShouldRequestAds               bool                  `json:"should_request_ads"`
	OriginalWidth                  int                   `json:"original_width"`
	OriginalHeight                 int                   `json:"original_height"`
	User                           User                  `json:"user"`
	CanViewerReshare               bool                  `json:"can_viewer_reshare"`
	CaptionIsEdited                bool                  `json:"caption_is_edited"`
	CommentLikesEnabled            bool                  `json:"comment_likes_enabled"`
	CommentThreadingEnabled        bool                  `json:"comment_threading_enabled"`
	HasMoreComments                bool                  `json:"has_more_comments"`
	NextMaxID                      int64                 `json:"next_max_id"`
	MaxNumVisiblePreviewComments   int                   `json:"max_num_visible_preview_comments"`
	PreviewComments                []PreviewComments     `json:"preview_comments"`
	CanViewMorePreviewComments     bool                  `json:"can_view_more_preview_comments"`
	CommentCount                   int                   `json:"comment_count"`
	InlineComposerDisplayCondition string                `json:"inline_composer_display_condition"`
	InlineComposerImpTriggerTime   int                   `json:"inline_composer_imp_trigger_time"`
	LikeCount                      int                   `json:"like_count"`
	HasLiked                       bool                  `json:"has_liked"`
	TopLikers                      []interface{}         `json:"top_likers"`
	PhotoOfYou                     bool                  `json:"photo_of_you"`
	Caption                        Caption               `json:"caption"`
	Injected                       Injected              `json:"injected"`
	CollapseComments               bool                  `json:"collapse_comments"`
	AdMetadata                     []AdMetadata          `json:"ad_metadata"`
	Link                           string                `json:"link"`
	LinkText                       string                `json:"link_text"`
	AdAction                       string                `json:"ad_action"`
	LinkHintText                   string                `json:"link_hint_text"`
	ITunesItem                     interface{}           `json:"iTunesItem"`
	AdLinkType                     int                   `json:"ad_link_type"`
	AdHeaderStyle                  int                   `json:"ad_header_style"`
	DrAdType                       int                   `json:"dr_ad_type"`
	AndroidLinks                   []AndroidLinks        `json:"android_links"`
	IabAutofillOptoutInfo          IabAutofillOptoutInfo `json:"iab_autofill_optout_info"`
	ForceOverlay                   bool                  `json:"force_overlay"`
	HideNuxText                    bool                  `json:"hide_nux_text"`
	OverlayText                    string                `json:"overlay_text"`
	OverlayTitle                   string                `json:"overlay_title"`
	OverlaySubtitle                string                `json:"overlay_subtitle"`
	DominantColor                  string                `json:"dominant_color"`
	FollowerCount                  int                   `json:"follower_count"`
	PostCount                      int                   `json:"post_count"`
	FbPageURL                      string                `json:"fb_page_url"`
	CanViewerSave                  bool                  `json:"can_viewer_save"`
	OrganicTrackingToken           string                `json:"organic_tracking_token"`
	ExpiringAt                     int                   `json:"expiring_at"`
	Preview                        string                `json:"preview"`
	SharingFrictionInfo            SharingFrictionInfo   `json:"sharing_friction_info"`
	IsInProfileGrid                bool                  `json:"is_in_profile_grid"`
	ProfileGridControlEnabled      bool                  `json:"profile_grid_control_enabled"`
	IsShopTheLookEligible          bool                  `json:"is_shop_the_look_eligible"`
	DeletedReason                  int                   `json:"deleted_reason"`
	InventorySource                string                `json:"inventory_source"`
	IsSeen                         bool                  `json:"is_seen"`
	IsEOF                          bool                  `json:"is_eof"`
	CommentsDisabled               bool                  `json:"comments_disabled"`
	Comments                       []interface{}         `json:"comments"`
}

type OrganicBidInfo struct {
	OrganicBidCpm                 float64     `json:"organicBidCpm"`
	UserDollarValue               float64     `json:"userDollarValue"`
	SerializedPacedOrganicBidMap  string      `json:"serializedPacedOrganicBidMap"`
	SerializedOrganicEventProbMap string      `json:"serializedOrganicEventProbMap"`
	SerializedEventQualityMap     string      `json:"serializedEventQualityMap"`
	SerializedImpOrganicBidMap    string      `json:"serializedImpOrganicBidMap"`
	SerializedImpOrganicCoeffMap  interface{} `json:"serializedImpOrganicCoeffMap"`
	SerializedImpOrganicScoreMap  interface{} `json:"serializedImpOrganicScoreMap"`
}
type AdsBidInfo struct {
	EcpmBid        float64        `json:"ecpmBid"`
	EcpmPrice      float64        `json:"ecpmPrice"`
	Ectr           float64        `json:"ectr"`
	Ecvr           float64        `json:"ecvr"`
	PostImpEcvr    int            `json:"postImpEcvr"`
	OrganicBidInfo OrganicBidInfo `json:"organicBidInfo"`
}
type AdsRankingInfo struct {
	OrganicRank  int `json:"organicRank"`
	EcpsBidRank  int `json:"ecpsBidRank"`
	TotalBidRank int `json:"totalBidRank"`
}
type AdsDebugInfo struct {
	AdsBidInfo     AdsBidInfo     `json:"adsBidInfo"`
	AdsRankingInfo AdsRankingInfo `json:"adsRankingInfo"`
}
type HideReasonsV2 struct {
	Text   string `json:"text"`
	Reason string `json:"reason"`
}
type CtdAdsInfo struct {
	BusinessResponsivenessTimeText string      `json:"business_responsiveness_time_text"`
	WelcomeMessageText             interface{} `json:"welcome_message_text"`
}
type Injected struct {
	Label                           string          `json:"label"`
	ShowIcon                        bool            `json:"show_icon"`
	HideLabel                       string          `json:"hide_label"`
	Invalidation                    interface{}     `json:"invalidation"`
	IsDemo                          bool            `json:"is_demo"`
	ViewTags                        []interface{}   `json:"view_tags"`
	IsHoldout                       bool            `json:"is_holdout"`
	TrackingToken                   string          `json:"tracking_token"`
	ShowAdChoices                   bool            `json:"show_ad_choices"`
	AdTitle                         string          `json:"ad_title"`
	AboutAdParams                   string          `json:"about_ad_params"`
	DirectShare                     bool            `json:"direct_share"`
	AdID                            int64           `json:"ad_id"`
	DisplayViewabilityEligible      bool            `json:"display_viewability_eligible"`
	AdsDebugInfo                    AdsDebugInfo    `json:"ads_debug_info"`
	ShouldShowSecondaryCtaOnProfile bool            `json:"should_show_secondary_cta_on_profile"`
	HideReasonsV2                   []HideReasonsV2 `json:"hide_reasons_v2"`
	HideFlowType                    int             `json:"hide_flow_type"`
	Cookies                         []string        `json:"cookies"`
	CtdAdsInfo                      CtdAdsInfo      `json:"ctd_ads_info"`
}
type AdMetadata struct {
	Value string `json:"value"`
	Type  int    `json:"type"`
}
type AndroidLinks struct {
	LinkType                                int         `json:"linkType"`
	WebURI                                  string      `json:"webUri"`
	AndroidClass                            string      `json:"androidClass"`
	Package                                 string      `json:"package"`
	DeeplinkURI                             string      `json:"deeplinkUri"`
	CallToActionTitle                       string      `json:"callToActionTitle"`
	RedirectURI                             string      `json:"redirectUri"`
	LeadGenFormID                           interface{} `json:"leadGenFormId"`
	IgUserID                                interface{} `json:"igUserId"`
	AppInstallObjectiveInvalidationBehavior interface{} `json:"appInstallObjectiveInvalidationBehavior"`
}
type IabAutofillOptoutInfo struct {
	Domain              string `json:"domain"`
	IsIabAutofillOptout bool   `json:"is_iab_autofill_optout"`
}
type EndOfFeedDemarcator struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}
type StoriesNetego struct {
	TrackingToken  string `json:"tracking_token"`
	HideUnitIfSeen string `json:"hide_unit_if_seen"`
	ClientPosition int    `json:"client_position"`
	ID             int64  `json:"id"`
}

// FeedMedia represent a set of media items
type FeedMedia struct {
	inst *Instagram

	err error

	uid       int64
	endpoint  string
	timestamp string

	Items               []Item         `json:"items"`
	TimelineItems       []TimelineItem `json:"feed_items"`
	NumResults          int            `json:"num_results"`
	MoreAvailable       bool           `json:"more_available"`
	AutoLoadMoreEnabled bool           `json:"auto_load_more_enabled"`
	Status              string         `json:"status"`
	// Can be int64 and string
	// this is why we recommend Next() usage :')
	NextID interface{} `json:"next_max_id"`

	IsTimelineMedia bool
}

// Delete deletes all items in media. Take care...
//
// See example: examples/media/mediaDelete.go
func (media *FeedMedia) Delete() error {
	for i := range media.Items {
		media.Items[i].Delete()
	}
	return nil
}

func (media *FeedMedia) instagram() *Instagram {
	return media.inst
}

// SetInstagram set instagram
func (media *FeedMedia) SetInstagram(inst *Instagram) {
	media.inst = inst
}

// SetID sets media ID
// this value can be int64 or string
func (media *FeedMedia) SetID(id interface{}) {
	media.NextID = id
}

// Sync updates media values.
func (media *FeedMedia) Sync() error {
	id := media.ID()
	insta := media.inst

	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": id,
		},
	)
	if err != nil {
		return err
	}

	body, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaInfo, id),
			Query:    generateSignature(data),
			IsPost:   false,
		},
	)
	if err != nil {
		return err
	}

	m := FeedMedia{}
	err = json.Unmarshal(body, &m)
	*media = m
	media.endpoint = urlMediaInfo
	media.inst = insta
	media.NextID = id
	media.setValues()
	return err
}

func (media *FeedMedia) setValues() {
	for i := range media.Items {
		setToItem(&media.Items[i], media)
	}
}

func (media FeedMedia) Error() error {
	return media.err
}

// ID returns media id.
func (media *FeedMedia) ID() string {
	switch s := media.NextID.(type) {
	case string:
		return s
	case int64:
		return strconv.FormatInt(s, 10)
	case json.Number:
		return string(s)
	}
	return ""
}

// Next allows pagination after calling:
// User.Feed
// Params: ranked_content is set to "true" by default, you can set it to false by either passing "false" or false as parameter.
// returns false when list reach the end.
// if FeedMedia.Error() is ErrNoMore no problem have been occurred.
func (media *FeedMedia) Next(params ...interface{}) bool {
	if media.err != nil {
		return false
	}

	insta := media.inst
	endpoint := media.endpoint
	next := media.ID()
	ranked := "true"

	if media.uid != 0 {
		endpoint = fmt.Sprintf(endpoint, media.uid)
	}

	for _, param := range params {
		switch s := param.(type) {
		case string:
			if _, err := strconv.ParseBool(s); err == nil {
				ranked = s
			}
		case bool:
			if !s {
				ranked = "false"
			}
		}
	}
	body, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			Query: map[string]string{
				"max_id":         next,
				"rank_token":     insta.rankToken,
				"min_timestamp":  media.timestamp,
				"ranked_content": ranked,
			},
			IsPost: media.IsTimelineMedia,
		},
	)
	if err == nil {
		m := FeedMedia{IsTimelineMedia: media.IsTimelineMedia}
		d := json.NewDecoder(bytes.NewReader(body))
		d.UseNumber()
		err = d.Decode(&m)
		if err == nil {
			*media = m
			media.inst = insta
			media.endpoint = endpoint
			if m.NextID == 0 || !m.MoreAvailable {
				media.err = ErrNoMore
			}
			media.setValues()
			return true
		}
	}
	return false
}

// UploadPhoto post image from io.Reader to instagram.
func (insta *Instagram) UploadPhoto(photo io.Reader, photoCaption string, quality int, filterType int) (Item, error) {
	out := Item{}

	config, err := insta.postPhoto(photo, photoCaption, quality, filterType, false)
	if err != nil {
		return out, err
	}
	data, err := insta.prepareData(config)
	if err != nil {
		return out, err
	}

	body, err := insta.sendRequest(&reqOptions{
		Endpoint: "media/configure/?",
		Query:    generateSignature(data),
		IsPost:   true,
	})
	if err != nil {
		return out, err
	}
	var uploadResult struct {
		Media    Item   `json:"media"`
		UploadID string `json:"upload_id"`
		Status   string `json:"status"`
	}
	err = json.Unmarshal(body, &uploadResult)
	if err != nil {
		return out, err
	}

	if uploadResult.Status != "ok" {
		return out, fmt.Errorf("invalid status, result: %s", uploadResult.Status)
	}

	return uploadResult.Media, nil
}

func (insta *Instagram) postPhoto(photo io.Reader, photoCaption string, quality int, filterType int, isSidecar bool) (map[string]interface{}, error) {
	uploadID := time.Now().Unix()
	photoName := fmt.Sprintf("pending_media_%d.jpg", uploadID)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	w.WriteField("upload_id", strconv.FormatInt(uploadID, 10))
	w.WriteField("_uuid", insta.uuid)
	w.WriteField("_csrftoken", insta.token)
	var compression = map[string]interface{}{
		"lib_name":    "jt",
		"lib_version": "1.3.0",
		"quality":     quality,
	}
	cBytes, _ := json.Marshal(compression)
	w.WriteField("image_compression", toString(cBytes))
	if isSidecar {
		w.WriteField("is_sidecar", toString(1))
	}
	fw, err := w.CreateFormFile("photo", photoName)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	rdr := io.TeeReader(photo, &buf)
	if _, err = io.Copy(fw, rdr); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", goInstaAPIUrl+"upload/photo/", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-IG-Capabilities", "3Q4=")
	req.Header.Set("X-IG-Connection-Type", "WIFI")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-type", w.FormDataContentType())
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", goInstaUserAgent)

	resp, err := insta.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code, result: %s", resp.Status)
	}
	var result struct {
		UploadID       string      `json:"upload_id"`
		XsharingNonces interface{} `json:"xsharing_nonces"`
		Status         string      `json:"status"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != "ok" {
		return nil, fmt.Errorf("unknown error, status: %s", result.Status)
	}
	width, height, err := getImageDimensionFromReader(&buf)
	if err != nil {
		return nil, err
	}
	config := map[string]interface{}{
		"media_folder": "Instagram",
		"source_type":  4,
		"caption":      photoCaption,
		"upload_id":    strconv.FormatInt(uploadID, 10),
		"device":       goInstaDeviceSettings,
		"edits": map[string]interface{}{
			"crop_original_size": []int{width * 1.0, height * 1.0},
			"crop_center":        []float32{0.0, 0.0},
			"crop_zoom":          1.0,
			"filter_type":        filterType,
		},
		"extra": map[string]interface{}{
			"source_width":  width,
			"source_height": height,
		},
	}
	return config, nil
}

// UploadAlbum post image from io.Reader to instagram.
func (insta *Instagram) UploadAlbum(photos []io.Reader, photoCaption string, quality int, filterType int) (Item, error) {
	out := Item{}

	var childrenMetadata []map[string]interface{}
	for _, photo := range photos {
		config, err := insta.postPhoto(photo, photoCaption, quality, filterType, true)
		if err != nil {
			return out, err
		}

		childrenMetadata = append(childrenMetadata, config)
	}
	albumUploadID := time.Now().Unix()

	config := map[string]interface{}{
		"caption":           photoCaption,
		"client_sidecar_id": albumUploadID,
		"children_metadata": childrenMetadata,
	}
	data, err := insta.prepareData(config)
	if err != nil {
		return out, err
	}

	body, err := insta.sendRequest(&reqOptions{
		Endpoint: "media/configure_sidecar/?",
		Query:    generateSignature(data),
		IsPost:   true,
	})
	if err != nil {
		return out, err
	}

	var uploadResult struct {
		Media           Item   `json:"media"`
		ClientSideCarID int64  `json:"client_sidecar_id"`
		Status          string `json:"status"`
	}
	err = json.Unmarshal(body, &uploadResult)
	if err != nil {
		return out, err
	}

	if uploadResult.Status != "ok" {
		return out, fmt.Errorf("invalid status, result: %s", uploadResult.Status)
	}

	return uploadResult.Media, nil
}
