package farmer

const (
	SCRIPT_CREATE = "create"
	SCRIPT_DEPLOY = "deploy"
	SCRIPT_TEST   = "test"
)

func (b *Box) runScript(key string) error {
	if _, ok := b.Scripts[key]; ok {
		return dockerExecOnContainer(b, []string{
			b.Home + "/" + b.Scripts[key],
		})
	}

	return nil
}
