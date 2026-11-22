package domain

const MarcadorOmisso = "log com senha omitido"

// LinhaLogMascarada contém somente a saída do scanner; atributos OTel ficam fora.
type LinhaLogMascarada struct {
	Texto         string
	IndiceLinha   int
	LinhaRemovida bool
}
type ResultadoValidarSeguranca struct {
	LinhasMascaradas []LinhaLogMascarada
	TotalSensivel    int
	Erro             error
}
type ContextoLoteMascarado struct {
	BatchID                     string
	Linhas                      []LinhaLogMascarada
	VersaoCatalogo              string
	DadosSensiveisIdentificados int
	FalhaSeguranca              bool
}

// CatalogoSeguranca é a fronteira opaca do artefato. A ADR P-010-2 deve
// definir regex, motor, placeholder e serialização; este tipo não os escolhe.
// Transformar recebe uma linha e retorna o texto e o total aplicado pelo catálogo.
type CatalogoSeguranca struct {
	Versao      string
	Transformar func(string) (LinhaLogMascarada, int, error)
}
