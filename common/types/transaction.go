package types 

type Transaction struct{
	Id string
	AccountName string 
} 

func (tx *Transaction) Hash() Hash {
	return Hash{}
}
