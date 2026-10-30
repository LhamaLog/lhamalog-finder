package domain

// Resultado é o efeito do passo de consumo sobre um item.
type Resultado struct {
	Item                  Item
	Concluido             bool
	Disponibilizar        bool
	LivreDeNovo           bool
	EncerradoSemConclusao bool
	ForaDeContrato        bool
	FalhaLeitura          bool
}

// Iniciar leva um item livre ao meio do passo.
func Iniciar(item Item) (Item, bool) {
	if !item.Aceito() || item.Estado != Livre {
		return item, false
	}
	item.Estado = EmPasso
	return item, true
}

// Fim encerra o passo.
// Com o worker e a instância no ar, o passo livre termina em conclusão,
// se o código ainda não foi concluído.
// Parada no meio do passo não grava conclusão e devolve o item a livre.
// Código já concluído não gera segunda conclusão nem reserva restante.
func Fim(item Item, workerNoAr, instanciaNoAr, codigoJaConcluido bool) Resultado {
	if !item.Aceito() {
		return Resultado{Item: item, ForaDeContrato: true}
	}
	if item.Estado == "" {
		item.Estado = Livre
	}
	if item.Estado != EmPasso {
		item, ok := Iniciar(item)
		if !ok {
			return Resultado{Item: item, ForaDeContrato: !item.Aceito()}
		}
	}
	if !workerNoAr || !instanciaNoAr {
		item.Estado = Livre
		return Resultado{Item: item, LivreDeNovo: true}
	}
	if codigoJaConcluido {
		item.Estado = NaoConcluido
		return Resultado{Item: item, EncerradoSemConclusao: true}
	}
	item.Estado = Concluido
	return Resultado{Item: item, Concluido: true, Disponibilizar: true}
}

// ConcluirJuntos aplica o fim a vários itens.
// Dois itens com o mesmo valor que terminam juntos produzem uma conclusão
// e nenhuma reserva restante. Qual item vence não está no modelo:
// o resultado exigido é a contagem, não a identidade do vencedor.
func ConcluirJuntos(itens []Item) []Resultado {
	vistos := make(map[Codigo]bool, len(itens))
	saida := make([]Resultado, len(itens))
	for i, item := range itens {
		if item.Estado != EmPasso {
			item, _ = Iniciar(item)
		}
		ja := vistos[item.Codigo]
		res := Fim(item, true, true, ja)
		saida[i] = res
		if res.Concluido {
			vistos[item.Codigo] = true
		}
	}
	return saida
}

// SemConsumo é fila vazia ou falha de leitura: nenhuma conclusão.
func SemConsumo(falhaLeitura bool) Resultado {
	return Resultado{Concluido: false, FalhaLeitura: falhaLeitura}
}
