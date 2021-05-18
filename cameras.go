package redeye

var (
	cameras map[string]*Camera
)

type Camera struct {
	Name string
	Addrport string
}

func NewCamera(name string) *Camera {
	cam := &Camera{Name: name}
	return cam;
}
