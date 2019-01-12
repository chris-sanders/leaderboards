The primary purpose of this project is for me to learn Go.

This project is a tool for backing up exporting and merge Local Leaderboard
vales for the Occulus game Beatsaber. The intention is to allow backing up the
Local scores to a database, sending that database to a server via sftp,
downloading other users databases, and then selectively merging the data into
your local data.

A sftp server to perform the sync with is required and not part of this project.

On first run a default settings.yaml will be written out. With the following
settings which must be edited.
Global.account: A name that uniquely identifies your file on the sftp server
Global.local_folder: The folder with game data, by default expands to the games
default folder location.
Global.logging: Not used currently
Import.local_limit: Maximum number of scores for each local player to add to your game
file.
Import.remote_limit: Maximum number of scores for each remote player to add to
your game file.
Sync.accounts: Which of the accounts to download from the sever, 'all' for all
files in the sync folder
Sync.folder: Path on the server where the files are sent and retrieved
Sync.password: Sftp password
Sync.port: Sftp port
Sync.server: Sftp host
Sync.user: Sftp user name

