package drive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

// CopyFile copies a file from sourceDrive to destDrive, placing it in destFolder.
func CopyFile(destDrive *Drive, sourceFile *drive.File, destFolder *drive.File, forceOverwrite bool, currentPath string, dryRun bool) (err error) {

	fileDescription := fileDescription(currentPath, sourceFile)

	// Check if file already exists in destination
	var exists bool
	var destFile *drive.File

	exists, destFile, err = destDrive.FileExists(sourceFile.Name, destFolder.Id)
	if err != nil {
		return err
	}

	if exists {
		if !forceOverwrite {
			destDrive.log.Printf(1, "S [%s]\n", fileDescription)
			return nil
		}

		destDrive.log.Printf(1, "D [%s]\n", fileDescription)

		if !dryRun {
			err = destDrive.Delete(destFile)
			if err != nil {
				return fmt.Errorf("failed to delete existing file %s: %w", destFile.Name, err)
			}
		}
	}

	copyFile := *sourceFile
	copyFile.Id = "" // Clear the ID to create a new file
	copyFile.Parents = []string{destFolder.Id}

	if dryRun {
		err = nil
	} else {
		_, err = destDrive.gapiService.Files.Copy(sourceFile.Id, &copyFile).Context(destDrive.ctx).Do()
	}
	if err == nil {
		destDrive.log.Printf(1, "C [%s]\n", fileDescription)
		return
	}

	return fmt.Errorf("failed to copy file [%s]:  %v", fileDescription, err)
}

// CopyRecursively copies content of the sourceFolder folder's to a specified existing folder in the destination Drive,
// skipping existing items unless forceOverwrite is true.
func CopyRecursively(sourceDrive *Drive, destDrive *Drive, sourceFolder *drive.File, destFolder *drive.File, forceOverwrite bool, dryRun bool, currentPath string) error {

	sourceDrive.log.Printf(1, "%s [%s]\n", fileTypeSymbol(sourceFolder), currentPath)

	sourceFiles, err := sourceDrive.FolderFiles(sourceFolder)
	if err != nil {
		return err
	}

	var destSubFolder drive.File
	var nfiles int

	for _, sourceFile := range sourceFiles {
		if isFolder(sourceFile) {
			// search folder in destination by name
			destSubFolder, nfiles, err = destDrive.ExistingFolderByName(sourceFile.Name, destFolder.Id)
			if err != nil {
				// If the subfolder does not exist, create it
				if nfiles == 0 {
					if dryRun {
						destSubFolder, err = drive.File{Id: "dummy id", Name: sourceFile.Name}, nil
					} else {
						destSubFolder, err = destDrive.createFolderFromTemplate(sourceFile, destFolder)
						if err != nil {
							return err
						}
					}
				} else {
					// search error or exists but not unique
					return err
				}
			}

			// exists and is unique: use it
			err = CopyRecursively(sourceDrive, destDrive, sourceFile, &destSubFolder, forceOverwrite, dryRun, currentPath+"/"+destSubFolder.Name)
			if err != nil {
				return fmt.Errorf("failed to copy subfolder [%s]: %w", sourceFile.Name, err)
			}
		} else {
			err = CopyFile(destDrive, sourceFile, destFolder, forceOverwrite, currentPath, dryRun)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
