package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/imroc/req"
	uuid "github.com/satori/go.uuid"
)

// File interface
type File interface {
	io.Reader
	io.Seeker
}

// SaveImage save image with a specific scale, depend on ffmpeg
func SaveImage(file File, scale string, path string) (string, error) {
	buffer := make([]byte, 512)
	file.Read(buffer)
	filetype := http.DetectContentType(buffer)
	var ext string
	switch filetype {
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		return "", errors.New("illegal image type")
	}

	// uuid
	id := uuid.Must(uuid.NewV4(), nil).String()

	// save origin file
	originName := "o-" + id + ext
	originPath := path + originName
	targetName := id + ext
	targetPath := path + targetName

	originFile, err := os.Create(originPath)
	if err != nil {
		return "", err
	}
	file.Seek(0, 0)
	if _, err := io.Copy(originFile, file); err != nil {
		return "", err
	}
	originFile.Close()
	defer os.Remove(originPath)

	// ffmpeg process
	cmd := exec.Command(
		"ffmpeg",
		"-i", originPath,
		"-y", "-strict", "-2",
		"-vf", "\"scale="+scale+":force_original_aspect_ratio=decrease\"",
		targetPath,
	)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return targetName, nil
}

// DownloadImage download image with a specific scale, depend on ffmpeg
func DownloadImage(url string, scale string, path string) (string, error) {
	uuidStr := uuid.Must(uuid.NewV4(), nil).String()
	tmpFile := path + uuidStr

	ri, err := req.Get(url)
	if err != nil {
		return "", err
	}
	if ri.Response().StatusCode != 200 {
		return "", errors.New("image download error")
	}
	defer os.Remove(tmpFile)
	err = ri.ToFile(tmpFile)
	if err != nil {
		return "", err
	}

	file, err := os.Open(tmpFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	targetName, err := SaveImage(file, scale, path)
	if err != nil {
		return "", err
	}

	return targetName, nil
}

// OptimizeImage make image limited in a specific scaling
func OptimizeImage(path string, scale string) (string, error) {
	dir := filepath.Dir(path)
	target := dir + "/o-" + filepath.Base(path)

	// target already exists, return
	if _, err := os.Stat(target); err == nil {
		return target, nil
	}

	// ffmpeg process
	cmd := exec.Command(
		"ffmpeg",
		"-i", path,
		"-y", "-strict", "-2",
		"-vf", "\"scale="+scale+":force_original_aspect_ratio=decrease\"",
		target,
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return target, nil
}

// SaveVideo save video with a specific scale, depend on ffmpeg
func SaveVideo(file File, scale string, path string) (string, string, error) {
	buffer := make([]byte, 512)
	file.Read(buffer)
	filetype := http.DetectContentType(buffer)
	var ext string
	switch filetype {
	case "video/x-flv":
		ext = ".flv"
	case "video/mp4":
		ext = ".mp4"
	case "video/3gpp":
		ext = ".3gp"
	case "video/quicktime":
		ext = ".mov"
	case "video/x-msvideo":
		ext = ".avi"
	case "video/x-ms-wmv":
		ext = ".wmv"
	default:
		return "", "", errors.New("illegal video type")
	}

	// uuid
	id := uuid.Must(uuid.NewV4(), nil).String()

	// save origin file
	originName := "o-" + id + ext
	originPath := path + originName
	targetName := id + ".mp4"
	targetPath := path + targetName
	posterName := path + ".jpg"
	posterPath := path + posterName

	originFile, err := os.Create(originPath)
	if err != nil {
		return "", "", err
	}
	file.Seek(0, 0)
	if _, err := io.Copy(originFile, file); err != nil {
		return "", "", err
	}
	originFile.Close()
	defer os.Remove(originPath)

	// ffmpeg process
	cmd := exec.Command(
		"ffmpeg",
		"-i", originPath,
		"-y", "-strict", "-2",
		"-ss", "00:00:00", "-t", "10",
		"-vf", "\"scale="+scale+":force_original_aspect_ratio=decrease\"",
		targetPath,
	)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", "", err
	}

	// poster
	cmd = exec.Command(
		"ffmpeg",
		"-i", targetPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-vf", "\"scale="+scale+":force_original_aspect_ratio=decrease\"",
		posterPath,
	)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", "", err
	}

	return targetName, posterName, nil
}
