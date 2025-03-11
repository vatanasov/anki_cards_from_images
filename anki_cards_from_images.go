package anki_cards_from_images

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type QuestionCouple struct {
	answer   string
	question string
}

func (q *QuestionCouple) ToRecord() []string {
	return []string{
		fmt.Sprintf("<img src=\"%s\">", q.question),
		fmt.Sprintf("<img src=\"%s\">", q.answer),
	}
}

type QuestionsList []QuestionCouple

func (q *QuestionsList) Add(question string, answer string) {
	*q = append(*q, QuestionCouple{question: question, answer: answer})
}

func Process(dirname string, w io.Writer) error {
	questionsList, err := processFiles(dirname)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(w)
	for _, question := range *questionsList {
		err := writer.Write(question.ToRecord())
		if err != nil {
			return fmt.Errorf("error writing record: %v", err)
		}
	}

	writer.Flush()

	return nil
}

func processFiles(dirname string) (*QuestionsList, error) {
	files := &QuestionsList{}
	err := filepath.WalkDir(
		dirname, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, "a.jpg") {
				return nil
			}

			if strings.HasSuffix(path, "q.jpg") {
				answerFile := strings.TrimSuffix(path, "q.jpg") + "a.jpg"
				// check if file exists
				if _, err := os.Stat(answerFile); err == nil {
					files.Add(path, answerFile)
					return nil
				}
				return fmt.Errorf("missing answer file: %s", answerFile)
			}

			return fmt.Errorf("unexpected file: %s", path)
		},
	)

	return files, err
}
