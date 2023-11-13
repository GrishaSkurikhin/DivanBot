package customerrors

type AlreadyRegistered struct{}

func (err AlreadyRegistered) Error() string {
	return "user already reg on film"
}

type UserNotRegistered struct{}

func (err UserNotRegistered) Error() string {
	return "user is not registered"
}