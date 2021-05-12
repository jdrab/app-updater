package main

import (
	"archive/zip"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/saracen/fastzip"

	"github.com/jdrab/app-updater/download"
	"github.com/jdrab/app-updater/platform"
	"github.com/jdrab/app-updater/serviceconfig"
)

// UpdateResponse is a json response to update request by an updater
type UpdateResponse struct {
	URL      string `json:"url"`
	Checksum string `json:"sha256"`
}

// func updateResponseParse(body []byte) (*UpdateResponse, error) {
// 	var resp = new(UpdateResponse)
// 	err := json.Unmarshal([]byte(body), &resp)
// 	return resp, err
// }

func init() {
	logFile, err := os.OpenFile("updater.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetPrefix(time.Now().Format(time.RFC3339) + " ")
	log.SetFlags(0)
	log.SetOutput(logFile)
	log.Println("init..")

}

var config = serviceconfig.Load()

func main() {
	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	defaultInstallationDir, _ := filepath.Abs(".")

	// required flags
	downloadURLFlag := cli.String("url", "", "https://server/update-1.0.2.zip (Required)")
	// whatever server requires to identify exact platform and version of client app
	clientVersionFlag := cli.String("client", "", "1.0.1 (platform;arch;os_version) (Required)")
	updateChecksumFlag := cli.String("checksum", "", "package sha256 sum (Required)")

	// "optional" flags
	installDirFlag := cli.String("installdir", defaultInstallationDir, "directory where to unzip archive,default to directory where updater is located")
	serviceFlag := cli.String("service", "my-service", "service name")
	//the app will be killed for now
	exeFlag := cli.String("app", config.AppName, "Client app name to be >killed< before unpacking update")

	verboseFlag := cli.Bool("verbose", false, "")

	if *verboseFlag {
		download.Verbose = true
		platform.Verbose = true
	}

	versionCmd := cli.Bool("version", false, "display version")

	// must be called after flags definition bug before their usage
	cli.Parse(os.Args[1:])

	reqFlags := make(map[string]string)
	reqFlags["downloadURL"] = *downloadURLFlag
	reqFlags["clientVersion"] = *clientVersionFlag
	reqFlags["updateChecksum"] = *updateChecksumFlag

	optFlags := make(map[string]string)
	optFlags["installDirCmd"] = *installDirFlag
	optFlags["serviceCmd"] = *serviceFlag
	optFlags["exeCmd"] = *exeFlag

	if *versionCmd {
		fmt.Printf("%s version %s\n", os.Args[0], config.Version)
		os.Exit(0)
	}

	// find if any required flag is missing
	for k, value := range reqFlags {
		if value == "" {
			fmt.Println("Usage:\n", os.Args[0])
			cli.PrintDefaults()
			fmt.Println("\nerror: missing required flag", strings.TrimSuffix(k, "Cmd"), "please read usage")
			os.Exit(1)
		}
	}

	if *verboseFlag {
		log.Printf("running installer in %s\n", *installDirFlag)
	}

	// merge reqFlags and optflags
	definedFlags := optFlags
	for k, v := range reqFlags {
		definedFlags[k] = v
	}

	tokens := strings.Split(*downloadURLFlag, "/")
	var filename string
	// get the filename from url like/this/is/the_archive.zip
	filename = tokens[len(tokens)-1] // return -1 slice

	if *verboseFlag {
		log.Printf("filename %v", filename)
	}

	// Check if it already exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Println("update zip does not exist, must be downloaded first")
		// Archive does not exist, let's download it
		filename, err = download.FromURL(filename, *downloadURLFlag, *verboseFlag)
		if err != nil {
			log.Fatalf("FATAL: %v", err)
		}

	}

	// If it's already downloaded verify checksum
	_, err := download.VerifyChecksum(filename, *updateChecksumFlag)
	if err != nil {
		log.Printf("error: file exists but checksum is invalid. re-downloading file from %v => %v\n", *downloadURLFlag, filename)

		filename, err = download.FromURL(filename, *downloadURLFlag, *verboseFlag)
		if err != nil {
			log.Fatalf("FATAL: %v", err)
			os.Exit(1)
		}
		log.Printf("download successful")

		_, err := download.VerifyChecksum(filename, *updateChecksumFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
	//End application and stop the service
	platform.KillProcessByName(*exeFlag)
	platform.StopService(*serviceFlag)

	// file is downloaded and checksum is valid at this point let's extract it
	archive, err := fastzip.NewExtractor(filename, *installDirFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	// Register faster decompressor
	archive.RegisterDecompressor(zip.Deflate, fastzip.FlateDecompressor())

	if *verboseFlag {
		log.Printf("extracting %v", filename)
		log.Printf("extracting to %v", *installDirFlag)
	}
	// Extract archive files
	if err = archive.Extract(context.Background()); err != nil {
		log.Fatal(err)
	}

	// start service as last thing
	platform.StartService(config.ServiceName)

}
