package utils

import (
	"errors"
	"strconv"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"context"

	"github.com/imroc/req"
	uuid "github.com/satori/go.uuid"
	"github.com/gabriel-vasile/mimetype"
)

// File interface
type File interface {
	io.Reader
	io.Seeker
}

// SaveImage save image with a specific scale, depend on ffmpeg
func SaveImage(ctx context.Context, file File, scale string, path string) (string, error) {
	file.Seek(0, 0)
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return "", err
	}

	var ext string
	switch mime.String() {
	case "image/jpeg", "image/jpg", "image/png":
		ext = mime.Extension()
	default:
		return "", errors.New("image must be jpg or png")
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
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-i", originPath,
		"-y", "-strict", "-2",
		"-vf", "scale="+scale+":force_original_aspect_ratio=decrease",
		targetPath,
	)
	log.Info().Str("cmd", cmd.String()).Send()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(output))
	}

	return targetName, nil
}

// DownloadImage download image with a specific scale, depend on ffmpeg
func DownloadImage(ctx context.Context, url string, scale string, path string) (string, error) {
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

	targetName, err := SaveImage(ctx, file, scale, path)
	if err != nil {
		return "", err
	}

	return targetName, nil
}

// OptimizeImage make image limited in a specific scaling
func OptimizeImage(ctx context.Context, path string, scale string) (string, error) {
	dir := filepath.Dir(path)
	target := dir + "/o-" + filepath.Base(path)

	// target already exists, return
	if _, err := os.Stat(target); err == nil {
		return target, nil
	}

	// ffmpeg process
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-i", path,
		"-y", "-strict", "-2",
		"-vf", "scale="+scale+":force_original_aspect_ratio=decrease",
		target,
	)
	log.Info().Str("cmd", cmd.String()).Send()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(output))
	}

	return target, nil
}

// SaveVideo save video with a specific scale, depend on ffmpeg
func SaveVideo(ctx context.Context, file File, scale string, time int, path string) (string, string, error) {
	file.Seek(0, 0)
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return "", "", err
	}

	var ext string
	switch mime.String() {
	case "video/mpeg", "video/quicktime", "video/mp4", "video/webm", "video/x-msvideo", "video/x-flv", "video/x-matroska":
		ext = mime.Extension()
	default:
		return "", "", errors.New("video must be mpeg, mov, mp4, webm, avi, flv or mkv")
	}

	// uuid
	id := uuid.Must(uuid.NewV4(), nil).String()

	// save origin file
	originName := "o-" + id + ext
	originPath := path + originName
	targetName := id + ".mp4"
	targetPath := path + targetName
	posterName := id + ".jpg"
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
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-i", originPath,
		"-y", "-strict", "-2",
		"-ss", "00:00:00", "-t", strconv.Itoa(time),
		"-vf", "scale="+scale+":force_original_aspect_ratio=decrease",
		targetPath,
	)
	log.Info().Str("cmd", cmd.String()).Send()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", errors.New(string(output))
	}

	// poster
	cmd = exec.CommandContext(
		ctx,
		"ffmpeg",
		"-i", targetPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-vf", "scale="+scale+":force_original_aspect_ratio=decrease",
		posterPath,
	)
	log.Info().Str("cmd", cmd.String()).Send()
	output, err = cmd.CombinedOutput()
	if err != nil {
		return "", "", errors.New(string(output))
	}

	return targetName, posterName, nil
}
