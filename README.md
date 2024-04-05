# pclub_vcs_task
CLI backup tool

## How to use
-download the file backup_tool and put folder in environment variables
-simply type backup_tool to create a backup
-to configure things use flags like
  -sd PATH to put path of source directory to backup from
  -ext TXT/JSON to modify extension of logger file between txt and json
  -e INT to modify number of times to encrypt(base 64) the files
-the files will be backed up to the folder ~/Documents/bd
-use command backup_tool reset to delete backup and logs
-use flag  -s VERSION_NUMBER to retrieve the backup version to another folder ~/Documents/bdshare
