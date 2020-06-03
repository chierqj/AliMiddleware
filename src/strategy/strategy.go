package strategy

// Strategy for sth.
type Strategy interface {
	ExecuteP1() map[string][]string
	ExecuteP2() map[string][]string
}
