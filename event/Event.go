package event

import (
	"fmt"
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
	BeforeCount
	AfterCount
	BeforeRaw
	AfterRaw
)

func eventRegister(iEvent IEvent) {
	events.add(iEvent.Event(), iEvent)
}

func Point(event EventType, data EventRecord) {
	e := events.get(event, data.Table)
	if e != nil {
		e.Handle(data)
	}
}

type IEvent interface {
	Event() EventType
	Table() string
	Handle(data EventRecord)
}

type EventType int

type EventRepository struct {
	syncMap sync.Map
}

func newEventRepository() *EventRepository {
	return &EventRepository{}
}

const keyFormat = "event:%d-table:%s"

func (this *EventRepository) get(event EventType, table string) IEvent {
	key := fmt.Sprintf(keyFormat, int(event), table)
	e, ok := this.syncMap.Load(key)
	if ok {
		return e.(IEvent)
	}
	return nil
}

func (this *EventRepository) add(event EventType, iEvent IEvent) {
	key := fmt.Sprintf(keyFormat, int(event), iEvent.Table())
	_, ok := this.syncMap.Load(key)
	if !ok {
		this.syncMap.Store(key, iEvent)
	}
}

type EventRecord struct {
	Table        string
	Sql          string
	Args         []interface{}
	Err          error
	LastInsertId int64
	RowsAffected int64
	Result       interface{}
}

func NewEventRecord(table string, sql string, args []interface{}, err error, lastInsertId int64, rowsAffected int64) EventRecord {
	return EventRecord{Table: table, Sql: sql, Args: args, Err: err, LastInsertId: lastInsertId, RowsAffected: rowsAffected}
}

func NewEventRecordResult(table string, sql string, args []interface{}, err error, result interface{}) EventRecord {
	return EventRecord{Table: table, Sql: sql, Args: args, Err: err, Result: result}
}

type EventHandle func(data EventRecord)

type Event struct {
	eventType EventType
	table     string
	handle    EventHandle
}

func AddEvent(eventType EventType, table string, handle EventHandle) {
	eventRegister(&Event{eventType: eventType, table: table, handle: handle})
}

func (this *Event) Event() EventType {
	return this.eventType
}

func (this *Event) Table() string {
	return this.table
}

func (this *Event) Handle(data EventRecord) {
	this.handle(data)
}