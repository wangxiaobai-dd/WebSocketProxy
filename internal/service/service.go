package service

type Service struct {
	ID      uint32 `json:"id"`
	IP      string `json:"ip"`
	Port    uint32 `json:"port"`
	ConnNum uint32 `json:"connNum"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start() {

}

func (s *Service) Stop() {

}

func (s *Service) Update() {

}
