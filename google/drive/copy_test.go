package drive

import (
	"testing"

	"github.com/cebas/go-util/util"
)

var (
	log    = util.NewLog()
	dryRun = false // Set to true for dry run, false for actual copy
)

func TestCopyRecursively(t *testing.T) {

	sourceFolderId := "1a2b3c4d5e6f7g8h9i0j" // Replace with your source folder ID
	destFolderId := "0j9i8h7g6f5e4d3c2b1a"   // Replace with your destination folder ID

	var err error

	driveWs, err := NewDrive("ws.json", true, &log)
	util.FatalErrorCheck(err)

	driveGm, err := NewDrive("gm.json", false, &log)
	util.FatalErrorCheck(err)

	sourceFolder, err := driveWs.ExistingFolderById(sourceFolderId)
	util.FatalErrorCheck(err)

	destFolder, err := driveGm.ExistingFolderById(destFolderId)
	util.FatalErrorCheck(err)

	//goland:noinspection GoBoolExpressions
	err = CopyRecursively(
		&driveWs, &driveGm,
		&sourceFolder, &destFolder,
		false, dryRun, "",
	)
	util.FatalErrorCheck(err)
}
