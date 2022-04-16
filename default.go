package logs_go

import "sync"

var (
	onceJ   sync.Once
	oncef   sync.Once
	stdLogJ LogJ
	stdLogf Logf
)

// DefaultLogJ
func DefaultLogJ() LogJ {
	onceJ.Do(func() {
		jcfg := NewLogJconfig()
		jcfg.Stdout = true
		stdLogJ, _ = jcfg.BuildLogJ()
	})
	return stdLogJ
}

// DefalutLogf
func DefalutLogf() Logf {
	oncef.Do(func() {
		fcfg := NewLogfConfig()
		fcfg.Stdout = true
		stdLogf, _ = fcfg.BuildLogf()
	})
	return stdLogf
}
