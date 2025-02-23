package scribe

type Scriber interface {
	BeginDescribe(desc string)
	EndDescribe()
	Done(done string)
}
