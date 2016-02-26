package mgox

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/kimiazhu/log4go"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "reflect"
    "time"
)

type ErrorMode int

const (
    error_mode_isolation = 1
    error_mode_sharing = 2
)

type Dao struct {
    db        *mgo.Database
    uid       interface{}
    LastError error
    errMode   ErrorMode
}

func New(user ...interface{}) *Dao {
    d := new(Dao)
    d.SetUser(user...)
    d.errMode = error_mode_isolation
    return d
}

func Connect(configs ...interface{}) *Dao {
    return New(configs...).Connect()
}

func (d *Dao) SetUser(user ...interface{}) {
    if len(user) > 0 {
        d.uid = user[0]
    }
}

func (d *Dao) Connect() *Dao {

    if d.db != nil {
        return d
    }

    if d.errMode == 0 {
        d.errMode = error_mode_isolation
    }

    if d.errMode == error_mode_isolation || (d.errMode == error_mode_sharing && d.LastError == nil) {
        var err error
        d.db, err = GetDatabase()
        if err != nil {
            log4go.Error("Connect db fail: ", err)
            panic(err)
        }
    }

    return d
}

func (d *Dao) Close() {
    if d.db != nil && d.db.Session != nil {
        d.db.Session.Close()
        d.db = nil
        log4go.Debug("Closed DB connection succssfully")
    }
}

func (d *Dao) ShareError() *Dao {
    d.errMode = error_mode_sharing
    return d
}

func (d *Dao) Insert(docs ...interface{}) error {

    if d.errMode == error_mode_sharing && d.LastError != nil {
        return d.LastError
    }

    for i, _ := range docs {

        if docs[i] == nil {
            err := errors.New("cannot insert empty document")
            d.LastError = err
            return err
        }

        mType := reflect.TypeOf(docs[i])
        mValue := reflect.ValueOf(docs[i])
        mType, mValue = getElem(mType, mValue)

        if mType.Kind() == reflect.Struct && mValue.IsValid() {
            field := mValue.FieldByName("Id")
            if field.IsValid() && field.String() == "" && field.CanSet() {
                field.Set(reflect.ValueOf(bson.NewObjectId()))
            }
            now := reflect.ValueOf(time.Now())
            field = mValue.FieldByName("FirstCreated")
            if field.CanSet() {
                field.Set(now)
            }
            field = mValue.FieldByName("LastModified")
            if field.CanSet() {
                field.Set(now)
            }

            // This will force mongo db store localcreated field in a local time , not UTC time.
            field = mValue.FieldByName("LocalCreated")
            if field.CanSet() {
                _, offset := time.Now().Local().Zone()
                //				zoneName,offset:=time.Now().Local().Zone()
                //				println("zone name is:%s,offset is:%d",zoneName,offset)
                newtime := time.Now().Add(time.Duration(offset * 1000000000))
                //				fmt.Println("newtime:",newtime)
                nowlocal := reflect.ValueOf(newtime)
                field.Set(nowlocal)
            }

            if d.uid != nil {
                field = mValue.FieldByName("FirstCreator")
                if field.CanSet() {
                    field.Set(reflect.ValueOf(d.uid))
                }
                field = mValue.FieldByName("LastModifier")
                if field.CanSet() {
                    field.Set(reflect.ValueOf(d.uid))
                }
            }

        } else if mType.Kind() == reflect.Map {
            m, ok := docs[i].(bson.M)
            if !ok {
                m, _ = docs[i].(map[string]interface{})
            }
            if len(m) > 0 {
                if _, ok := m["_id"]; !ok {
                    m["_id"] = bson.NewObjectId()
                }
                now := time.Now()
                m["firstcreated"] = now
                m["lastmodified"] = now
                if d.uid != nil {
                    m["firstcreator"] = d.uid
                    m["lastmodifier"] = d.uid
                }
            }
        }
    }

    if d.db == nil {
        d.Connect()
        defer d.Close()
    }

    log4go.Debug("insert")

    err := d.db.C(getCollectionName(docs[0])).Insert(docs...)
    if err != nil {
        d.LastError = err
    }
    return err
}

