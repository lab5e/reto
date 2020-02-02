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

	"github.com/ExploratoryEngineering/reto/pkg/toolbox"
)

// GenerateSHA256File generates a text file with SHA256 checksums for the files
// in the release
func GenerateSHA256File(ctx *Context) error {
	fmt.Println("Generating SHA256 checksums...")

	return generateChecksumFile(ctx, nil)
}

func checksumFileName(name, version string) string {
	return fmt.Sprintf("%s/%s/sha256sum_%s_%s.txt", archiveDir, version, name, version)
}

func generateChecksumFile(ctx *Context, changelogBuf []byte) error {
	checksumFilename := checksumFileName(ctx.Config.Name, ctx.Version)
	f, err := os.Create(checksumFilename)
	if err != nil {
		toolbox.PrintError("Could not create the checksum file %s: %v", checksumFilename, err)
		return err
	}
	defer f.Close()

	// Generate checksum for changelogBuf
	if changelogBuf != nil {
		csum := fmt.Sprintf("%x  changelog.md\n", sha256.Sum256(changelogBuf))
		fmt.Print(csum)
		f.Write([]byte(csum))
	}
	for _, v := range ctx.Config.Files {
		buf, err := ioutil.ReadFile(v.Name)
		if err != nil {
			toolbox.PrintError("Unable to read %s: %v", v.Name, err)
			return err
		}
		sum := sha256.Sum256(buf)

		line := fmt.Sprintf("%x  %s\n", sum, filepath.Base(v.Name))
		fmt.Print(line)
		if _, err := f.Write([]byte(line)); err != nil {
			toolbox.PrintError("Could not write checksum to %s: %v", checksumFilename, err)
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
		toolbox.PrintError("Couldn't read %s: %v", checksumFile, err)
		return err
	}
	for _, line := range strings.Split(string(buf), "\n") {
		fields := strings.Split(line, "  ")
		if len(fields) == 2 {
			checksums = append(checksums, checksum{File: fields[1], Checksum: fields[0]})
		}
	}
	if len(checksums) == 0 {
		toolbox.PrintError("Could not find any checksums in file %s", checksumFile)
		return errors.New("no checksums")
	}

	zipArchive, err := zip.OpenReader(archive)
	if err != nil {
		toolbox.PrintError("Could not open archive %s: %v", archive, err)
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
					toolbox.PrintError("Could not open archived file %s: %v", archivedFile.Name, err)
					return err
				}
				buf, err := ioutil.ReadAll(r)
				r.Close()
				if err != nil {
					toolbox.PrintError("Couldn't read archived file %s: %v", archivedFile.Name, err)
					return err
				}
				cs := fmt.Sprintf("%x", sha256.Sum256(buf))
				if cs == csum.Checksum {
					fmt.Printf("%s is OK\n", csum.File)
				} else {
					toolbox.PrintError("**** WARNING the checksum for %s does not match the checksum file", archivedFile.Name)
					errs++
				}
				found = true
			}
		}
		if !found {
			errs++
			fmt.Printf("WARNING! %s is not in the signature file!\n", archivedFile.Name)
		}
	}
	if errs > 0 {
		toolbox.PrintError("================================================")
		toolbox.PrintError("!!!! Archive has files with checksum errors !!!!")
		toolbox.PrintError("================================================")
		return errors.New("checksum error")
	}
	return nil

}
