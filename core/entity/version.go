package entity

type Release struct {
	Name        string `json:"name"`
	Draft       bool   `json:"draft"`
	PreRelease  bool   `json:"pre_release"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
}

type ReleaseSlices []Release

func (r ReleaseSlices) Len() int {
	return len(r)
}

func (r ReleaseSlices) Less(i, j int) bool {
	return r[i].PublishedAt < r[j].PublishedAt
}

func (r ReleaseSlices) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
