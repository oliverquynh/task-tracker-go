# Task Tracker

This is a solution written in Go for the task tracker project at https://roadmap.sh/projects/task-tracker

## Installation

Clone the this repository and go into its directory

```bash
cd /path/to/workspace

git clone https://github.com/oliverquynh/task-tracker-go.git

cd task-tracker-go
```

Build the task-tracker cli program

```bash
go build task-tracker.go
```

## Usage

Use can try these commands in sequence to figure out how this CLI works.

```bash
# 1. List existing tasks.
./task-tracker list

# 2. Add some tasks and check if they are saved.
./task-tracker add "Study Go basic syntax"
./task-tracker add "Build a CLI program in Go"
./task-tracker add "Build a web-based in Go"
./task-tracker add "Share my journey to GitHub"
./task-tracker list

# 3. Edit some tasks and check if they are saved.

./task-tracker edit 1 "Study Go including JSON encoding/decoding"
./task-tracker list
./task-tracker mark 1 "in-progress"
./task-tracker list
./task-tracker mark 1 "done"
./task-tracker list

# 4. Delete a task and check if it is deleted.
./task-tracker delete 3
./task-tracker list
```
