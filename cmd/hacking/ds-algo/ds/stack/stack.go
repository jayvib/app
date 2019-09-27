package stack

type Stack interface {
	Push(data interface{})
	Pop() (data interface{})
	Top() (data interface{})
	Len() int
	IsEmpty() bool
}

type Container struct {
	list Stack
}

func (c *Container) Push(d interface{}) {
	c.list.Push(d)
}

func (c *Container) Pop() (data interface{}) {
	return c.list.Pop()
}

func (c *Container) Len() int {
	return c.list.Len()
}

func (c *Container) Top() (data interface{}) {
	return c.list.Top()
}

