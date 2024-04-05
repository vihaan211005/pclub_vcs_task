### Overview

- The tool logs the date and time of the last backup and when the next backup is done only copy files that have been modified later to that date and time
- The tool has a disadvantage that files if deleted/renamed in source directory will cause issues as when backup is called it only checks file currently present is source directory
- The configurations like extension of logger file, number of times to encrypt, source directory can be changed using flags -ext -e -sd which are stored in file conf.json created where the tool is saved
- The log files store the time stamp of each backup version
- To check modification metadata is stored which can be accessed via fs.DirEntry.info().ModTime()
- When backed up the file.ext is saved as file.(v1).ext and when it is changed and backed up in a later version a new file file.(v3).ext is created. Version indexing starts from 0
- To transfer data of previous version the file with max version less than equal to back up version is transferred. Example if v2 is backep up in the case given above file.(v1).ext is transferred
