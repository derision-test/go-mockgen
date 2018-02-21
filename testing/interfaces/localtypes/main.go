package localtypes

type LocalRewrite interface {
	Test(x X, y Y, z Z) (*X, *Y, *Z)
}
