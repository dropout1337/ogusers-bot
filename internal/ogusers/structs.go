package ogusers

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`

	Uid        string `json:"uid"`
	Username   string `json:"username"`
	Rank       string `json:"rank"`
	Avatar     string `json:"avatar"`
	Registered string `json:"registered"`
	Signature  string `json:"signature"`
	Revolution bool   `json:"revolution"`
	Luminary   bool   `json:"luminary"`
	Rep        string `json:"rep"`
	Vouches    string `json:"vouches"`
	Likes      string `json:"likes"`
	Awards     string `json:"awards"`
	Items      string `json:"items"`
	Threads    string `json:"threads"`
	Posts      string `json:"posts"`
	Visitors   string `json:"visitors"`
	Credits    string `json:"credits"`
	Title      string `json:"title"`
	Lastonline string `json:"lastonline"`
	Ban        bool   `json:"ban"`
}
