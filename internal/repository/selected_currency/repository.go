package selected_currency

type Repository interface {
	SaveSelectedCurrency(currency string, userId int64) error
	GetSelectedCurrency(userId int64) (SelectedCurrency, error)
}
