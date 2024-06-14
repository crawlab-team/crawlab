package color

import (
	"encoding/hex"
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/data"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/trace"
	"math/rand"
	"strconv"
	"strings"
)

func NewService() (svc interfaces.ColorService, err error) {
	var cl []*entity.Color
	cm := map[string]*entity.Color{}

	if err := json.Unmarshal([]byte(data.ColorsDataText), &cl); err != nil {
		return nil, trace.TraceError(err)
	}

	for _, c := range cl {
		cm[c.Name] = c
	}

	return &Service{
		cl: cl,
		cm: cm,
	}, nil
}

type Service struct {
	cl []*entity.Color
	cm map[string]*entity.Color
}

func (svc *Service) Inject() (err error) {
	return nil
}

func (svc *Service) GetByName(name string) (res interfaces.Color, err error) {
	res, ok := svc.cm[name]
	if !ok {
		return res, errors.ErrorModelNotFound
	}
	return res, err
}

func (svc *Service) GetRandom() (res interfaces.Color, err error) {
	if len(svc.cl) == 0 {
		hexStr, err := svc.getRandomColorHex()
		if err != nil {
			return res, err
		}
		return &entity.Color{Hex: hexStr}, nil
	}

	idx := rand.Intn(len(svc.cl))
	return svc.cl[idx], nil
}

func (svc *Service) getRandomColorHex() (res string, err error) {
	n := 6
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		arr[i], err = svc.getRandomHexChar()
		if err != nil {
			return res, err
		}
	}
	return strings.Join(arr, ""), nil
}

func (svc *Service) getRandomHexChar() (res string, err error) {
	n := rand.Intn(16)
	b := []byte(strconv.Itoa(n))
	h := make([]byte, 1)
	hex.Encode(h, b)
	return string(h), nil
}
