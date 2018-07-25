package cooked

type untrusted interface {
	IsValid() error
}


func check(u untrusted) error {
	return u.IsValid()
}
