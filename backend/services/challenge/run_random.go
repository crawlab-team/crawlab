package challenge

type RunRandomService struct {
}

func (s *RunRandomService) Check() (bool, error) {
	return true, nil
}
