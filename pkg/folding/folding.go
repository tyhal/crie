package folding

import "os"

type Folder interface {
	Start(display string, open bool) string
	Stop(id string)
}

func isSet(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

func NewFolder(structured bool) Folder {
	switch {
	case isSet("GITHUB_ACTIONS"):

		return &githubFolder{}
	case isSet("GITLAB_CI"):
		return &gitlabFolder{}
	default:
		return &plainFolder{}
	}
}
