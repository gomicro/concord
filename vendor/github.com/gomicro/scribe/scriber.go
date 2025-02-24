package scribe

type Scriber interface {
	BeginDescribe(desc string)
	EndDescribe()
	Print(done string)
}
