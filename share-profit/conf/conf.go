package conf


type Config struct {
	Repo *Repo
	Mysql *MySql
	Lotus *Lotus
	Redis *Redis
	Project *Project
}

type Repo struct{
	UploadPath string
	MaxSize int64
}

type MySql struct {
	Name string
	Host string
	Username string
	Password string
}

type Lotus struct {
	Host string
	Token string
}

type Redis struct {
	Host string
	Password string
}

type Project struct {
	Host string
}
