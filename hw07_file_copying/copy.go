package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameInOutFiles        = errors.New("from file and out file are the same")
	ErrNoInFile              = errors.New("from file isn't set")
	ErrNoOutFile             = errors.New("to file isn't set")
	ErrWriteToOutput         = errors.New("error writing to output")
	ErrReadFromInput         = errors.New("error reading from input")
)

func copyWithProgressBar(in io.Reader, out io.Writer, bytesToCopy int64) error {
	// start new bar
	bar := pb.Full.Start64(bytesToCopy)

	// create proxy reader
	barReader := bar.NewProxyReader(io.LimitReader(in, bytesToCopy))

	// copy from proxy reader
	wb, err := io.Copy(out, barReader)

	// finish bar
	bar.Finish()

	if (err != nil) && (!errors.Is(err, io.EOF)) {
		return err
	}

	if wb != bytesToCopy {
		return ErrReadFromInput
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	var inFile *os.File

	switch {
	case fromPath == "":
		return ErrNoInFile
	case toPath == "":
		return ErrNoOutFile
	case fromPath == toPath:
		return ErrSameInOutFiles
	}

	if fromPath == toPath {
		return ErrSameInOutFiles
	}

	inFile, err1 := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err1 != nil {
		return err1
	}
	defer inFile.Close()

	var inFileInfo os.FileInfo
	inFileInfo, err1 = inFile.Stat()

	if err1 != nil {
		return err1
	}

	if (inFileInfo.Mode() & os.ModeType) != 0 {
		return ErrUnsupportedFile
	}

	size := inFileInfo.Size()

	if size < 0 {
		return ErrUnsupportedFile
	}

	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	var outFile *os.File
	var err2 error
	outFile, err2 = os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0))
	if err2 != nil {
		switch {
		case os.IsPermission(err2):
			err3 := os.Chmod(toPath, os.FileMode(0o600))
			if err3 != nil {
				return err2
			}
			var err4 error
			outFile, err4 = os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0))
			if err4 != nil {
				return err4
			}
		default:
			return err2
		}
	}
	defer outFile.Close()
	defer outFile.Chmod(os.FileMode(0o444))

	bytesToCopy := size - offset

	if (limit > 0) && (bytesToCopy > limit) {
		bytesToCopy = limit
	}

	if bytesToCopy <= 0 {
		return nil
	}

	if offset > 0 {
		n, err := inFile.Seek(offset, 0)
		if (n != offset) || (err != nil) {
			return ErrReadFromInput
		}
	}
	err5 := copyWithProgressBar(inFile, outFile, bytesToCopy)
	if err5 != nil {
		defer os.Remove(toPath)
		return err5
	}
	return nil
}
