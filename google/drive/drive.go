package drive

import (
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"github.com/cebas/go-util/google/auth"
	"github.com/cebas/go-util/util"
)

type Drive struct {
	ctx         context.Context
	gapiService *drive.Service
	log         *util.Log
}

const mimeTypeFolder = "application/vnd.google-apps.folder"

func NewDrive(credentialsFile string, readonly bool, log *util.Log) (newDrive Drive, err error) {

	ctx := context.Background()
	newDrive.ctx = ctx

	googleAuth := auth.NewGauth(ctx, credentialsFile, log)

	// DriveScope =         "https://www.googleapis.com/auth/drive" -          View and manage all your Drive files.
	// DriveReadonlyScope = "https://www.googleapis.com/auth/drive.readonly" - View and download all your Drive files.
	scope := util.If[string](readonly, drive.DriveReadonlyScope, drive.DriveScope)
	httpClient, err := googleAuth.HttpClient(scope)
	if err != nil {
		return newDrive, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	newDrive.gapiService, err = drive.NewService(ctx, option.WithHTTPClient(httpClient), option.WithScopes(scope))
	if err != nil {
		return newDrive, fmt.Errorf("failed to create Google Drive service: %w", err)
	}

	newDrive.log = log

	return
}

// ExistingFolderById retrieves a folder by ID.
func (d *Drive) ExistingFolderById(folderId string) (drive.File, error) {
	file, err := d.gapiService.Files.Get(folderId).Context(d.ctx).Do()
	if err != nil {
		return drive.File{}, fmt.Errorf("unable to get folder: %w", err)
	}
	if file.MimeType != mimeTypeFolder {
		return drive.File{}, fmt.Errorf("id %s is not a folder", folderId)
	}

	return *file, nil
}

// ExistingFolderByName retrieves a folder by name in a specific folder
func (d *Drive) ExistingFolderByName(folderName, parentFolderId string) (drive.File, int, error) {
	query := fmt.Sprintf(
		"mimeType='%s' and name='%s' and '%s' in parents and trashed=false",
		mimeTypeFolder,
		folderName,
		parentFolderId,
	)
	files, err := d.gapiService.Files.List().Q(query).Fields("files(id, name, mimeType)").Context(d.ctx).Do()
	if err != nil {
		return drive.File{}, -1, fmt.Errorf("unable to list folders in [%s]: %w", parentFolderId, err)
	}
	if len(files.Files) == 0 {
		return drive.File{}, 0, fmt.Errorf("folder [%s] not found in parent [%s]", folderName, parentFolderId)
	}
	if len(files.Files) > 1 {
		return drive.File{}, 2, fmt.Errorf("multiple folders found with name [%s] in parent [%s]", folderName, parentFolderId)
	}

	folder := files.Files[0]

	return *folder, 0, nil
}

func (d *Drive) FileExists(fileName, folderId string) (bool, *drive.File, error) {
	query := fmt.Sprintf("name='%s' and '%s' in parents and trashed=false", fileName, folderId)
	files, err := d.gapiService.Files.List().Q(query).Fields("files(id)").Context(d.ctx).Do()

	if err != nil {
		gErr, ok := err.(*googleapi.Error)

		// if err is googleapi.Error and the error code is 404, it means the file does not exist
		if ok && gErr.Code == 404 {
			// File does not exist, return false
			return false, nil, nil
		} else if ok && gErr.Code == 403 {
			// If the error is a permission error, return an error
			return false, nil, fmt.Errorf("permission denied when checking file existence: %v", err)
		} else {
			// For any other error, log it and return an error
			return false, &drive.File{}, fmt.Errorf("failed to check file existence: %v", err)
		}
	}

	if len(files.Files) > 0 {
		return true, files.Files[0], nil
	}

	// If no files found, return false
	return false, nil, nil
}

// Delete deletes a file/folder
func (d *Drive) Delete(file *drive.File) error {
	err := d.gapiService.Files.Delete(file.Id).Context(d.ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to delete: %s [%w]", file.Name, err)
	}
	return nil
}

// queries:
// https://developers.google.com/workspace/drive/api/guides/search-files#examples
// https://developers.google.com/workspace/drive/api/guides/ref-search-terms#file-properties
//	mimeType = 'application/vnd.google-apps.folder'
//	'%s' in parents
//	trashed = false

func (d *Drive) FolderFiles(folder *drive.File) ([]*drive.File, error) {
	query := fmt.Sprintf("'%s' in parents and trashed = false", folder.Id)

	files, err := d.gapiService.Files.List().
		Q(query).
		Fields("files(id,name,description,mimeType,createdTime,modifiedTime,folderColorRgb,starred)").
		OrderBy("folder desc,name").
		Context(d.ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list files in folder [%s]: %v", folder.Id, err)
	}
	return files.Files, nil
}

// createFolderFromTemplate creates a new folder with f as a parent,
// using originalFolder as a template
func (d *Drive) createFolderFromTemplate(originalFolder *drive.File, parentFolder *drive.File) (drive.File, error) {
	folderTemplate := *originalFolder
	folderTemplate.Parents = []string{parentFolder.Id}
	folderTemplate.Id = "" // Clear the ID to create a new folder

	newFolder, err := d.gapiService.Files.Create(&folderTemplate).Context(d.ctx).Do()
	if err != nil {
		return drive.File{}, fmt.Errorf("failed to create folder [%s]: %v", originalFolder.Name, err)
	}

	return *newFolder, nil
}

// CreateFile creates a file in the folder.
//func (f *Folder) CreateFile(ctx context.Context, name, content string) (File, error) {
//	file := &drive.File{
//		Name:    name,
//		Parents: []string{f.id},
//	}
//	createdFile, err := f.srv.Files.Create(file).
//		Media(strings.NewReader(content)).
//		Context(ctx).
//		Do()
//	if err != nil {
//		return nil, fmt.Errorf("failed to create file: %w", err)
//	}
//	return &googleFile{
//		srv:  f.srv,
//		id:   createdFile.Id,
//		name: createdFile.Name,
//	}, nil
//}
// CreateFolder creates a subfolder.
//func (f *Folder) CreateFolder(ctx context.Context, name string) (Folder, error) {
//	file := &drive.File{
//		Name:     name,
//		MimeType: "application/vnd.google-apps.folder",
//		Parents:  []string{f.id},
//	}
//	createdFile, err := f.srv.Files.Create(file).Context(ctx).Do()
//	if err != nil {
//		return nil, fmt.Errorf("failed to create subfolder: %w", err)
//	}
//	return &Folder{
//		srv:  f.srv,
//		id:   createdFile.Id,
//		name: createdFile.Name,
//	}, nil
//}
