package updater

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/cmd/updater/logging"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
)

type FileData struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type FolderStructure struct {
	Name    string             `json:"name,omitempty"`
	Parent  *FolderStructure   `json:"-"`
	Files   []*FileData        `json:"files,omitempty"`
	Folders []*FolderStructure `json:"folders,omitempty"`
}

const (
	outFileName = "versions.json"
)

var (
	inDir      string
	outDir     string
	shouldHash bool

	generateCmd = &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen", "g"},
		Short:   "Generates metadata files for project releases or patch given a project name",
		Long: `Generates metadata files for project releases or patch given a project name.

	The project name is case sensitive.
	The in directory should be the base directory for a packaged release or patch.`,
		Example: "updater g SRO",
		Run: func(cmd *cobra.Command, args []string) {
			GenerateMetadataFile(inDir, outDir, shouldHash)
		},
	}
)

func GenerateMetadataFile(inputDir string, outputDir string, hashing bool) {
	root := &FolderStructure{
		Name:    inputDir,
		Parent:  nil,
		Files:   []*FileData{},
		Folders: []*FolderStructure{},
	}

	queue := []*FolderStructure{root}
	var current *FolderStructure
	hasher := md5.New()
	var f *os.File
	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]
		files, err := ioutil.ReadDir(current.FullPath())
		logging.HandleError(err)
		for _, file := range files {
			if file.IsDir() {
				newFS := &FolderStructure{
					Name:    file.Name(),
					Parent:  current,
					Files:   []*FileData{},
					Folders: []*FolderStructure{},
				}
				current.Folders = append(current.Folders, newFS)
				queue = append(queue, newFS)

			} else if file.Name() != outFileName {
				fmt.Printf("Processing file %s%s\n", current.FullPath(), file.Name())
				data := &FileData{
					Name: file.Name(),
					Hash: "",
				}

				if hashing {
					f, err = os.Open(current.FullPath() + "/" + file.Name())
					logging.HandleError(err)
					if _, err := io.Copy(hasher, f); err != nil {
						logging.HandleError(err)
					}
					data.Hash = hex.EncodeToString(hasher.Sum(nil))
				}

				current.Files = append(current.Files, data)
				if f != nil {
					err := f.Close()
					logging.HandleError(err)
				}
			}
		}
	}
	// Remove the root name because it should always be relative and therefore irrelevant
	root.Name = ""
	jsonBytes, err := json.Marshal(&root)
	logging.HandleError(err)
	outFile, err := os.Create(outputDir + "/" + outFileName)
	logging.HandleError(err)
	_, err = outFile.Write(jsonBytes)
	logging.HandleError(err)
	fmt.Printf("Created file %s\n", outFile.Name())
}

func initGenerate() {
	generateCmd.Flags().StringVarP(&inDir, "in", "i", ".", "release or patch base directory")
	generateCmd.Flags().StringVarP(&outDir, "out", "o", ".", "metadata file output directory")
	generateCmd.Flags().BoolVarP(&shouldHash, "hash", "", true, "should a hash be generated for each file")
	rootCmd.AddCommand(generateCmd)
}

func (fs *FolderStructure) FullPath() string {
	path := ""
	for current := fs; current != nil; current = current.Parent {
		path = current.Name + "/" + path
	}
	return path
}
