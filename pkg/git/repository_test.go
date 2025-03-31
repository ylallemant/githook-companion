package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseGitURI(t *testing.T) {
	cases := []struct {
		name                 string
		uri                  string
		expectedProtocol     string
		expectedHostname     string
		expectedPath         string
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name:             "Azure style uri",
			uri:              "https://PROJECT@dev.azure.com/PROJECT/test/_git/some-repo",
			expectedProtocol: "https",
			expectedHostname: "dev.azure.com",
			expectedPath:     "/PROJECT/test/_git/some-repo",
			expectError:      false,
		},
		{
			name:             "Azure style ssh uri",
			uri:              "git@ssh.dev.azure.com:v3/PROJECT/test/some-repo",
			expectedProtocol: "https",
			expectedHostname: "dev.azure.com",
			expectedPath:     "/PROJECT/test/_git/some-repo",
			expectError:      false,
		},
		{
			name:             "git protocol uri",
			uri:              "git://github.com/test/some-repo.git",
			expectedProtocol: "https",
			expectedHostname: "github.com",
			expectedPath:     "/test/some-repo",
			expectError:      false,
		},
		{
			name:             "git ssh uri",
			uri:              "git@github.com:test/some-repo.git",
			expectedProtocol: "https",
			expectedHostname: "github.com",
			expectedPath:     "/test/some-repo",
			expectError:      false,
		},
		{
			name:             "git uri",
			uri:              "https://github.com/test/some-repo.git",
			expectedProtocol: "https",
			expectedHostname: "github.com",
			expectedPath:     "/test/some-repo",
			expectError:      false,
		},
		{
			name:             "clean git uri",
			uri:              "https://github.com/test/some-repo",
			expectedProtocol: "https",
			expectedHostname: "github.com",
			expectedPath:     "/test/some-repo",
			expectError:      false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			parsed, err := parseGitURI(c.uri)

			if c.expectError {
				assert.Nil(tt, parsed)
				assert.NotNil(tt, err)
				assert.Equal(tt, c.expectedErrorMessage, err.Error(), "wrong error message")
			} else {
				assert.NotNil(tt, parsed)
				assert.Nil(tt, err)

				assert.Equal(tt, c.expectedProtocol, parsed.Scheme, "wrong protocol")
				assert.Equal(tt, c.expectedHostname, parsed.Hostname(), "wrong hostname")
				assert.Equal(tt, c.expectedPath, parsed.Path, "wrong path")
			}
		})
	}
}

func Test_RepositoryFromUri(t *testing.T) {
	cases := []struct {
		name                 string
		uri                  string
		expected             string
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name:        "Azure style uri",
			uri:         "https://PROJECT@dev.azure.com/PROJECT/test/_git/some-repo",
			expected:    "some-repo",
			expectError: false,
		},
		{
			name:        "Azure style ssh uri",
			uri:         "git@ssh.dev.azure.com:v3/PROJECT/test/some-repo",
			expected:    "some-repo",
			expectError: false,
		},
		{
			name:        "git protocol uri",
			uri:         "git://github.com/test/some-repo.git",
			expected:    "some-repo",
			expectError: false,
		},
		{
			name:        "git ssh uri",
			uri:         "git@github.com:test/some-repo.git",
			expected:    "some-repo",
			expectError: false,
		},
		{
			name:        "git uri",
			uri:         "https://github.com/test/some-repo.git",
			expected:    "some-repo",
			expectError: false,
		},
		{
			name:        "clean git uri",
			uri:         "https://github.com/test/some-repo",
			expected:    "some-repo",
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			repository, err := RepositoryFromUri(c.uri)

			if c.expectError {
				assert.NotNil(tt, err)
				assert.Equal(tt, c.expectedErrorMessage, err.Error(), "wrong error message")
			} else {
				assert.Nil(tt, err)

				assert.Equal(tt, c.expected, repository, "wrong repository")
			}
		})
	}
}

