package ast

type Statement interface {
	Element
	isStatement()
}

// ReturnStatement

type ReturnStatement struct {
	Expression Expression
	StartPos   *Position
	EndPos     *Position
}

func (s *ReturnStatement) StartPosition() *Position {
	return s.StartPos
}

func (s *ReturnStatement) EndPosition() *Position {
	return s.EndPos
}

func (*ReturnStatement) isStatement() {}

func (s *ReturnStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitReturnStatement(s)
}

// IfStatement

type IfStatement struct {
	Test     Expression
	Then     *Block
	Else     *Block
	StartPos *Position
	EndPos   *Position
}

func (s *IfStatement) StartPosition() *Position {
	return s.StartPos
}

func (s *IfStatement) EndPosition() *Position {
	return s.EndPos
}

func (*IfStatement) isStatement() {}

func (s *IfStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitIfStatement(s)
}

// WhileStatement

type WhileStatement struct {
	Test     Expression
	Block    *Block
	StartPos *Position
	EndPos   *Position
}

func (s *WhileStatement) StartPosition() *Position {
	return s.StartPos
}

func (s *WhileStatement) EndPosition() *Position {
	return s.EndPos
}

func (*WhileStatement) isStatement() {}

func (s *WhileStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitWhileStatement(s)
}

// AssignmentStatement

type AssignmentStatement struct {
	Target   Expression
	Value    Expression
	StartPos *Position
	EndPos   *Position
}

func (s *AssignmentStatement) StartPosition() *Position {
	return s.StartPos
}

func (s *AssignmentStatement) EndPosition() *Position {
	return s.EndPos
}

func (*AssignmentStatement) isStatement() {}

func (s *AssignmentStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitAssignment(s)
}

// ExpressionStatement

type ExpressionStatement struct {
	Expression Expression
}

func (s *ExpressionStatement) StartPosition() *Position {
	return s.Expression.StartPosition()
}

func (s *ExpressionStatement) EndPosition() *Position {
	return s.Expression.EndPosition()
}

func (*ExpressionStatement) isStatement() {}

func (s *ExpressionStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitExpressionStatement(s)
}
