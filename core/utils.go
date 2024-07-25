package core

import (
	"bytes"
	"io"
)

func ApplyTransforms(content []byte, transforms ...func(io.Reader, io.Writer) error) ([]byte, error) {
	var processedContent bytes.Buffer
	currentContent := bytes.NewReader(content)

	for _, transform := range transforms {
		processedContent.Reset()
		err := transform(currentContent, &processedContent)
		if err != nil {
			return nil, err
		}
		currentContent = bytes.NewReader(processedContent.Bytes())
	}

	return processedContent.Bytes(), nil
}
