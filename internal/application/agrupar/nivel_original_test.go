package agrupar

import (
	"github.com/LhamaLog/lhamalog-finder/internal/domain"
	"testing"
)

func TestParticaoPreservaNivelOriginalSeparadoDaClasse(t *testing.T) {
	a := domain.LogClassificado{Apontamento: "a", Classe: domain.ERROR, TenantID: "t", ProductID: "p", ServiceName: "svc", Host: "h", NivelCampo: "INFO", Mensagem: "ERROR falha"}
	b := a
	b.Apontamento = "b"
	b.NivelCampo = "WARN"
	grupos, err := ParticionarLote([]domain.LogClassificado{a, b})
	if err != nil || len(grupos) != 2 {
		t.Fatalf("niveis distintos agrupados: grupos=%v erro=%v", grupos, err)
	}
	for i, nivel := range []string{"INFO", "WARN"} {
		if grupos[i].Chave.Assinatura != nivel+"|error falha|svc|h" || grupos[i].Chave.Classe != domain.ERROR {
			t.Fatalf("nivel ou classe alterados: %+v", grupos[i].Chave)
		}
	}
}

func TestParticaoMantemNivelAusenteSemInventarClasseComoNivel(t *testing.T) {
	log := domain.LogClassificado{Apontamento: "a", Classe: domain.INFO, TenantID: "t", ProductID: "p", ServiceName: "svc", Host: "h", Mensagem: "processado"}
	grupos, err := ParticionarLote([]domain.LogClassificado{log})
	if err != nil || len(grupos) != 1 || grupos[0].Chave.Assinatura != "|processado|svc|h" {
		t.Fatalf("nivel opcional substituido ou rejeitado: grupos=%v erro=%v", grupos, err)
	}
}
