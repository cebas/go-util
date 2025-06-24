package drive

import (
	"fmt"
	"strings"

	"google.golang.org/api/drive/v3"
)

func fileDescription(path string, file *drive.File) string {
	return fmt.Sprintf("%s %s/%s", fileTypeSymbol(file), path, file.Name)
}

func isFolder(file *drive.File) bool {
	return file.MimeType == mimeTypeFolder
}

// FileTypeSymbol returns a Unicode symbol for the file's MIME type or extension.
func fileTypeSymbol(file *drive.File) string {
	switch file.MimeType {
	case "application/vnd.google-apps.folder":
		return "📁"
	case "application/vnd.google-apps.document":
		return "📘"
	case "application/vnd.google-apps.spreadsheet":
		return "📊"

	case "text/plain":
		return "📝"
	case "application/pdf":
		return "📕"
	case "image/jpeg", "image/png", "image/gif":
		return "🖼️"
	case "video/mp4", "video/x-msvideo":
		return "🎥"
	case "audio/mpeg", "audio/wav":
		return "🎵"
	case "application/zip":
		return "📦"
	case "application/vnd.jgraph.mxfile": // draw.io file
		return "🗺️"

	default:
		// Fallback to extension-based mapping
		parts := strings.Split(file.Name, ".")
		if len(parts) > 1 {
			switch strings.ToLower(parts[len(parts)-1]) {
			case "go":
				return "🐹"
			case "json":
				return "🗃️"
			case "md":
				return "🗒️"
			case "csv":
				return "📋"
			}
		}
		return "📄" // Generic file
	}
}
