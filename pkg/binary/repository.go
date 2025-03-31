package binary

var repository string

var (
	defaultRepository = "https://github.com/ylallemant/githook-companion"
)

func GetRepository() string {
	return getOr(repository, defaultRepository)
}
