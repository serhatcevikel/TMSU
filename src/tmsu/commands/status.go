/*
Copyright 2011 Paul Ruane.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
package main
*/

package main

import (
    "path/filepath"
    "fmt"
    "os"
)

type StatusCommand struct {}

func (this StatusCommand) Name() string {
    return "status"
}

func (this StatusCommand) Summary() string {
    return "lists file status"
}

func (this StatusCommand) Help() string {
    return `tmsu status
tmsu status FILE...

Shows the status of files.`
}

func (this StatusCommand) Exec(args []string) error {
    tagged := make([]string, 0, 10)
    untagged := make([]string, 0, 10)
    var error error

    if len(args) == 0 {
        tagged, untagged, error = this.status([]string{"."}, tagged, untagged)
    } else {
        tagged, untagged, error = this.status(args, tagged, untagged)
    }

    if error != nil {
        return error
    }

    for _, path := range tagged {
        fmt.Printf("T %v\n", path)
    }

    for _, path := range untagged {
        fmt.Printf("U %v\n", path)
    }

    return nil
}

func (this StatusCommand) status(paths []string, tagged []string, untagged []string) ([]string, []string, error) {
    db, error := OpenDatabase(databasePath())
    if error != nil {
        return nil, nil, error
    }

    return this.statusRecursive(db, paths, tagged, untagged)
}

func (this StatusCommand) statusRecursive(db *Database, paths []string, tagged []string, untagged []string) ([]string, []string, error) {
    for _, path := range paths {
        fileInfo, error := os.Lstat(path)
        if error != nil {
            return nil, nil, error
        }

        if fileInfo.IsRegular() {
            absPath, error := filepath.Abs(path)
            if error != nil {
                return nil, nil, error
            }

            file, error := db.FileByPath(absPath)
            if error != nil {
                return nil, nil, error
            }

            if file == nil {
                untagged = append(untagged, absPath)
            } else {
                tagged = append(tagged, absPath)
            }
        } else if fileInfo.IsDirectory() {
            file, error := os.Open(path)
            if error != nil {
                return nil, nil, error
            }
            defer file.Close()

            dirNames, error := file.Readdirnames(0)
            if error != nil {
                return nil, nil, error
            }

            childPaths := make([]string, len(dirNames))
            for index, dirName := range dirNames {
                childPaths[index] = filepath.Join(path, dirName)
            }

            tagged, untagged, error = this.statusRecursive(db, childPaths, tagged, untagged)
            if error != nil {
                return nil, nil, error
            }
        }
    }

    return tagged, untagged, nil
}