func (d *Dao) Remove(collectionName interface{}, selector interface{}) error {

    if d.errMode == error_mode_sharing && d.LastError != nil {
        return d.LastError
    }

    if d.db == nil {
        d.Connect()
        defer d.Close()
    }

    log4go.Debug("remove")

    c := d.db.C(getCollectionName(collectionName))

    var id bson.ObjectId
    var err error

    if strId, ok := selector.(string); ok {
        if strId == "" {
            err = errors.New("id can't be empty")
            d.LastError = err
            return err
        }
        id = bson.ObjectIdHex(strId)
    } else if oId, ok := selector.(bson.ObjectId); ok {
        if oId == "" {
            err = errors.New("id can't be empty")
            d.LastError = err
            return err
        }
        id = oId
    } else if _, ok := selector.(bson.M); ok {

    } else if _, ok := selector.(map[string]interface{}); ok {

    } else {
        err = errors.New("unrecognized selector for remove")
        d.LastError = err
        return err
    }

    if id != "" {
        err = c.RemoveId(id)
    } else {
        _, err = c.RemoveAll(selector)
    }

    if err != nil {
        d.LastError = err
    }

    return err
}

func (d *Dao) isID(selectors ...interface{}) bool {
    if len(selectors) != 1 {
        return false
    }
    return !d.isM(selectors[0])
}

func (d *Dao) isM(v interface{}) bool {
    if _, ok := v.(bson.M); ok {
        return true
    } else if _, ok := v.(map[string]interface{}); ok {
        return true
    }
    return false
}

func (d *Dao) getM(values ...interface{}) bson.M {
    m := bson.M{}
    for i := 0; i < len(values); i++ {
        if _m, ok := values[i].(bson.M); ok {
            for key, value := range _m {
                m[key] = value
            }
        } else if _m, ok := values[i].(map[string]interface{}); ok {
            for key, value := range _m {
                m[key] = value
            }
        } else {
            if i == len(values) - 1 {
                break
            }
            m[values[i].(string)] = values[i + 1]
            i++
        }
    }
    return m
}

func (d *Dao) update(operator string, collectionName interface{}, selector interface{}, updates ...interface{}) error {

    if d.errMode == error_mode_sharing && d.LastError != nil {
        return d.LastError
    }

    if len(updates) == 0 {
        d.LastError = errors.New("updates can't be empty")
        return d.LastError
    }

    isStruct := false
    mType := reflect.TypeOf(updates[0])
    mValue := reflect.ValueOf(updates[0])
    mType, mValue = getElem(mType, mValue)
    if mType.Kind() == reflect.Struct && mValue.IsValid() {

        isStruct = true

        field := mValue.FieldByName("LastModified")
        if field.CanSet() {
            field.Set(reflect.ValueOf(time.Now()))
        }

        if d.uid != nil {
            field := mValue.FieldByName("LastModifier")
            if field.CanSet() {
                field.Set(reflect.ValueOf(d.uid))
            }
        }
    }

    var update bson.M

    if !isStruct {

        update = d.getM(updates...)

        if operator == "$inc" || operator == "$push" {
            m := bson.M{"lastmodified": time.Now()}
            if d.uid != nil {
                m["lastmodifier"] = d.uid
            }
            update = bson.M{"$set": m, operator: update}
        } else {
            if operator == "$set" || operator == "$update" {
                update["lastmodified"] = time.Now()
                if d.uid != nil {
                    update["lastmodifier"] = d.uid
                }
            }

            if operator != "update" {
                update = bson.M{operator: update}
            }
        }
    }

    if d.db == nil {
        d.Connect()
        defer d.Close()
    }

    var err error

    _collectionName := getCollectionName(collectionName)
    c := d.db.C(_collectionName)

    isId := false
    if _, ok := selector.(string); ok {
        isId = true
    } else if _, ok := selector.(bson.ObjectId); ok {
        isId = true
    } else if _, ok := selector.(bson.M); ok {

    } else if _, ok := selector.(map[string]interface{}); ok {

    } else {
        err = errors.New("unrecognized selector for update")
        d.LastError = err
        return err
    }

    if isStruct {
        if isId {
            log4go.Debug(fmt.Sprintf("[%s]collection=%s,id=%s,struct=%s", operator, _collectionName, selector, updates[0]))
            err = c.UpdateId(selector, updates[0])
        } else {
            log4go.Debug(fmt.Sprintf("[%s]collection=%s,selector=%s,struct=%s", operator, _collectionName, selector, updates[0]))
            _, err = c.UpdateAll(selector, updates[0])
        }
    } else {
        if isId {
            log4go.Debug(fmt.Sprintf("id=%s", selector))
            log4go.Debug(fmt.Sprintf("[%s]collection=%s,id=%s,update=%s", operator, _collectionName, selector, update))
            err = c.UpdateId(selector, update)
        } else {
            log4go.Debug(fmt.Sprintf("[%s]collection=%s,selector=%s,struct=%s", operator, _collectionName, selector, update))
            _, err = c.UpdateAll(selector, update)
        }
    }

    if err != nil {
        d.LastError = err
    }

    return err
}

