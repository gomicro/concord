package report

type ChangeSet struct {
	changes []change
}

func (c *ChangeSet) Add(pre, post string) {
	c.changes = append(c.changes, change{
		pre:  pre,
		post: post,
	})
}

func (c *ChangeSet) PrintPre() {
	for i := range c.changes {
		PrintAdd(c.changes[i].pre)
		Println()
	}
}

func (c *ChangeSet) PrintPost() {
	for i := range c.changes {
		PrintSuccess(c.changes[i].post)
		Println()
	}
}

func (c *ChangeSet) HasChanges() bool {
	return len(c.changes) > 0
}

type change struct {
	pre  string
	post string
}
