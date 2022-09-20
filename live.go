package hn

// LiveService communicates with the news
// related pageTypes in the Hacker News API
type LiveService interface {
	Stories(string) ([]int, error)
	MaxItem() (int, error)
	Updates() (*Updates, error)
}

// liveService implements LiveService.
type liveService struct {
	client *Client
}

// Updates contains the latest updated items and profiles
type Updates struct {
	Items    []int    `json:"items"`
	Profiles []string `json:"profiles"`
}

// TopStories gets the ids of the stories on the Front (Top) page, in order.
func (c *Client) TopStories() ([]int, error) {
	return c.Live.Stories("top")
}

// NewStories gets the ids of the stories on the New page, in order.
func (c *Client) NewStories() ([]int, error) {
	return c.Live.Stories("new")
}

// BestStories gets the ids of the stories on the Best page, in order.
func (c *Client) BestStories() ([]int, error) {
	return c.Live.Stories("best")
}

// AskStories gets the ids of the stories on the Ask page, in order.
func (c *Client) AskStories() ([]int, error) {
	return c.Live.Stories("ask")
}

// ShowStories gets the ids of the stories on the Show page, in order.
func (c *Client) ShowStories() ([]int, error) {
	return c.Live.Stories("show")
}

// JobStories gets the ids of the stories on the Show page, in order.
func (c *Client) JobStories() ([]int, error) {
	return c.Live.Stories("job")
}


// Stories gets the ids of the stories for the given pageType, where pageType is one of "top","new","best","ask","show",or "jobs"
func (c *Client) Stories(pageType string) ([]int, error) {
	return c.Live.Stories(pageType)
}


// Stories retrieves the current stories for the given page (where page is one of "top","new","best","ask", or "show")
func (s *liveService) Stories(pageType string) ([]int, error) {
	req, err := s.client.NewRequest(s.path(pageType))
	if err != nil {
		return nil, err
	}

	var value []int
	_, err = s.client.Do(req, &value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *liveService) path(pageType string) string {

	validPageTypes := map[string]string{
		"top":  "topstories.json",
		"new":  "newstories.json",
		"best": "beststories.json",
		"ask":  "askstories.json",
		"show": "showstories.json",
		"job":  "jobstories.json",
	}

	path, ok := validPageTypes[pageType]
	if !ok {
		panic("Invalid pageType: " + pageType)
	}
	return path
}

// MaxItem is a convenience method proxying Live.MaxItem
func (c *Client) MaxItem() (int, error) {
	return c.Live.MaxItem()
}

// MaxItem retrieves the current largest item id
func (s *liveService) MaxItem() (int, error) {
	req, err := s.client.NewRequest(s.maxItemPath())
	if err != nil {
		return 0, err
	}

	var value int
	_, err = s.client.Do(req, &value)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (s *liveService) maxItemPath() string {
	return "maxitem.json"
}

// Updates is a convenience method proxying Live.Updates
func (c *Client) Updates() (*Updates, error) {
	return c.Live.Updates()
}

// Updates retrieves the current largest item id
func (s *liveService) Updates() (*Updates, error) {
	req, err := s.client.NewRequest(s.updatesPath())
	if err != nil {
		return nil, err
	}

	var value Updates
	_, err = s.client.Do(req, &value)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (s *liveService) updatesPath() string {
	return "updates.json"
}
