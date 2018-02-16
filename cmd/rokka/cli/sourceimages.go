package cli

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

var sourceImagesListOptions rokka.ListSourceImagesOptions
var dynamicMetadataOptions rokka.DynamicMetadataOptions
var userMetadataName string

func listSourceImages(c *rokka.Client, args []string) (interface{}, error) {
	return c.ListSourceImages(args[0], sourceImagesListOptions)
}

func getSourceImage(c *rokka.Client, args []string) (interface{}, error) {
	return c.GetSourceImage(args[0], args[1])
}

func createSourceImage(c *rokka.Client, args []string) (interface{}, error) {
	if _, err := os.Stat(args[1]); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(args[1])
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return c.CreateSourceImage(args[0], path.Base(args[1]), file)
}

func addDynamicMetadata(c *rokka.Client, args []string) (interface{}, error) {
	b := bytes.NewBufferString(args[3])
	return c.AddDynamicMetadata(args[0], args[1], args[2], b, dynamicMetadataOptions)
}

func deleteDynamicMetadata(c *rokka.Client, args []string) (interface{}, error) {
	return c.DeleteDynamicMetadata(args[0], args[1], args[2], dynamicMetadataOptions)
}

func updateUserMetadata(c *rokka.Client, args []string) (interface{}, error) {
	org := args[0]
	hash := args[1]
	b := bytes.NewBufferString(args[2])

	if userMetadataName != "" {
		if err := c.UpdateUserMetadataByName(org, hash, userMetadataName, b); err != nil {
			return nil, err
		}
	} else {
		if err := c.UpdateUserMetadata(org, hash, b); err != nil {
			return nil, err
		}
	}
	return c.GetSourceImage(org, hash)
}

func deleteUserMetadata(c *rokka.Client, args []string) (interface{}, error) {
	org := args[0]
	hash := args[1]

	if userMetadataName != "" {
		if err := c.DeleteUserMetadataByName(org, hash, userMetadataName); err != nil {
			return nil, err
		}
	} else {
		if err := c.DeleteUserMetadata(org, hash); err != nil {
			return nil, err
		}
	}
	return c.GetSourceImage(org, hash)
}

// sourceImagesCmd represents the sourceImages command
var sourceImagesCmd = &cobra.Command{
	Use:                   "sourceimages",
	Short:                 "Create, list, search and show source images",
	Run:                   nil,
	Aliases:               []string{"si"},
	DisableFlagsInUseLine: true,
}

var sourceImagesListCmd = &cobra.Command{
	Use:                   "list [org]",
	Short:                 "List/Search source images",
	Args:                  cobra.ExactArgs(1),
	Aliases:               []string{"l"},
	DisableFlagsInUseLine: true,
	Run: run(listSourceImages, "Name\tHash\tDetails\n{{range .Items}}{{.Name}}\t{{.Hash}}\t{{.MimeType}}, {{.Width}}x{{.Height}}\n{{end}}\nTotal: {{.Total}}\n"),
}

const sourceImageTemplate = "Hash:\t{{.Hash}} ({{.ShortHash}})\nName:\t{{.Name}}\nDetails:\t{{.MimeType}}, {{.Width}}x{{.Height}}, {{.Size}}Bytes\nCreated at:\t{{datetime .Created}}\nBinary hash:\t{{.BinaryHash}}{{if .UserMetadata}}\nUser metadata:{{range $key, $value := .UserMetadata}}\n  {{$key}}:\t{{$value}}{{end}}{{end}}{{if .DynamicMetadata}}\nDynamic metadata:{{range $key, $value := .DynamicMetadata}}\n  {{$key}}:\t{{$value}}{{end}}{{end}}\n"

var sourceImagesGetCmd = &cobra.Command{
	Use:                   "get [org] [hash]",
	Short:                 "Get details of a source image by hash",
	Args:                  cobra.ExactArgs(2),
	Aliases:               []string{"g"},
	DisableFlagsInUseLine: true,
	Run: run(getSourceImage, sourceImageTemplate),
}

var sourceImagesCreateCmd = &cobra.Command{
	Use:                   "create [org] [file]",
	Short:                 "Upload a new image",
	Args:                  cobra.ExactArgs(2),
	Aliases:               []string{"c"},
	DisableFlagsInUseLine: true,
	Run: run(createSourceImage, fmt.Sprintf("{{range .Items}}%s{{end}}", sourceImageTemplate)),
}

var sourceImagesDynamicMetadataCmd = &cobra.Command{
	Use:                   "dynamic-metadata",
	Short:                 "Add or remove dynamic metadata on a source image",
	Run:                   nil,
	Aliases:               []string{"dm"},
	DisableFlagsInUseLine: true,
}

