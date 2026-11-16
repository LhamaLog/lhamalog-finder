package out

import (
	"context"

	"github.com/LhamaLog/lhamalog-finder/internal/domain"
)

// Efeito recebe a referência validada e o id da entrada, que acompanha a referência.
// A execução está em application.ServicoRecuperacao, que devolve a marca de efeito concluído.
// concluido autoriza a confirmação pela primeira tarefa. Esta porta não confirma a entrada.
// Indisponibilidade, falha de página, lote incompleto, configuração ausente ou a porta ausente
// mantêm a entrada pendente.
type Efeito interface {
	Executar(ctx context.Context, ref domain.Referencia) (concluido bool, err error)
}
