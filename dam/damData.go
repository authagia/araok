package dam

const DamAPI = "https://www.clubdam.com/dkwebsys/search-api/SearchMusicByKeywordApi"

type SearchResponse struct {
	Result Result `json:"result"`
	Data   Data   `json:"data"`
	List   []Song `json:"list"`
}

type Result struct {
	StatusCode string `json:"statusCode"`
	Message    string `json:"message"`
}

type Data struct {
	Keyword    string `json:"keyword"`
	PageNo     int    `json:"pageNo"`
	PageCount  int    `json:"pageCount"`
	HasNext    string `json:"hasNext"`
	HasPreview string `json:"hasPreview"`
	TotalCount int    `json:"totalCount"`
	OverFlag   string `json:"overFlag"`
}

type Song struct {
	RequestNo         string `json:"requestNo"`
	Title             string `json:"title"`
	TitleYomi         string `json:"titleYomi"`
	ArtistCode        int    `json:"artistCode"`
	Artist            string `json:"artist"`
	ArtistYomi        string `json:"artistYomi"`
	ReleaseDate       string `json:"releaseDate"`
	NewReleaseFlag    string `json:"newReleaseFlag"`
	FutureReleaseFlag string `json:"futureReleaseFlag"`

	HighlightLyrics string `json:"highlightLyrics"`

	HonninFlag     string `json:"honninFlag"`
	AnimeFlag      string `json:"animeFlag"`
	LiveFlag       string `json:"liveFlag"`
	KidsFlag       string `json:"kidsFlag"`
	MamaotoFlag    string `json:"mamaotoFlag"`
	NamaotoFlag    string `json:"namaotoFlag"`
	DuetFlag       string `json:"duetFlag"`
	GuideVocalFlag string `json:"guideVocalFlag"`
	ProokeFlag     string `json:"prookeFlag"`
	DuetDxFlag     string `json:"duetDxFlag"`

	DamTomoMovieFlag           string `json:"damTomoMovieFlag"`
	DamTomoRecordingFlag       string `json:"damTomoRecordingFlag"`
	DamTomoPublicVocalFlag     string `json:"damTomoPublicVocalFlag"`
	DamTomoPublicMovieFlag     string `json:"damTomoPublicMovieFlag"`
	DamTomoPublicRecordingFlag string `json:"damTomoPublicRecordingFlag"`

	ScoreFlag    string `json:"scoreFlag"`
	MyListFlag   string `json:"myListFlag"`
	Shift        string `json:"shift"`
	PlaybackTime int    `json:"playbackTime"`
}
