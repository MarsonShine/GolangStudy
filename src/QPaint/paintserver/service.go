package paintserver

import (
	"net/http"

	"github.com/http/restrpc"
)

type Service struct {
	doc *Document
}

func NewService(doc *Document) (p *Service) {
	p = &Service{doc: doc}
	return
}

func SomeMetho() {
	return
}

var routeTable = [][2]string{
	{"POST /drawings", "PostDrawings"},
	{"GET /drawings/*", "GetDrawing"},
	{"DELETE /drawings/*", "DeleteDrawing"},
	{"POST /drawings/*/sync", "PostDrawingSync"},
	{"POST /drawings/*/shapes", "PostShapes"},
	{"GET /drawings/*/shapes/*", "GetShape"},
	{"POST /drawings/*/shapes/*", "PostShape"},
	{"DELETE /drawings/*/shapes/*", "DeleteShape"},
}

func (p *Service) PostShapes(aShape *serviceShape, env *restrpc.Env) (err error) {
	id := env.Args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		return
	}

	return drawing.Add(aShape.Get())
}

func (p *Service) PostDrawingSync(ds *serviceDrawingSync, env *restrpc.Env) (err error) {
	return
}

type serviceShape struct {
	ID      string       `json:"id"`
	Path    *pathData    `json:"path,omitempty"`
	Line    *lineData    `json:"line,omitempty"`
	Rect    *rectData    `json:"rect,omitempty"`
	Ellipse *ellipseData `json:"ellipse,omitempty"`
}

func (p *serviceShape) Get() Shape {
	if p.Path != nil {
		return &Path{shapeBase: shapeBase{p.ID}, pathData: *p.Path}
	}
	if p.Line != nil {
		return &Line{shapeBase: shapeBase{p.ID}, lineData: *p.Line}
	}
	if p.Rect != nil {
		return &Rect{shapeBase: shapeBase{p.ID}, rectData: *p.Rect}
	}
	if p.Ellipse != nil {
		return &Ellipse{shapeBase: shapeBase{p.ID}, ellipseData: *p.Ellipse}
	}
	return nil
}

type serviceDrawingSync struct {
	Changes []serviceShape `json:"changes"`
	Shapes  []ShapeID      `json:"shapes"`
}

func Main() {
	doc := NewDocument()
	service := NewService(doc)
	router := restrpc.Router{}
	http.ListenAndServe(":9999", router.Register(service, routeTable))
}
