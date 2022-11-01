package hn

import "context"

// LiveService communicates with the news
// related pageTypes in the Hacker News API
type LiveService interface {
	Stories(context.Context, string) ([]int, error)
	MaxItem(context.Context) (int, error)
	Updates(context.Context) (*Updates, error)
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
func (c *Client) TopStories(ctx context.Context) ([]int, error) {
	return c.Live.Stories(ctx, "top")
}

// NewStories gets the ids of the stories on the New page, in order.
func (c *Client) NewStories(ctx context.Context) ([]int, error) {
	return c.Live.Stories(ctx, "new")
}

// BestStories gets the ids of the stories on the Best page, in order.
func (c *Client) BestStories(ctx context.Context) ([]int, error) {
	return c.Live.Stories(ctx, "best")
}

// AskStories gets the ids of the stories on the Ask page, in order.
func (c *Client) AskStories(ctx context.Context) ([]int, error) {
	return c.Live.Stories(ctx, "ask")
}

// ShowStories gets the ids of the stories on the Show page, in order.
func (c *Client) ShowStories(ctx context.Context) ([]int, error) {
	return c.Live.Stories(ctx, "show")
}

// JobStories gets the ids of the stories on the Show page, in order.
func (c *Client) JobStories(ctx context.Context) ([]int, error) {
	return c.Live.Stories(ctx, "job")
}

// Stories gets the ids of the stories for the given pageType, where pageType is one of "top","new","best","ask","show",or "jobs"
func (c *Client) Stories(ctx context.Context, pageType string) ([]int, error) {
	return c.Live.Stories(ctx, pageType)
}

// Stories retrieves the current stories for the given page (where page is one of "top","new","best","ask", or "show")
func (s *liveService) Stories(ctx context.Context, pageType string) ([]int, error) {
	req, err := s.client.NewRequest(ctx, s.path(pageType))
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
func (c *Client) MaxItem(ctx context.Context) (int, error) {
	return c.Live.MaxItem(ctx)
}

// MaxItem retrieves the current largest item id
func (s *liveService) MaxItem(ctx context.Context) (int, error) {
	req, err := s.client.NewRequest(ctx, s.maxItemPath())
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
func (c *Client) Updates(ctx context.Context) (*Updates, error) {
	return c.Live.Updates(ctx)
}

// Updates retrieves the current largest item id
func (s *liveService) Updates(ctx context.Context) (*Updates, error) {
	req, err := s.client.NewRequest(ctx, s.updatesPath())
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
