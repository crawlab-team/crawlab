package challenge

type CreateCustomizedSpiderService struct {
}

func (s *CreateCustomizedSpiderService) Check() (bool, error) {
	return true, nil
}
