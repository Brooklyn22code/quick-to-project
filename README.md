Do you think everytime you need to find your project by finding it in the file explorer or open folder in vscode is very slow and looks stupid. 
This is a very easy CLI tool help you doing that without touching mouse.

To run this CLI, make sure you download golang compiler, vscode and SQLite.

You first need clone the repo to your computer, and type in the terminal of your repo direction:
go mod tidy

Then type in terminal:
sqlite3 project.db

Finally you put the direction where you clone the repo in the environment variables.

You are done right now, open terminal and type command quickpro, then you should see:
command>

There are 5 commands in this CLI tools:

create: create a project, it will ask for name of project and path of project,

list: it will gove you every project's name, path and id,

open: it will ask project id, then it will open vscode for this project,

delete: it will ask project id, and delete the project,

exit: it will exit the CLI

