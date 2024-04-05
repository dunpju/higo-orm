package event

import (
	"sync"
)

func init() {
	events = newEventRepository()
}

var (
	events *EventRepository
)

const (
	BeforeInsert EventType = iota + 1
	AfterInsert
	BeforeUpdate
	AfterUpdate
	BeforeDelete
	AfterDelete
	BeforeSelect
	AfterSelect
)

func eventRegister(iEvent IEvent) {
	events.add(iEvent.Event(), iEvent)
}

func EventPoint(event EventType, data EventData) {
	e := events.get(event)
	if e != nil {
		e.Handle(data)
	}
}

type IEvent interface {
	Event() EventType
	Handle(data EventData)
}

type EventType int

type EventRepository struct {
	sort    []EventType
	syncMap sync.Map
}

func newEventRepository() *EventRepository {
	return &EventRepository{sort: []EventType{}}
}

func (this *EventRepository) Len() int {
	return len(this.sort)
}

func (this *EventRepository) get(event EventType) IEvent {
	e, ok := this.syncMap.Load(event)
	if ok {
		return e.(IEvent)
	}
	return nil
}

func (this *EventRepository) add(event EventType, iEvent IEvent) {
	_, ok := this.syncMap.Load(event)
	if !ok {
		this.sort = append(this.sort, event)
	}
	this.syncMap.Store(event, iEvent)
}

type EventData struct {
	Table        string
	Sql          string
	Args         []interface{}
	Err          error
	LastInsertId int64
	RowsAffected int64
}

func NewEventData(table string, sql string, args []interface{}, err error, lastInsertId int64, rowsAffected int64) EventData {
	return EventData{Table: table, Sql: sql, Args: args, Err: err, LastInsertId: lastInsertId, RowsAffected: rowsAffected}
}

type EventHandle func(data EventData)

type Event struct {
	eventType EventType
	handle    EventHandle
}

func AddEvent(eventType EventType, handle EventHandle) {
	eventRegister(&Event{eventType: eventType, handle: handle})
}

func (this *Event) Event() EventType {
	return this.eventType
}

func (this *Event) Handle(data EventData) {
	this.handle(data)
}
