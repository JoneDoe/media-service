package services

import (
	"istorage/logger"
	"istorage/models"
	"istorage/repositories/cloud"
	"istorage/services/pool"
)

func Prowler(db *DbEngine, data []models.DownloadRequestItem) map[string]interface{} {
	p := pool.NewPool(len(data))

	for _, src := range data {
		p.AddTask(pool.NewTask(src, p.Channel))
	}

	p.Run()

	var filesList = make([]models.DownloadResponse, 0)
	var errors = make([]models.DownloadResponseError, 0)

	out := map[string]interface{}{
		"files":  &filesList,
		"errors": &errors,
	}

	for task := range p.Channel {
		p.ReceiveAnswer()

		if task.Err != nil {
			errors = append(errors, models.DownloadResponseError{
				Source:  task.Item.Source,
				Message: task.Err.Error(),
			})
		} else {
			file := task.Result

			go func(lf models.LocalFile, db *DbEngine) {
				if err := cloud.StoreFromLocal(lf); err != nil {
					logger.Error(err)
				}

				if err := db.CreateRecord(lf.ToJson()); err != nil {
					logger.Error(err)
				}

				lf.Remove()
			}(file, db)

			filesList = append(filesList, models.DownloadResponse{
				FileName: file.OriginalName(),
				Uuid:     file.Uuid(),
				Url:      cloud.GetFullUrl(file.FileSystemPath()),
			})
		}
	}

	return out
}
