package debug

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	NSFile = "/Users/thebevrishot/Workspaces/chain/spam-ns.json"
)

var nsmap = map[string]string{}
var inited = false

func finit() {
	if inited {
		return
	}

	reversedNsMap := map[string]string{}

	_, err := os.Stat(NSFile)
	if os.IsNotExist(err) {
		return
	}

	f, err := os.Open(NSFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.NewDecoder(f).Decode(&reversedNsMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range reversedNsMap {
		nsmap[v] = k
	}

	inited = true
}

func lookup(cd []byte) string {
	h := hex.EncodeToString(cd)
	if v, ok := nsmap[h]; ok {
		return v
	}

	return h
}

type Event struct {
	Table     string `json:"table"`
	Name      string `json:"name"`
	RequestID int64  `json:"request_id"`
	CreatedAt int64  `json:"created_at"`
	GasUsed   int64  `json:"gas_used"`
	Time      int64  `json:"time"`
	Memo      string `json:"memo,omitempty"`
}

type RequestStat struct {
	// Type                 string        `json:""`
	RequestID            int64 `json:"request_id"`
	CreatedAt            int64 `json:"created_at"`
	PrepareGasUsage      int64 `json:"prepare_gas_usage"`
	PrepareTime          int64 `json:"prepare_time"`
	PrepareOwasmGasUsage int64 `json:"prepare_owasm_gas_usage"`
	PrepareOwasmTime     int64 `json:"prepare_owasm_time"`
	ResolveAtBlock       int64 `json:"resolve_at_block"`
	ResolveGasUsage      int64 `json:"resolve_gas_usage"`
	ResolveTime          int64 `json:"resolve_time"`
	ExecuteGasUsage      int64 `json:"execute_gas_usage"`
	ExecuteTime          int64 `json:"execute_time"`
}

func NewRequest(
	requestID int64,
	createdAt int64,
	usedGas int64,
	usedTime int64,
	owasmGas int64,
	owasmTime int64,
) RequestStat {
	return RequestStat{
		RequestID:            requestID,
		CreatedAt:            createdAt,
		PrepareGasUsage:      usedGas,
		PrepareTime:          usedTime,
		PrepareOwasmGasUsage: owasmGas,
		PrepareOwasmTime:     owasmTime,
	}
}

func (r *RequestStat) Resolve(
	resolveAt int64,
	resolveGasUsage int64,
	resolveTime int64,
	executeGasUsage int64,
	executeTime int64,
) {
	r.ResolveAtBlock = resolveAt
	r.ResolveGasUsage = resolveGasUsage
	r.ResolveTime = resolveTime
	r.ExecuteGasUsage = executeGasUsage
	r.ExecuteTime = executeTime
}

type Stat struct {
	name     string
	count    int64
	time     int64
	countAcc int64
	timeAcc  int64
	gas      int64
	gasAcc   int64
}

func NewStat(name string) Stat {
	s := Stat{name: name}
	s.ResetHard()
	return s
}

func (s *Stat) Add(start time.Time, gas int64) {
	s.count++
	s.countAcc++
	s.gas += gas

	dur := time.Since(start)
	s.time += dur.Microseconds()
	s.timeAcc += dur.Microseconds()
	s.gasAcc += gas
}

func (s *Stat) AddStat(m map[string]interface{}) {
	m[s.name+".count"] = s.count
	m[s.name+".time"] = s.time
	m[s.name+".countAcc"] = s.countAcc
	m[s.name+".timeAcc"] = s.timeAcc
	m[s.name+".gas"] = s.gas
	m[s.name+".gasAcc"] = s.gasAcc
}

func (s *Stat) Reset() {
	s.count = 0
	s.time = 0
	s.gas = 0
}

func (s *Stat) ResetHard() {
	s.Reset()
	s.countAcc = 0
	s.timeAcc = 0
	s.gasAcc = 0
}

type BlockState struct {
	CurrentBlock     int64
	CurrentBlockTime time.Time
	requestsStat     Stat
	resolvesStat     Stat
	preparesStat     Stat
	executeStat      Stat
	reportStat       Stat

	requests map[int64]RequestStat
	events   []Event

	db *DB
	// TODO
	// dbClient *r.Session
	// table    string

	hook *Hook
}

func NewBlockState(dbAddress string, table string) *BlockState {
	s := &BlockState{
		requestsStat: NewStat("requests"),
		resolvesStat: NewStat("resolve"),
		preparesStat: NewStat("prepare"),
		executeStat:  NewStat("execute"),
		reportStat:   NewStat("reports"),
		requests:     map[int64]RequestStat{},
		events:       []Event{},
	}
	if dbAddress != "" {
		finit()

		session, err := r.Connect(r.ConnectOpts{
			Address: dbAddress,
		})

		if err != nil {
			panic("init db fail : " + err.Error())
		}

		if err := r.DBCreate(table).Exec(session); err != nil {
			fmt.Println("warning", "init table fail : ", err.Error())
		}

		if err := r.DB(table).TableCreate("blocks").Exec(session); err != nil {
			fmt.Println("warning", "init table fail : ", err.Error())
		}

		if err := r.DB(table).TableCreate("requests").Exec(session); err != nil {
			fmt.Println("warning", "init table fail : ", err.Error())
		}

		if err := r.DB(table).TableCreate("events").Exec(session); err != nil {
			fmt.Println("warning", "init table fail : ", err.Error())
		}

		if err := r.DB(table).TableCreate("block-stats").Exec(session); err != nil {
			fmt.Println("warning", "init table fail : ", err.Error())
		}

		s.db = NewDB(session, table, 100)

		s.hook = NewHook(s.db, table)
	}

	return s
}

func (s *BlockState) Hook() *Hook {
	return s.hook
}

func (s *BlockState) SetPendingAmount(x int) {
	if s.hook != nil {
		s.hook.hookBlock.Resolved = x
	}
}

func (s *BlockState) Reset() {
	s.requestsStat.Reset()
	s.resolvesStat.Reset()
	s.preparesStat.Reset()
	s.executeStat.Reset()
	s.reportStat.Reset()
}

func (s *BlockState) Record() {
	if s.db != nil {
		m := map[string]interface{}{
			"currentBlock":     s.CurrentBlock,
			"currentBlockTime": s.CurrentBlockTime,
		}
		s.requestsStat.AddStat(m)
		s.resolvesStat.AddStat(m)
		s.preparesStat.AddStat(m)
		s.executeStat.AddStat(m)
		s.reportStat.AddStat(m)

		s.db.Insert("blocks", m)

		recorded := []int64{}
		for _, req := range s.requests {
			if req.ResolveAtBlock == 0 {
				continue
			}
			s.db.Insert("requests", req)

			recorded = append(recorded, req.RequestID)
		}

		for _, id := range recorded {
			delete(s.requests, id)
		}

		for _, e := range s.events {
			e.Table = "events"
			s.db.Insert("events", e)
		}
	}
}

func (s *BlockState) Request(start time.Time, gas int64) {
	s.requestsStat.Add(start, gas)
}

func (s *BlockState) Resolve(start time.Time, gas int64) {
	s.resolvesStat.Add(start, gas)
}

func (s *BlockState) Prepare(start time.Time, gas int64) {
	s.preparesStat.Add(start, gas)
}

func (s *BlockState) Execute(start time.Time, gas int64) {
	s.executeStat.Add(start, gas)
}

func (s *BlockState) Report(start time.Time, gas int64) {
	s.reportStat.Add(start, gas)
}

func (s *BlockState) NewRequest(
	// requestID int64,
	createdAt int64,
	usedGas int64,
	used time.Duration,
	owasmGas int64,
	owasmUsed time.Duration,
	calldata []byte,
) {
	memo := lookup(calldata)
	s.events = append(s.events, Event{
		Name:      "request",
		RequestID: 0,
		CreatedAt: createdAt,
		GasUsed:   usedGas,
		Time:      used.Microseconds(),
		Memo:      memo,
	}, Event{
		Name:      "owasm-request",
		RequestID: 0,
		CreatedAt: createdAt,
		GasUsed:   owasmGas,
		Time:      owasmUsed.Microseconds(),
		Memo:      memo,
	})
	// s.requests[requestID] = NewRequest(requestID, createdAt, usedGas, used, owasmGas, owasmUsed)
}

func (s *BlockState) FinishResolve(
	requestID int64,
	resolveAt int64,
	resolveGasUsage int64,
	used time.Duration,
	executeGasUsage int64,
	executeUsed time.Duration,
	calldata []byte,
) {
	memo := lookup(calldata)
	s.events = append(s.events, Event{
		Name:      "resolve",
		RequestID: requestID,
		CreatedAt: resolveAt,
		GasUsed:   resolveGasUsage,
		Time:      used.Microseconds(),
		Memo:      memo,
	}, Event{
		Name:      "owasm-execute",
		RequestID: requestID,
		CreatedAt: resolveAt,
		GasUsed:   executeGasUsage,
		Time:      executeUsed.Microseconds(),
		Memo:      memo,
	})
	// if r, ok := s.requests[requestID]; ok {
	// 	r.Resolve(resolveAt, resolveGasUsage, used, executeGasUsage, executeUsed)
	// 	s.requests[requestID] = r
	// } else {
	// 	panic("No request ID")
	// }
}
