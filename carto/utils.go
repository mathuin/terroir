package carto

type CartoError error

func notnil(err error) bool {
	return (err != nil && err.Error() != "No Error")
}

type Float32Arr []float32

func (arr Float32Arr) min() (m float32) {
	m = arr[0]
	for _, v := range arr {
		if v < m {
			m = v
		}
	}
	return
}

func (arr Float32Arr) max() (m float32) {
	m = arr[0]
	for _, v := range arr {
		if v > m {
			m = v
		}
	}
	return
}

type Float64Arr []float64

func (arr Float64Arr) min() (m float64) {
	m = arr[0]
	for _, v := range arr {
		if v < m {
			m = v
		}
	}
	return
}

func (arr Float64Arr) max() (m float64) {
	m = arr[0]
	for _, v := range arr {
		if v > m {
			m = v
		}
	}
	return
}
