package paintserver

import (
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserID 用户 ID
type UserID = uint64

// M 只是一个缩略写法
type M = bson.M

// DBName 是默认 QPaint Database Name
var DBName = "qpaint"

type QShape map[string]interface{}

func (shape QShape) GetID() string {
	if val, ok1 := shape["id"]; ok1 {
		if id, ok2 := val.(string); ok2 { // val 转为 string 类型，val.(type) 将 val 转为 type 类型
			return id
		}
	}
	return ""
}

func (shape QShape) setID(id string) {
	shape[id] = id
}

// ---------------------------------------------------

type shapeOnDrawing struct {
	front *shapeOnDrawing
	back  *shapeOnDrawing
	data  Shape
}

func (p *shapeOnDrawing) init() {
	p.front, p.back = p, p
}

func (p *shapeOnDrawing) insertFront(shape *shapeOnDrawing) {
	shape.back = p.back
	shape.front = p.front
	p.front.back = shape
	p.front = shape
}

func (p *shapeOnDrawing) insertBack(shape *shapeOnDrawing) {
	shape.front = p
	shape.back = p.back
	p.back.front = shape
	p.back = shape
}

func (p *shapeOnDrawing) moveFront() {
	front := p.front
	p.delete()
	front.insertFront(p)
}

func (p *shapeOnDrawing) moveBack() {
	back := p.back
	p.delete()
	back.insertBack(p)
}

func (p *shapeOnDrawing) delete() {
	p.front.back = p.back
	p.back.front = p.front
	p.back, p.front = nil, nil
}

// ---------------------------------------------------

type Drawing struct {
	id      bson.ObjectId
	session *mgo.Session
}

func newDrawing(id bson.ObjectId, session *mgo.Session) *Drawing {
	p := &Drawing{
		id:      id,
		session: session,
	}
	return p
}

func (p *Drawing) GetID() string {
	return p.id.Hex()
}

// Sync 同步 drawing 的修改
func (p *Drawing) Sync(shapes []ShapeID, changes []Shape) (err error) {
	c := p.session.Copy()
	defer c.Close()
	shapeColl := c.DB(DBName).C("shape")
	for _, change := range changes {
		spid := change.GetID()
		_, err = shapeColl.Upsert(M{
			"dgid": p.id,
			"spid": spid,
		}, M{
			"dgid":  p.id,
			"spid":  spid,
			"shape": change,
		})
		if err != nil {
			return mgoError(err)
		}
	}
	drawingColl := c.DB(DBName).C("drawing")
	return mgoError(drawingColl.UpdateId(p.id, M{
		"$set": M{"shapes": shapes},
	}))
}

func (p *Drawing) Add(shape Shape) (err error) {
	c := p.session.Copy()
	defer c.Close()
	shapeColl := c.DB(DBName).C("shape")
	spid := shape.GetID()
	err = shapeColl.Insert(M{
		"dgid":  p.id,
		"spid":  spid,
		"shape": shape,
	})
	if err != nil {
		return mgoError(err)
	}
	drawingColl := c.DB(DBName).C("drwaing")
	return mgoError(drawingColl.Update(p.id, M{
		"$push": M{"shapes": spid},
	}))
}

func (p *Drawing) List() (shapes []Shape, err error) {
	c := p.session.Copy()
	defer c.Close()
	var result []struct {
		ID    string `bson:"spid"`
		Shape QShape `bson:"shape"`
	}
	shapeColl := c.DB(DBName).C("shape")
	shapeColl.Find(M{
		"dgid": p.id,
	}).Select(M{
		"spid": 1, "shape": 1,
	}).All(&result)
	if err != nil {
		return nil, mgoError(err)
	}
	shapes = make([]Shape, len(result))
	for i, item := range result {
		item.Shape.setID(item.ID)
		shapes[i] = item.Shape
	}
	return
}

// Get 取出某个图形。
func (p *Drawing) Get(id ShapeID) (shape Shape, err error) {
	c := p.session.Copy()
	defer c.Close()
	var o struct {
		Shape QShape `bson:"shape"`
	}
	shapeColl := c.DB(DBName).C("shape")
	err = shapeColl.Find(M{
		"dgid": p.id,
		"spid": id,
	}).Select(M{
		"shape": 1,
	}).One(&o)
	if err != nil {
		return nil, mgoError(err)
	}
	o.Shape.setID(id)
	return o.Shape, nil
}

// Set 修改某个图形。
func (p *Drawing) Set(id ShapeID, shape Shape) (err error) {
	if shape.GetID() != "" {
		return syscall.EINVAL
	}
	c := p.session.Copy()
	defer c.Close()
	shapeColl := c.DB(DBName).C("shape")
	return mgoError(shapeColl.Update(M{
		"dgid": p.id,
		"spid": id,
	}, M{
		"$set": M{"shape": shape},
	}))
}

func (p *Drawing) SetZorder(id ShapeID, zorder string) (err error) {
	return nil // TODO
}

// Delete 删除某个图形。
func (p *Drawing) Delete(id ShapeID) (err error) {
	c := p.session.Copy()
	defer c.Close()
	drawingColl := c.DB(DBName).C("drawing")
	err = drawingColl.UpdateId(p.id, M{
		"$pull": M{"shapes": id},
	})
	if err != nil {
		return mgoError(err)
	}
	shapeColl := c.DB(DBName).C("shape")
	return mgoError(shapeColl.Remove(M{
		"dgid": p.id,
		"spid": id,
	}))
}

// ---------------------------------------------------

type Document struct {
	mutex   sync.Mutex
	data    map[string]*Drawing
	session *mgo.Session
}

func NewDocument() *Document {
	drawings := make(map[string]*Drawing)
	return &Document{
		data: drawings,
	}
}

func (p *Document) Add(uid UserID) (drawing *Drawing, err error) {
	c := p.session.Copy()
	defer c.Close()
	drawingColl := c.DB(DBName).C("drawing")
	id := bson.NewObjectId()
	err = drawingColl.Insert(M{
		"_id":    id,
		"uid":    uid,
		"shapes": []ShapeID{},
	})

	if err != nil {
		return
	}
	return newDrawing(id, p.session), nil
}

func (p *Document) Get(id string) (drawing *Drawing, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	drawing, ok := p.data[id]
	if !ok {
		return nil, syscall.ENOENT
	}
	return
}

func (p *Document) Delete(id string) (err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.data, id)
	return
}

// ---------------------------------------------------

var (
	idDrawingBase int64 = 10000
)

func makeDrawingID() string {
	id := atomic.AddInt64(&idDrawingBase, 1)
	return strconv.Itoa(int(id))
}

// ---------------------------------------------------
