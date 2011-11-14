package main

import (
	"errors"
	"fmt"
	"path/filepath"
)

type TagCommand struct{}

func (this TagCommand) Name() string {
	return "tag"
}

func (this TagCommand) Summary() string {
	return "applies one or more tags to a file"
}

func (this TagCommand) Help() string {
	return `  tmsu tag FILE TAG...

Tags the file FILE with the tag(s) specified.`
}

func (this TagCommand) Exec(args []string) error {
	if len(args) < 2 {
		return errors.New("File to tag and tags to apply must be specified.")
	}

	error := this.tagPath(args[0], args[1:])
	if error != nil {
		return error
	}

	return nil
}

// implementation

func (this TagCommand) tagPath(path string, tagNames []string) error {
	db, error := OpenDatabase(databasePath())
	if error != nil {
		return error
	}
	defer db.Close()

	absPath, error := filepath.Abs(path)
	if error != nil {
		return error
	}

	file, error := this.addFile(db, absPath)
	if error != nil {
		return error
	}

	for _, tagName := range tagNames {
		_, _, error = this.applyTag(db, path, file.Id, tagName)
		if error != nil {
			return error
		}
	}

	return nil
}

func (this TagCommand) applyTag(db *Database, path string, fileId uint, tagName string) (*Tag, *FileTag, error) {
	tag, error := db.TagByName(tagName)
	if error != nil {
		return nil, nil, error
	}

	if tag == nil {
		fmt.Printf("New tag '%v'\n", tagName)
		tag, error = db.AddTag(tagName)
		if error != nil {
			return nil, nil, error
		}
	}

	fileTag, error := db.FileTagByFileIdAndTagId(fileId, tag.Id)
	if error != nil {
		return nil, nil, error
	}

	if fileTag == nil {
		_, error := db.AddFileTag(fileId, tag.Id)
		if error != nil {
			return nil, nil, error
		}
	}

	return tag, fileTag, nil
}

func (this TagCommand) addFile(db *Database, path string) (*File, error) {
	fingerprint, error := Fingerprint(path)
	if error != nil {
		return nil, error
	}

	file, error := db.FileByPath(path)
	if error != nil {
		return nil, error
	}

	if file == nil {
		file, error = db.FileByFingerprint(fingerprint)
		if error != nil {
			return nil, error
		}

		if file != nil {
			fmt.Printf("Warning: file is a duplicate of a previously tagged file.\n")
		}

		file, error = db.AddFile(path, fingerprint)
		if error != nil {
			return nil, error
		}
	} else {
		if file.Fingerprint != fingerprint {
			db.UpdateFileFingerprint(file.Id, fingerprint)
		}
	}

	return file, nil
}