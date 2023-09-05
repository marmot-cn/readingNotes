type DynamicSubject interface {
    Execute(command string)
}

type RealDynamicSubject struct{}

func (rds *RealDynamicSubject) Execute(command string) {
    fmt.Printf("Executing command: %s\n", command)
}

type DynamicProxy struct {
    realSubject DynamicSubject
}

func (dp *DynamicProxy) Execute(command string) {
    fmt.Println("DynamicProxy: Logging command:", command)
    dp.realSubject.Execute(command)
}
