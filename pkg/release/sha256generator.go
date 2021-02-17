package release

import (
	"archive/zip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lab5e/reto/pkg/toolbox"
)

func checksumFileName(name, version string) string {
	return fmt.Sprintf("%s/%s/sha256sum_%s_%s.txt", archiveDir, version, name, version)
}

func generateChecksumFile(ctx *Context, files []string) error {
	checksumFilename := checksumFileName(ctx.Config.Name, ctx.Version)
	f, err := os.Create(checksumFilename)
	if err != nil {
		fmt.Printf("Could not create the checksum file %s: %v\n", checksumFilename, err)
		return err
	}
	defer f.Close()

	for _, v := range files {
		buf, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Printf("Unable to read %s: %v\n", v, err)
			return err
		}
		sum := sha256.Sum256(buf)

		line := fmt.Sprintf("%x  %s\n", sum, filepath.Base(v))
		fmt.Print(line)
		if _, err := f.Write([]byte(line)); err != nil {
			fmt.Printf("Could not write checksum to %s: %v\n", checksumFilename, err)
			return err
		}
	}
	return nil
}

// VerifyChecksums verifies that the files in the archive matches the checksums
// in the checksum file.
func VerifyChecksums(checksumFile, archive string) error {
	type checksum struct {
		File, Checksum string
	}
	var checksums []checksum
	buf, err := ioutil.ReadFile(checksumFile)
	if err != nil {
		fmt.Printf("Couldn't read %s: %v\n", checksumFile, err)
		return err
	}
	for _, line := range strings.Split(string(buf), "\n") {
		fields := strings.Split(line, "  ")
		if len(fields) == 2 {
			checksums = append(checksums, checksum{File: fields[1], Checksum: fields[0]})
		}
	}
	if len(checksums) == 0 {
		fmt.Printf("Could not find any checksums in file %s\n", checksumFile)
		return errors.New("no checksums")
	}

	zipArchive, err := zip.OpenReader(archive)
	if err != nil {
		fmt.Printf("Could not open archive %s: %v\n", archive, err)
		return err
	}
	defer zipArchive.Close()

	errs := 0
	for _, archivedFile := range zipArchive.File {
		found := false
		for _, csum := range checksums {
			if csum.File == archivedFile.Name {
				r, err := archivedFile.Open()
				if err != nil {
					fmt.Printf("Could not open archived file %s: %v\n", archivedFile.Name, err)
					return err
				}
				buf, err := ioutil.ReadAll(r)
				r.Close()
				if err != nil {
					fmt.Printf("Couldn't read archived file %s: %v\n", archivedFile.Name, err)
					return err
				}
				cs := fmt.Sprintf("%x", sha256.Sum256(buf))
				if cs == csum.Checksum {
					fmt.Printf("%s is OK\n", csum.File)
				} else {
					fmt.Printf("%sWARNING%s the checksum for %s does not match the checksum file\n", toolbox.Yellow, toolbox.Reset, archivedFile.Name)
					errs++
				}
				found = true
			}
		}
		if !found {
			errs++
			fmt.Printf("%sWARNING!%s %s is not in the signature file!\n", toolbox.Yellow, toolbox.Reset, archivedFile.Name)
		}
	}
	if errs > 0 {
		fmt.Printf("================================================\n")
		fmt.Printf("!!!! Archive has files with %schecksum errors%s !!!!\n", toolbox.Red, toolbox.Reset)
		fmt.Printf("================================================\n")
		return errors.New("checksum error")
	}
	return nil

}
