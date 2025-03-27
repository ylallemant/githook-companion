package config

import (
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/command"
	"github.com/ylallemant/githook-companion/pkg/git"
)

func parentVersion(ctx api.ConfigContext, branch string) (string, error) {
	return git.CommitHashFromPath(ctx.ParentPath(), branch)
}

func parentRemoteVersion(ctx api.ConfigContext, branch string) (string, error) {
	rem := gogit.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{ctx.LocalConfig().ParentConfig.GitRepository},
	})

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&gogit.ListOptions{
		// Returns all references, including peeled references.
		PeelingOption: gogit.AppendPeeled,
	})
	if err != nil {
		log.Warn().Msgf("unable to fetch remote hash for branch \"%s\": %s", branch, err.Error())
		return "", nil
	}

	// filters the references list and only the searched branch last commit hash
	var hash string
	for _, ref := range refs {
		if ref.Name().IsBranch() && ref.Name().Short() == branch {
			hash = ref.Hash().String()
		}
	}

	return hash, nil
}

func EnsureVersionSync(ctx api.ConfigContext) error {
	if !ctx.HasParent() {
		// no repository to pull
		return nil
	}

	branch, err := git.CurrentBranchFromPath(ctx.ParentPath())
	if err != nil {
		return err
	}
	log.Debug().Msgf("local parent config branch  \"%s\"", branch)

	localHash, err := parentVersion(ctx, branch)
	if err != nil {
		return err
	}
	log.Debug().Msgf("local parent config version  \"%s\"", localHash)

	remoteHash, err := parentRemoteVersion(ctx, branch)
	if err != nil {
		return err
	}
	log.Debug().Msgf("remote parent config version \"%s\"", remoteHash)

	if remoteHash == "" {
		log.Debug().Msg("skip sync process")
		return nil
	}

	log.Debug().Msgf("parent config has to be synchronized: %v", localHash != remoteHash)
	if localHash != remoteHash {
		err = git.Pull(ctx.ParentPath())
		if err != nil {
			return err
		}

		configInit := command.New("githook-companion")
		configInit.AddArg("init")

		_, err := configInit.Execute()
		if err != nil {
			return errors.Wrapf(err, "failed to init pulled configuration")
		}
	}

	return nil
}