var sourceImagesAddDynamicMetadataCmd = &cobra.Command{
	Use:   "add [org] [hash] [name] [json]",
	Short: "Add dynamic metadata",
	Long: `Adding dynamic metadata generates a new image and returns the location of the new image.
If the deletePrevious flag is supplied, the previous image will be deleted.`,
	Args:                  cobra.ExactArgs(4),
	Aliases:               []string{"a"},
	DisableFlagsInUseLine: true,
	Run: run(addDynamicMetadata, "Location: {{.Location}}"),
}

var sourceImagesDeleteDynamicMetadataCmd = &cobra.Command{
	Use:   "delete [org] [hash] [name]",
	Short: "Delete dynamic metadata",
	Long: `Deleting dynamic metadata generates a new image and returns the location of the new image.
If the deletePrevious flag is supplied, the previous image will be deleted.`,
	Args:                  cobra.ExactArgs(3),
	Aliases:               []string{"d"},
	DisableFlagsInUseLine: true,
	Run: run(deleteDynamicMetadata, "Location: {{.Location}}"),
}

var sourceImagesUserMetadataCmd = &cobra.Command{
	Use:                   "user-metadata",
	Short:                 "Update or remove user metadata on a source image",
	Run:                   nil,
	Aliases:               []string{"dm"},
	DisableFlagsInUseLine: true,
}

var sourceImagesUpdateUserMetadataCmd = &cobra.Command{
	Use:   "update [org] [hash] [json]",
	Short: "Update user metadata",
	Long: `Update patches the currently set user metadata on the sourceimage, updating existing field names or adding new ones.
In case the --name flag is specified, the value will be set on that specific field name.`,
	Args:                  cobra.ExactArgs(3),
	Aliases:               []string{"u"},
	DisableFlagsInUseLine: true,
	Run: run(updateUserMetadata, sourceImageTemplate),
}

var sourceImagesDeleteUserMetadataCmd = &cobra.Command{
	Use:   "delete [org] [hash]",
	Short: "Delete user metadata",
	Long: `Delete removes user metadata on the sourceimage.
In case the --name flag is specified, only that field is removed.`,
	Args:                  cobra.ExactArgs(2),
	Aliases:               []string{"d"},
	DisableFlagsInUseLine: true,
	Run: run(deleteUserMetadata, sourceImageTemplate),
}

func init() {
	rootCmd.AddCommand(sourceImagesCmd)

	sourceImagesCmd.AddCommand(sourceImagesListCmd)
	sourceImagesCmd.AddCommand(sourceImagesGetCmd)
	sourceImagesCmd.AddCommand(sourceImagesCreateCmd)

	sourceImagesCmd.AddCommand(sourceImagesDynamicMetadataCmd)
	sourceImagesDynamicMetadataCmd.AddCommand(sourceImagesAddDynamicMetadataCmd)
	sourceImagesDynamicMetadataCmd.AddCommand(sourceImagesDeleteDynamicMetadataCmd)

	sourceImagesCmd.AddCommand(sourceImagesUserMetadataCmd)
	sourceImagesUserMetadataCmd.AddCommand(sourceImagesUpdateUserMetadataCmd)
	sourceImagesUserMetadataCmd.AddCommand(sourceImagesDeleteUserMetadataCmd)

	sourceImagesListCmd.Flags().IntVarP(&sourceImagesListOptions.Limit, "limit", "l", 20, "Limit")
	sourceImagesListCmd.Flags().IntVarP(&sourceImagesListOptions.Offset, "offset", "o", 0, "Offset")
	sourceImagesListCmd.Flags().StringVar(&sourceImagesListOptions.Hash, "hash", "", "Hash")
	sourceImagesListCmd.Flags().StringVar(&sourceImagesListOptions.BinaryHash, "binaryHash", "", "Binary hash")
	sourceImagesListCmd.Flags().StringVar(&sourceImagesListOptions.Size, "size", "", "Size in kilobytes")
	sourceImagesListCmd.Flags().StringVar(&sourceImagesListOptions.Format, "format", "", "Format")
	sourceImagesListCmd.Flags().StringVar(&sourceImagesListOptions.Width, "width", "", "Width")
	sourceImagesListCmd.Flags().StringVar(&sourceImagesListOptions.Height, "height", "", "Height")
	sourceImagesListCmd.Flags().StringVar(&sourceImagesListOptions.Created, "created", "", "Created")
	sourceImagesListCmd.Flags().StringVar(&sourceImagesListOptions.Sort, "sort", "", "Sort")

	sourceImagesAddDynamicMetadataCmd.Flags().BoolVar(&dynamicMetadataOptions.DeletePrevious, "deletePrevious", false, "Delete previous image")
	sourceImagesDeleteDynamicMetadataCmd.Flags().BoolVar(&dynamicMetadataOptions.DeletePrevious, "deletePrevious", false, "Delete previous image")

	sourceImagesUpdateUserMetadataCmd.Flags().StringVar(&userMetadataName, "name", "", "Update only the specified field")
	sourceImagesDeleteUserMetadataCmd.Flags().StringVar(&userMetadataName, "name", "", "Delete only the specified field")
}
