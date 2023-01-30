package database

import(

)

var (
	ErrCantFindProduct = errors.New("Can't find the product")
	ErrCantDecodeProducts = errors.New("Can't find the product")
	ErrUserIdIsNotValid = errors.New("This user is not valid")
	ErrCantUpdateUser = errors.New("Can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("Can't remove this item from the cart")
	ErrCantGetItem = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem = errors.New("can't update the purchase")
)