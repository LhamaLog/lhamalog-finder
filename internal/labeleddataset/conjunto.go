package labeleddataset

import (
	"encoding/json"
	"sort"
	"strings"
)

// DocumentacaoConjunto recebe a justificativa revisada pelo responsável pela base
// (007-12 B). Não inventa nome ou formato de um manifesto físico.
type DocumentacaoConjunto struct{ Conflitos string }

type entradaRotulada struct {
	Nivel, Mensagem, Lote, Registro, Tenant, Produto string
}

func entradaDe(c CamposEntrada) entradaRotulada {
	return entradaRotulada{c.LevelText, c.MessageText, c.BatchID, c.RecordID, c.TenantID, c.ProductID}
}

type chaveRotulada struct {
	Tipo    string
	Entrada entradaRotulada
	Membros string
}

// Compara identidade e conteúdo, sem produzir assinatura de ocorrência (008).
func chaveDe(r Registro) chaveRotulada {
	membros := make([]string, 0, len(r.Members))
	for _, m := range r.Members {
		valor, _ := json.Marshal(entradaDe(m)) // estrutura contém somente strings
		membros = append(membros, string(valor))
	}
	sort.Strings(membros)
	valor, _ := json.Marshal(membros)
	return chaveRotulada{r.ItemType, entradaDe(r.CamposEntrada), string(valor)}
}

// ValidarConjunto também cobre duplicatas em arquivos ou partições diferentes.
// Documentação pode acompanhar o registro ou ser fornecida pelo chamador; seu
// conteúdo e autoria continuam sujeitos à revisão da base, não a esta heurística.
func ValidarConjunto(registros []Registro, documentacao ...DocumentacaoConjunto) error {
	justificado := false
	for _, d := range documentacao {
		justificado = justificado || strings.TrimSpace(d.Conflitos) != ""
	}
	type estado struct {
		classe                ClassLabel
		conflito, documentado bool
	}
	vistos := make(map[chaveRotulada]estado)
	for _, r := range registros {
		if err := ValidateRegistro(r); err != nil {
			return err
		}
		chave := chaveDe(r)
		e, existe := vistos[chave]
		if !existe {
			e.classe = r.ClassLabel
		} else if e.classe != r.ClassLabel {
			e.conflito = true
		}
		e.documentado = e.documentado || strings.TrimSpace(r.Notas) != "" || strings.TrimSpace(r.Documentacao) != ""
		vistos[chave] = e
	}
	for _, e := range vistos {
		if e.conflito && !e.documentado && !justificado {
			return ErrRegistro
		}
	}
	return nil
}
