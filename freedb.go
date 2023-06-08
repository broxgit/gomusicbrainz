package gomusicbrainz

type Freedb struct{}

func (c *WS2Client) SearchFreedb(searchTerm string, limit, offset int) (*FreedbSearchResponse, error) {
	return nil, nil
}

type FreedbSearchResponse struct{}
