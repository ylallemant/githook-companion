package config

import (
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	gogithttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/command"
	"github.com/ylallemant/githook-companion/pkg/git"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

func parentVersion(ctx api.ConfigContext, branch string) (string, error) {
	return git.CommitHashFromPath(ctx.ParentPath(), branch)
}

func parentRemoteVersion(ctx api.ConfigContext, branch string) (string, error) {
	var err error
	uri := ctx.LocalConfig().ParentConfig.GitRepository

	listOptions := &gogit.ListOptions{
		// Returns all references, including peeled references.
		PeelingOption: gogit.AppendPeeled,
	}

	if ctx.LocalConfig().ParentConfig.Private {
		authMethod, err := git.AuthMethodFromUri(ctx.LocalConfig().ParentConfig.GitRepository)
		if err != nil {
			return "", errors.Wrap(err, "failed to add credentials to repository uri")
		}

		log.Debug().Msgf("add auth method to request options")
		listOptions.Auth = authMethod
	}

	client.InstallProtocol("https", gogithttp.NewClient(globals.DefaultApiClient))

	rem := gogit.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{uri},
	})

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(listOptions)
	if err != nil {
		log.Warn().Msgf("unable to fetch remote hash for branch \"%s\": %s", branch, err.Error())
		SetTimedLockWithDescription("network-problems", networkLockDescription, 2*time.Minute, ctx)
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

	active, _ := TimeLockActive("network-problems", ctx)
	if active {
		// no sync possible
		return nil
	}

	hasCredentials, err := git.HasCredentialsForUri(ctx.LocalConfig().ParentConfig.GitRepository)
	if err != nil {
		return err
	}

	if !hasCredentials {
		active, _ := TimeLockActive("config-sync", ctx)
		if active {
			return nil
		}

		SetTimedLockWithDescription("config-sync", configLockDescription, 10*time.Minute, ctx)
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
