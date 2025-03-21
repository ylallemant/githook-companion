package environment

import (
	"fmt"
	"os"
	"runtime"
)

func Debug() string {
	dirname, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get local directory: %s", err.Error()))
	}

	oldpwd := FindEnvVar("OLDPWD")
	// absolute path of the called executable
	underscore := FindEnvVar("_")

	calledFromTerminal, _ := CalledFromTerminal()

	pwdEqual := dirname == FindEnvVar("PWD")

	return fmt.Sprintf(
		`
++++ CENVIRONMENT +++++++++++++++++
  OS:   "%s"
  ARCH: "%s"

  os.Getwd(): "%s"
  PWD equal:   %v
  OLDPWD:     "%s"
  _:          "%s"

  TERM_PROGRAM      "%s"
  PACKAGE_PATH:     "%s"
  XPC_SERVICE_NAME: "%s"

  pathValid          => %v
  binaryPathValid    => %v
  CalledFromTerminal => %v
  `,
		runtime.GOOS,
		runtime.GOARCH,
		dirname,
		pwdEqual,
		oldpwd,
		underscore,
		FindEnvVar("TERM_PROGRAM"),
		FindEnvVar("PACKAGE_PATH"),
		FindEnvVar("XPC_SERVICE_NAME"),
		pathValid(dirname, oldpwd),
		binaryPathValid(underscore),
		calledFromTerminal,
	)
}
