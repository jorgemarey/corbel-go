package corbel

// ResourcesService handles the interface for retrival resource's representation
// on Corbel.
//
// Full API info: http://docs.corbelresources.apiary.io/
type ResourcesService struct {
	client *Client
}

//
type ACL map[string]string

// Resource te
type Resource struct {
	ACL map[string]string
}
