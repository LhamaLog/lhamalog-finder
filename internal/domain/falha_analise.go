package domain

// FalhaAnalise preserva o motivo sem texto de log. SemRetry é definido
// explicitamente para produto ausente/inconsistente (008-16 A); não cria retry.
type FalhaAnalise struct {
	Motivo   string
	SemRetry bool
}

func (f FalhaAnalise) Error() string    { return f.Motivo }
func (f FalhaAnalise) Permanente() bool { return f.SemRetry }
