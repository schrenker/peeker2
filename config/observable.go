package config

type Observable interface {
	GetNames() []string
}

type DiskSlice []Disk

type Disk struct {
	Path string `yaml:"path"`
}

func (d DiskSlice) GetNames() []string {
	tmp := make([]string, len(d))
	for i := range d {
		tmp[i] = d[i].Path
	}
	return tmp
}

type ServiceSlice []Service

type Service struct {
	Name string `yaml:"name"`
}

func (s ServiceSlice) GetNames() []string {
	tmp := make([]string, len(s))
	for i := range s {
		tmp[i] = s[i].Name
	}
	return tmp
}
