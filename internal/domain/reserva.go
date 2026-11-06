package domain

// Quadro guarda a reserva exclusiva de cada item.
// A parada de uma instância entra como fato. Não há prazo numérico.
type Quadro struct {
	dono   map[ID]string
	itens  map[ID]Item
	livres map[ID]struct{}
}

// NovoQuadro começa com os itens livres informados.
func NovoQuadro(itens []Item) *Quadro {
	q := &Quadro{
		dono:   map[ID]string{},
		itens:  make(map[ID]Item, len(itens)),
		livres: make(map[ID]struct{}, len(itens)),
	}
	for _, item := range itens {
		item.Estado = Livre
		q.itens[item.ID] = item
		q.livres[item.ID] = struct{}{}
	}
	return q
}

// Tomar entrega no máximo um item ainda não reservado.
// O segundo tomador recebe outro item livre, ou nenhum.
func (q *Quadro) Tomar(instancia string) (Item, bool) {
	for id := range q.livres {
		delete(q.livres, id)
		item := q.itens[id]
		item.Estado = EmPasso
		q.itens[id] = item
		q.dono[id] = instancia
		return item, true
	}
	return Item{}, false
}

// Marcar reserva o item para a instância. Falso se outra instância já o tem.
func (q *Quadro) Marcar(instancia string, item Item) bool {
	if dono, ok := q.dono[item.ID]; ok && dono != instancia {
		return false
	}
	delete(q.livres, item.ID)
	item.Estado = EmPasso
	q.itens[item.ID] = item
	q.dono[item.ID] = instancia
	return true
}

// LiberarItem devolve um item à condição livre, sem conclusão.
func (q *Quadro) LiberarItem(id ID) {
	delete(q.dono, id)
	item, ok := q.itens[id]
	if !ok {
		return
	}
	item.Estado = Livre
	q.itens[id] = item
	q.livres[id] = struct{}{}
}

// ConcluirItem tira o item da reserva e da lista livre, sem novo dono.
func (q *Quadro) ConcluirItem(id ID) {
	delete(q.dono, id)
	delete(q.livres, id)
	if item, ok := q.itens[id]; ok {
		item.Estado = Concluido
		q.itens[id] = item
	}
}

// Dono informa a instância que reserva o item, se houver.
func (q *Quadro) Dono(id ID) (string, bool) {
	inst, ok := q.dono[id]
	return inst, ok
}

// SoltarInstancia devolve à condição livre os itens da instância que parou.
// O item não fica concluído. Nenhum prazo numérico se aplica.
func (q *Quadro) SoltarInstancia(instancia string) []Item {
	var soltos []Item
	for id, dono := range q.dono {
		if dono != instancia {
			continue
		}
		delete(q.dono, id)
		item := q.itens[id]
		item.Estado = Livre
		q.itens[id] = item
		q.livres[id] = struct{}{}
		soltos = append(soltos, item)
	}
	return soltos
}

// Encerrar tira o item da reserva sem devolvê-lo a livre e sem concluir.
func (q *Quadro) Encerrar(id ID) {
	delete(q.dono, id)
	delete(q.livres, id)
	if item, ok := q.itens[id]; ok {
		item.Estado = NaoConcluido
		q.itens[id] = item
	}
}

// Livre informa se o item pode ser tomado.
func (q *Quadro) Livre(id ID) bool {
	_, ok := q.livres[id]
	return ok
}

// Estado devolve o estado guardado do item.
func (q *Quadro) Estado(id ID) string {
	return q.itens[id].Estado
}
