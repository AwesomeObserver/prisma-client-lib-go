package prisma

import (
	"fmt"
	"strings"
)

type operationType uint8

const (
	query operationType = iota + 1
	mutation
	subscription
)

type argumentList []argument

func (l argumentList) format(b *strings.Builder) {
	if len(l) == 0 {
		return
	}
	b.WriteByte('(')
	for i, arg := range l {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(arg.name)
		b.WriteString(": ")
		b.WriteString(arg.value)
	}
	b.WriteByte(')')
}

type fieldList []field

func (l fieldList) format(b *strings.Builder) {
	for _, f := range l {
		f.format(b)
		b.WriteByte('\n')
	}
}

type operation struct {
	typ       operationType
	name      string
	arguments argumentList
	fields    fieldList
}

type argument struct {
	name  string
	value string
}

func formatOperation(op *operation) (string, error) {
	// TODO(dh): verify that all names are valid (e.g. don't contain spaces)

	b := &strings.Builder{}
	switch op.typ {
	case query:
		b.WriteString("query")
	case mutation:
		b.WriteString("mutation")
	case subscription:
		b.WriteString("subscription")
	default:
		return "", fmt.Errorf("invalid operation type %q", op.typ)
	}

	b.WriteByte(' ')
	b.WriteString(op.name)
	op.arguments.format(b)
	b.WriteString(" {\n")
	op.fields.format(b)
	b.WriteByte('}')

	return b.String(), nil
}

type field interface {
	format(b *strings.Builder)
	isField()
}

type ScalarField struct {
	name      string
	arguments argumentList
}

func (f ScalarField) format(b *strings.Builder) {
	b.WriteString(f.name)
	f.arguments.format(b)
}

func (ScalarField) isField() {}

type ObjectField struct {
	name      string
	arguments argumentList
	fields    fieldList
}

func (f ObjectField) format(b *strings.Builder) {
	b.WriteString(f.name)
	f.arguments.format(b)
	b.WriteString(" {\n")
	f.fields.format(b)
	b.WriteString("}")
}

func (ObjectField) isField() {}
