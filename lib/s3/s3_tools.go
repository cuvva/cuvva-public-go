package s3

import (
	"errors"
	"regexp"
)

var regSafeCharactersExclusion = regexp.MustCompile(`[^a-zA-Z0-9!\-_.*'()/]`)

func ReplaceUnsafeKeyCharacters(s3Key, replacementCharacter string) (*string, error) {
	// Based on https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingMetadata.html
	// if provided replacement character was not safe it would have characters replaced and not equal to original anymore
	testReplacementString := regSafeCharactersExclusion.ReplaceAllString(replacementCharacter, "-")
	if testReplacementString != replacementCharacter {
		return nil, errors.New("replacement character is not safe")
	}

	out := regSafeCharactersExclusion.ReplaceAllString(s3Key, replacementCharacter)

	return &out, nil
}
