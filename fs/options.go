package fs

import "time"

type Option func(m Manager)

func WithFilerUrl(url string) Option {
	return func(m Manager) {
		m.SetFilerUrl(url)
	}
}

func WithFilerAuthKey(authKey string) Option {
	return func(m Manager) {
		m.SetFilerAuthKey(authKey)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(m Manager) {
		m.SetTimeout(timeout)
	}
}

func WithWorkerNum(num int) Option {
	return func(m Manager) {
		m.SetWorkerNum(num)
	}
}

func WithRetryInterval(interval time.Duration) Option {
	return func(m Manager) {
		m.SetRetryInterval(interval)
	}
}

func WithRetryNum(num int) Option {
	return func(m Manager) {
		m.SetRetryNum(num)
	}
}

func WithMaxQps(qps int) Option {
	return func(m Manager) {
		m.SetMaxQps(qps)
	}
}