func (d *Dao) Set(collectionName interface{}, selector interface{}, updates ...interface{}) error {
    return d.update("$set", collectionName, selector, updates...)
}

func (d *Dao) Inc(collectionName interface{}, selector interface{}, updates ...interface{}) error {
    return d.update("$inc", collectionName, selector, updates...)
}

func (d *Dao) Push(collectionName interface{}, selector interface{}, updates ...interface{}) error {
    return d.update("$push", collectionName, selector, updates...)
}

func (d *Dao) Pull(collectionName interface{}, selector interface{}, updates ...interface{}) error {
    return d.update("$pull", collectionName, selector, updates...)
}

func (d *Dao) Replace(collectionName interface{}, selector interface{}, updates ...interface{}) error {
    return d.update("update", collectionName, selector, updates...)
}

func (d *Dao) ReplaceDoc(doc interface{}) error {
    return d.ReplaceDocById(getObjectId(doc), doc)
}

func (d *Dao) ReplaceDocById(id interface{}, doc interface{}) error {
    return d.update("update", doc, id, doc)
}

type Query struct {
    collection string
    dao        *Dao
    ignoreNFE  bool
    queries    []interface{}
    sorts      []string
    page       *Page
    columns    []string
    exColumns  []string
    distinct   string
    limit      int
}

func (d *Dao) Find(queries ...interface{}) *Query {
    query := &Query{}
    query.dao = d
    query.queries = queries
    return query
}

func (d *Dao) Get(queries ...interface{}) *Query {
    return d.Find(queries...)
}

func (q *Query) IgnoreNFE() *Query {
    q.ignoreNFE = true
    return q
}

func (q *Query) Page(page *Page) *Query {
    if page != nil && page.Cursor >= 0 {
        q.page = page
    }
    return q
}

func (q *Query) Sort(sorts ...string) *Query {
    q.sorts = sorts
    return q
}

func (q *Query) Select(columns ...string) *Query {
    q.columns = columns
    return q
}

func (q *Query) SelectExclude(exColumns ...string) *Query {
    q.exColumns = exColumns
    return q
}

func (q *Query) Distinct(collectionName interface{}, distinct string, result interface{}) error {
    q.distinct = distinct
    q.collection = getCollectionName(collectionName)
    return q.Result(result)
}

