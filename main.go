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

	"github.com/jdrab/app-updater/config"
	"github.com/jdrab/app-updater/download"
	"github.com/jdrab/app-updater/platform"
)

var Version string = "0.2.0"
var runtimeApp string
var runtimeService string

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
		log.Println(err)
	}
	log.SetPrefix(time.Now().Format(time.RFC3339) + " ")
	log.SetFlags(0)
	log.SetOutput(logFile)
	log.Println("init")

}

/**
 * @var		mixed	config
 * @global
 */
//

var conf = config.Load()

func main() {

	if runtimeApp != "" {
		conf.App = runtimeApp
	}
	if runtimeService != "" {
		conf.Service = runtimeService
	}

	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	defaultInstallationDir, _ := filepath.Abs(".")

	// required flags
	downloadURLFlag := cli.String("url", "", "https://server/update-1.0.2.zip (Required)")
	// whatever server requires to identify exact platform and version of client app
	clientVersionFlag := cli.String("client", "", "1.0.1 (platform;arch;os_version) (Required)")
	updateChecksumFlag := cli.String("checksum", "", "package sha256 sum (Required)")

	// "optional" flags
	installDirFlag := cli.String("installdir", defaultInstallationDir, "directory where to unzip archive,default to directory where updater is located")
	serviceFlag := cli.String("service", conf.Service, "service name")
	//the app will be killed for now
	exeFlag := cli.String("app", conf.App, "Client app name to be >killed< before unpacking update")

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
		fmt.Printf("%s %s\n", os.Args[0], Version)
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
	// Check if it already exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Println("downloading update..")
		// Archive does not exist, let's download it
		_, err = download.FromURL(filename, *downloadURLFlag, *verboseFlag)

		if err != nil {
			log.Printf("fatal error: %v", err)
			os.Exit(2)
		}

	}

	// If it's already downloaded verify checksum
	_, err = download.VerifyChecksum(filename, *updateChecksumFlag)
	if err != nil {
		log.Printf("error: file exists but checksum is invalid. re-downloading file \nfrom\t%v \nto\t%v\n", *downloadURLFlag, filename)

		filename, err = download.FromURL(filename, *downloadURLFlag, *verboseFlag)
		if err != nil {
			log.Printf("fatal error: %v", err)
			os.Exit(2)
		}
		log.Printf("download successful")

		_, err := download.VerifyChecksum(filename, *updateChecksumFlag)
		if err != nil {
			log.Printf("error: checksum verification failed - %v", err)
			os.Exit(3)
		}
	}
	//End application and stop the service
	platform.KillProcessByName(*exeFlag)
	platform.StopService(*serviceFlag)

	// file is downloaded and checksum is valid at this point let's extract it
	archive, err := fastzip.NewExtractor(filename, *installDirFlag)
	if err != nil {
		log.Printf("fatal error: failed to create extractor - %v", err)
		os.Exit(4)
	}
	defer archive.Close()

	// Register faster decompressor
	archive.RegisterDecompressor(zip.Deflate, fastzip.FlateDecompressor())

	if *verboseFlag {
		log.Printf("extracting %v", filename)
		log.Printf("extracting to %v", *installDirFlag)
	}

	// Extract archive files
	err = archive.Extract(context.Background())
	if err != nil {
		log.Printf("fatal error: update extraction failed - %v", err)
		os.Exit(5)
	}

	// start service as last thing
	platform.StartService(conf.Service)

	// delete the update file if everything succeed
	log.Printf("deleting update file")
	err = os.Remove(filename)

	if err != nil {
		log.Printf("error: deleting update file %v failed - %v", filename, err)
	}

	log.Printf("done\n")

}
