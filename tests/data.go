//go:generate csval

package tests

// FooStruct is a sample struct to be validated
type FooStruct struct {
	Name        string       `csval:"req"`
	Email       string       `csval:"req,email"`
	Password    string       `csval:"min(3),max(11)"`
	ConfirmPass string       `csval:"req,equals(Password)"`
	Age         int          `csval:"min(18),max(65)"`
	Sub         BarSubStruct `csval:"obj"`
}

// BarSubStruct is a sample struct to be validated
type BarSubStruct struct {
	IP   string `csval:"req,ip"`
	Port int
}
