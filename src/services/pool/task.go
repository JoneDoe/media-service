package pool

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"istorage/dir"
	"istorage/models"
)

type Task struct {
	Err     error
	Channel chan *Task
	Item    models.DownloadRequestItem
	Result  models.LocalFile
}

func NewTask(req models.DownloadRequestItem, ch chan *Task) *Task {
	return &Task{Item: req, Channel: ch}
}

func (t *Task) Run(wg *sync.WaitGroup) {
	file := createFileTarget(t.Item.BuildDestinationFileName())
	f := models.NewLocalFile(file)

	err := downloadFile(t.Item.Source, file, httpClient())
	if err == nil {
		fm, _ := dir.NewManager(&dir.Config{
			MimeType: models.FileTypeImage,
			FileName: f.Name(),
			Context:  "default",
		})

		f.SetPath(fm.Abs())

		t.Result = f
	} else {
		t.Err = err
		f.Remove()
	}

	t.Channel <- t

	wg.Done()
}

func httpClient() *http.Client {
	client := http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func createFileTarget(fileName string) *os.File {
	filePath := filepath.Join(os.TempDir(), fileName)
	file, _ := os.Create(filePath)

	return file
}

func downloadFile(url string, file *os.File, client *http.Client) error {
	resp, err := client.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if size == 0 {
		return errors.New("can`t write file")
	}

	defer file.Close()

	if err != nil {
		return err
	}

	return nil
}
