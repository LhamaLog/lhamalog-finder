package domain

// AssinaturaLog documenta os componentes do cálculo canônico application/assinatura.
// Valor não é parseado enquanto o escape de barra invertida estiver pendente.
type AssinaturaLog struct{ Nivel, MensagemNormalizada, ServiceName, Host, Valor, VersaoRegra string }
