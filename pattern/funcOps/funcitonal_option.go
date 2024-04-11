package funcOps

type config struct {
	// Required
	foo, bar string
	// Optional
	fizz, bazz int
}

func (c *config) WithFizz(fizz int) *config {
	c.fizz = fizz
	return c
}

func (c *config) WithBazz(bazz int) *config {
	c.bazz = bazz
	return c
}

// NewConfig 是新建Config的限制方法, 必须传入foo和bar
func NewConfig(foo, bar string) *config {
	return &config{
		foo: foo,
		bar: bar,
	}
}

func Do(c *config) {}

// calling 这里是功能选项模式的实现, NewConfig时, 通过是否继续调用With开头的方法,
// 来进行调用方的个性化配置
func calling() {
	c := NewConfig("foo", "bar").WithBazz(0)
	Do(c)
}
