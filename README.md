# app-updater
Experimental app updater in golang


**Exit codes**
* 0 - success
* 1 - usage shown - missing required flags
* 2 - download failed
* 3 - checksum failed
* 4 - extractor error
* 5 - unzip failed
 
Errors running or stopping service are not fatal