func (q *Query) Result(result interface{}) error {

    d := q.dao

    if d.errMode == error_mode_sharing && d.LastError != nil {
        return d.LastError
    }

    if d.db == nil {
        d.Connect()
        defer d.Close()
    }

    collectionName := q.collection
    if collectionName == "" {
        collectionName = getCollectionName(result)
    }

    c := d.db.C(collectionName)

    var selector bson.M
    var mgoQuery *mgo.Query
    var err error
    var log bytes.Buffer
    log.WriteString(fmt.Sprintf("[query]collection=%s", collectionName))

    if q.dao.isID(q.queries...) {
        if IsSlice(result) {
            panic("result argument can't be a slice address")
        }
        log.WriteString(fmt.Sprintf(",id=%s", q.queries[0]))
        mgoQuery = c.FindId(q.queries[0])
    } else {
        selector = d.getM(q.queries...)
        log.WriteString(fmt.Sprintf(",selector=%s", selector))
        mgoQuery = c.Find(selector)
    }

    if q.sorts != nil && len(q.sorts) > 0 {
        mgoQuery = mgoQuery.Sort(q.sorts...)
        log.WriteString(fmt.Sprintf(",sort=%s", q.sorts))
    }

    if q.limit > 0 {
        mgoQuery = mgoQuery.Limit(q.limit)
        log.WriteString(fmt.Sprintf(",limit=%d", q.limit))
    }

    if q.page != nil {
        log.WriteString(fmt.Sprintf(",page=%s", q.page))
        var err error
        q.page.Total, err = mgoQuery.Count()
        if err != nil {
            d.LastError = err
            return err
        }
        if q.page.Count == 0 {
            q.page.Count = PAGE_RECORD_COUNT
        }
        q.page.Next = q.page.Total
        mgoQuery = mgoQuery.Skip(q.page.Cursor).Limit(q.page.Count)
    }

    if q.columns != nil && len(q.columns) > 0 {
        _columns := bson.M{}
        for _, column := range q.columns {
            _columns[column] = 1
        }
        mgoQuery.Select(_columns)
        log.WriteString(fmt.Sprintf(",columns=%s", _columns))
    }

    if q.exColumns != nil && len(q.exColumns) > 0 {
        _nColumns := bson.M{}
        for _, c := range q.exColumns {
            _nColumns[c] = 0
        }
        mgoQuery.Select(_nColumns)
        log.WriteString(fmt.Sprintf(",excludeColumns=%s", _nColumns))
    }

    if q.distinct != "" {
        log.WriteString(fmt.Sprintf(",distinct=%s", q.distinct))
        err = mgoQuery.Distinct(q.distinct, result)
    } else if IsSlice(result) {
        err = mgoQuery.All(result)
        if q.page != nil {
            len := GetValueLen(result)
            if len > 0 {
                q.page.Next = q.page.Cursor + len
            }
        }
    } else {
        err = mgoQuery.One(result)
        if q.page != nil {
            len := 1
            if err == mgo.ErrNotFound {
                len = 0
            }
            q.page.Next = q.page.Cursor + len
        }
    }

    if q.ignoreNFE && err == mgo.ErrNotFound {
        err = nil
    }

    log4go.Trace(log.String())

    if q.page != nil {
        b, _ := json.Marshal(q.page)
        log4go.Trace("page: %s", string(b))
    }
    b, _ := json.Marshal(result)
    log4go.Trace("result: %s", string(b))

    if err != nil {
        d.LastError = err
    }

    return err
}

func (q *Query) Limit(n int) *Query {
    q.limit = n
    return q
}

func (q *Query) Count(collectionName interface{}) (int, error) {

    d := q.dao

    if d.errMode == error_mode_sharing && d.LastError != nil {
        return -1, d.LastError
    }

    if d.db == nil {
        d.Connect()
        defer d.Close()
    }

    log4go.Debug("count")

    c := d.db.C(getCollectionName(collectionName))
    var mgoQuery *mgo.Query
    if d.isID(q.queries...) {
        mgoQuery = c.FindId(q.queries[0])
    } else {
        selector := d.getM(q.queries...)
        log4go.Debug("queries: %s", selector)
        mgoQuery = c.Find(selector)
    }
    n, err := mgoQuery.Count()
    log4go.Debug(n)
    if err != nil {
        d.LastError = err
    }
    return n, err
}

func (q *Query) First(result interface{}) error {
    log4go.Debug("first")
    q.page = &Page{Count: 1}
    return q.Result(result)
}

func (q *Query) Last(result interface{}) error {

    d := q.dao

    if d.errMode == error_mode_sharing && d.LastError != nil {
        return d.LastError
    }

    if d.db == nil {
        d.Connect()
        defer d.Close()
    }

    count, err := q.Count(result)
    if err != nil {
        d.LastError = err
        return err
    }
    if count == 0 {
        count = 1
    }
    log4go.Debug("last")
    q.page = &Page{Cursor: count - 1, Count: 1}
    return q.Result(result)
}

func (q *Query) Exist(collectionName interface{}) (bool, error) {
    n, err := q.Count(collectionName)
    return n > 0, err
}
