package transform

type Transformer interface {
	Transform(Data []byte)
}
