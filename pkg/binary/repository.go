package binary

var Repository string

var (
	defaultRepository = "https://github.com/ylallemant/githook-companion"
)

func GetRepository() string {
	return getOr(Repository, defaultRepository)
}
