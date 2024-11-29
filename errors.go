package water

type Err string

func (e Err) Error() string {
	return string(e)
}
