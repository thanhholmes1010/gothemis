package themisallaka

type Actor interface {
	Receive(allakaCtx *AllakaContext)
}
