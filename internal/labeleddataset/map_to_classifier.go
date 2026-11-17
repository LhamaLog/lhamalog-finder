package labeleddataset

import "github.com/LhamaLog/lhamalog-finder/internal/domain"

// Predizer usa a implementação real da 006. O catálogo só-catálogo aguarda 006-7 D.
func Predizer(r Registro) (ClassLabel, error) {
	if err := ValidateRegistro(r); err != nil {
		return "", err
	}
	if r.ItemType == "log" {
		return domain.ClassificarLog(r.LevelText, r.MessageText, nil).Classe, nil
	}
	prevista := domain.INFO
	for _, m := range r.Members {
		prevista = domain.MaisGrave(prevista, domain.ClassificarLog(m.LevelText, m.MessageText, nil).Classe)
	}
	return prevista, nil
}
