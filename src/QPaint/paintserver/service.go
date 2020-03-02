package paintserver

import (
	"log"
	"net/http"
	"syscall"

	"github.com/http/restrpc"
	"gopkg.in/mgo.v2"
)

func mgoError(err error) error {
	if err == nil {
		return nil
	}
	if err == mgo.ErrNotFound {
		return syscall.ENOENT
	}
	if mgo.IsDup(err) {
		return syscall.EEXIST
	}
	return err
}

type RouteTable map[string]func(w http.ResponseWriter, req *http.Request, args []string)

type Service struct {
	doc        *Document
	routeTable RouteTable
}

func NewService(doc *Document) (p *Service) {
	p = &Service{doc: doc}
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

func (p *Service) PostDrawings(w http.ResponseWriter, req *http.Request, args []string) (m M, err error) {
	log.Println(req.Method, req.URL)
	drawing, err := p.doc.Add(1) // TODO
	if err != nil {
		return
	}
	return M{"id": drawing.id}, nil
}

func (p *Service) DeleteDrawing(env *restrpc.Env) (err error) {
	id := env.Args[0]
	return p.doc.Delete(id)
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
