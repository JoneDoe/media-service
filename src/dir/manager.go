package dir

type FileManager interface {
	Filepath() string
	Abs() string
}

type Config struct {
	FileName, MimeType, Context string
}

func NewManager(cfg *Config) (FileManager, error) {
	return &DefaultManager{
		FileName: cfg.FileName,
		Context:  cfg.Context,
	}, nil
}
