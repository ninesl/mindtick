#### About

Mindtick is a lightweight CLI tool designed to help you track your progress, tasks, and accomplishments in a structured, timestamped log. With Mindtick, you can easily log your wins, notes, fixes, and tasks to maintain a clear and organized timeline of your work.

All data is stored in a SQLite database file named `store.mindtick`, which is created in your current directory when you run `mindtick new`. This file acts as a portable and self-contained activity log, ensuring your thoughts and progress are always consolidated and easy to manage.

Each entry in the `store.mindtick` file is recorded with:
- A **timestamp** (automatically generated).
- A **message type** (`win`, `note`, `fix`, `task`).
- The **message content** you provide.

When you run a `mindtick` command, the tool searches for the `store.mindtick` file starting from your current directory and traversing up the directory tree. This means that you can use Mindtick from any subdirectory of your project, and the tool will always find and operate on the nearest `store.mindtick` file.

## Installation

Make sure you have the most recent version of go installed *(1.24.4 at the time of writing)*

You can find the most recent release of go here: https://go.dev/dl/

After go has been installed, run `go get github.com/ninesl/mindtick@latest`

If your `$PATH` variables are setup correctly *(will be by default)* run `mindtick` in a new terminal window to ensure Mindtick was installed correctly.

NOTE: ANSI color escape codes are used to color code tags and messages in Mindtick. Make sure you are using a terminal that supports
this type of text rendering (the default Windows Powershell and Commandprompt do not)

## Commands

### Usage: `mindtick <command>`

| Command   | Description                                         |
|-----------|-----------------------------------------------------|
| `help`    | Display this help message.                         |
| `new`     | Create a new `store.mindtick` file in the current directory. |
| `delete`  | Delete the `store.mindtick` file in the current directory. |
| `view`    | Display all messages in the current `store.mindtick` file. |
| `win`     | Add a win message: `mindtick win -<message>`.      |
| `note`    | Add a note message: `mindtick note -<message>`.    |
| `fix`     | Add a fix message: `mindtick fix -<message>`.      |
| `task`    | Add a task message: `mindtick task -<message>`.    |

## Examples

The output after running `mindtick view`

![example mindtick view](readme_assets/view.png)

Adding a message to your current `store.mindtick`

`mindtick win -your message goes here`

![example adding fix tag msg](readme_assets/addingmsg.png)

| Planned Features                              |                                                |
|--------------------------------------|------------------------------------------------------------|
| `search {tags}`                      | Display messages filtered by specific tags.                |
| `export {tags} {filetype}`           | Export messages to `.pdf`, `.csv`, or `.txt` based on tags. |
| `delete <id>`                        | Delete a specific message by its unique ID.                |
| `edit <id> <new message>`            | Edit an existing message by its ID.                        |
| `filter {tags}`                      | Filter messages in various commands by tags.               |
| `{win, note, fix, task}`             | Filter messages by type (win, note, fix, task).            |
| `{keyword}`                          | Search messages by a specific keyword or substring.        |
| `{today, yesterday, week, YYYY-MM-DD}` | Filter messages by date ranges.                           |

I also am planning on implementing user-created tags and the ability to turn off color codes (for terminals that can't render it)
