package price

type PriceSender interface {
	Send(event PriceEvent) error
}
