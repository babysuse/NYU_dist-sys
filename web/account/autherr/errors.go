package autherr

type WrongAuth struct{}

func (e *WrongAuth) Error() string {
	return "Wrong username or password"
}