func Test_OwnerFromUri(t *testing.T) {
	cases := []struct {
		name                 string
		uri                  string
		expected             string
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name:        "Azure style uri",
			uri:         "https://PROJECT@dev.azure.com/PROJECT/test/_git/some-repo",
			expected:    "PROJECT",
			expectError: false,
		},
		{
			name:        "Azure style ssh uri",
			uri:         "git@ssh.dev.azure.com:v3/PROJECT/test/some-repo",
			expected:    "PROJECT",
			expectError: false,
		},
		{
			name:        "git protocol uri",
			uri:         "git://github.com/test/some-repo.git",
			expected:    "test",
			expectError: false,
		},
		{
			name:        "git ssh uri",
			uri:         "git@github.com:test/some-repo.git",
			expected:    "test",
			expectError: false,
		},
		{
			name:        "git uri",
			uri:         "https://github.com/test/some-repo.git",
			expected:    "test",
			expectError: false,
		},
		{
			name:        "clean git uri",
			uri:         "https://github.com/test/some-repo",
			expected:    "test",
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			repository, err := OwnerFromUri(c.uri)

			if c.expectError {
				assert.NotNil(tt, err)
				assert.Equal(tt, c.expectedErrorMessage, err.Error(), "wrong error message")
			} else {
				assert.Nil(tt, err)

				assert.Equal(tt, c.expected, repository, "wrong repository")
			}
		})
	}
}

func Test_OwnerAndRepositoryFromUri(t *testing.T) {
	cases := []struct {
		name                 string
		uri                  string
		expectedOwner        string
		expectedRepository   string
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name:               "Azure style uri",
			uri:                "https://PROJECT@dev.azure.com/PROJECT/test/_git/some-repo",
			expectedOwner:      "PROJECT",
			expectedRepository: "some-repo",
			expectError:        false,
		},
		{
			name:               "Azure style ssh uri",
			uri:                "git@ssh.dev.azure.com:v3/PROJECT/test/some-repo",
			expectedOwner:      "PROJECT",
			expectedRepository: "some-repo",
			expectError:        false,
		},
		{
			name:               "git protocol uri",
			uri:                "git://github.com/test/some-repo.git",
			expectedOwner:      "test",
			expectedRepository: "some-repo",
			expectError:        false,
		},
		{
			name:               "git ssh uri",
			uri:                "git@github.com:test/some-repo.git",
			expectedOwner:      "test",
			expectedRepository: "some-repo",
			expectError:        false,
		},
		{
			name:               "git uri",
			uri:                "https://github.com/test/some-repo.git",
			expectedOwner:      "test",
			expectedRepository: "some-repo",
			expectError:        false,
		},
		{
			name:               "clean git uri",
			uri:                "https://github.com/test/some-repo",
			expectedOwner:      "test",
			expectedRepository: "some-repo",
			expectError:        false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			owner, repository, err := OwnerAndRepositoryFromUri(c.uri)

			if c.expectError {
				assert.NotNil(tt, err)
				assert.Equal(tt, c.expectedErrorMessage, err.Error(), "wrong error message")
			} else {
				assert.Nil(tt, err)

				assert.Equal(tt, c.expectedOwner, owner, "wrong repository")
				assert.Equal(tt, c.expectedRepository, repository, "wrong repository")
			}
		})
	}
}

func Test_RepositorySignatureFromUri(t *testing.T) {
	cases := []struct {
		name                 string
		uri                  string
		expected             string
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name:        "Azure style uri",
			uri:         "https://PROJECT@dev.azure.com/PROJECT/test/_git/some-repo",
			expected:    "dev.azure.com/PROJECT/test/_git/some-repo",
			expectError: false,
		},
		{
			name:        "Azure style ssh uri",
			uri:         "git@ssh.dev.azure.com:v3/PROJECT/test/some-repo",
			expected:    "dev.azure.com/PROJECT/test/_git/some-repo",
			expectError: false,
		},
		{
			name:        "git protocol uri",
			uri:         "git://github.com/test/some-repo.git",
			expected:    "github.com/test/some-repo",
			expectError: false,
		},
		{
			name:        "git ssh uri",
			uri:         "git@github.com:test/some-repo.git",
			expected:    "github.com/test/some-repo",
			expectError: false,
		},
		{
			name:        "git uri",
			uri:         "https://github.com/test/some-repo.git",
			expected:    "github.com/test/some-repo",
			expectError: false,
		},
		{
			name:        "clean git uri",
			uri:         "https://github.com/test/some-repo",
			expected:    "github.com/test/some-repo",
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			repository, err := RepositorySignatureFromUri(c.uri)

			if c.expectError {
				assert.NotNil(tt, err)
				assert.Equal(tt, c.expectedErrorMessage, err.Error(), "wrong error message")
			} else {
				assert.Nil(tt, err)

				assert.Equal(tt, c.expected, repository, "wrong repository")
			}
		})
	}
}
