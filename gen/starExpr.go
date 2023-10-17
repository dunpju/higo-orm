package gen

type starExprCollect struct {
	collect []string
}

func newStarExprCollect() *starExprCollect {
	return &starExprCollect{collect: make([]string, 0)}
}

func (this *starExprCollect) append(starExpr string) {
	has := false
	var index int
	for i, s := range this.collect {
		if s == starExpr {
			has = true
			index = i
			break
		}
	}
	if !has {
		this.collect = append(this.collect, starExpr)
	} else {
		this.collect[index] = starExpr
	}
}
