package generator

import (
	pParser "protocoll/internal/domain/parser"
)

type Result struct {
	Variables []string
	Package   *pParser.Package
}
