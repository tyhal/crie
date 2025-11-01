package folding

import "os"

type Folder interface {
	Start(id string)
	Stop()
	Log()
}

func isSet(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

func NewFolder(structured bool) Folder {
	switch {
	case structured:

		return &structuredFolder{}
	case isSet("GITHUB_ACTIONS"):

		return &githubFolder{}
	case isSet("GITLAB_CI"):
		return &gitlabFolder{}
	default:
		return &plainFolder{}
	}
}
