package database

import (
	"errors"
	"fmt"
	"github.com/statping-ng/statping-ng/types"
	"github.com/statping-ng/statping-ng/utils"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GroupBy struct {
	db    Database
	query *GroupQuery
}

type GroupByer interface {
	ToTimeValue() (*TimeVar, error)
}

type By string

func (b By) String() string {
	return string(b)
}

type GroupQuery struct {
	Start     time.Time
	End       time.Time
	Group     time.Duration
	Order     string
	Limit     int
	Offset    int
	FillEmpty bool

	db Database
}

func (b GroupQuery) Find(data interface{}) error {
	return b.db.Order("id DESC").Find(data).Error()
}

func (b GroupQuery) Database() Database {
	return b.db
}

var (
	ByCount   = By("COUNT(id) as amount")
	ByAverage = func(column string, multiplier int) By {
		switch database.DbType() {
		case "mysql":
			return By(fmt.Sprintf("CAST(AVG(%s) as UNSIGNED INT) as amount", column))
		case "postgres":
			return By(fmt.Sprintf("cast(AVG(%s) as int) as amount", column))
		default:
			return By(fmt.Sprintf("cast(AVG(%s) as int) as amount", column))
		}
	}
)

type TimeVar struct {
	g    *GroupQuery
	data []*TimeValue
}

func (t *TimeVar) ToValues() ([]*TimeValue, error) {
	return t.data, nil
}

// GraphData will return all hits or failures, without outage_type selection
func (b *GroupQuery) GraphData(by By) ([]*TimeValue, error) {
	b.db = b.db.MultipleSelects(
		b.db.SelectByTime(b.Group),
		by.String(),
	).Group("timeframe").Order("timeframe", true)

	caller, err := b.ToTimeValue()
	if err != nil {
		return nil, err
	}

	if b.FillEmpty {
		return caller.FillMissing(b.Start, b.End)
	}
	return caller.ToValues()
}

// GraphDataForFailures will return failures data with outage_type selection
func (b *GroupQuery) GraphDataForFailures(by By) ([]*TimeValue, error) {
    selectExpr := fmt.Sprintf(`%s, 
        CASE 
            WHEN SUM(CASE WHEN outage_type = 'critical' THEN 1 ELSE 0 END) > 0 THEN 'critical'
            WHEN SUM(CASE WHEN outage_type = 'major' THEN 1 ELSE 0 END) > 0 THEN 'major'
            WHEN SUM(CASE WHEN outage_type = 'minor' THEN 1 ELSE 0 END) > 0 THEN 'minor'
            ELSE ''
        END as outage_type`, by.String())

    b.db = b.db.MultipleSelects(b.db.SelectByTime(b.Group), selectExpr).Group("timeframe").Order("timeframe", true)

    caller, err := b.ToTimeValueForFailures()

    if err != nil {
        log.Errorf("GraphDataForFailures: Error in query execution: %v", err)
        return nil, err
    }

    if b.FillEmpty {
        filled, err := caller.FillMissing(b.Start, b.End)
        if err != nil {
            log.Errorf("GraphDataForFailures: Error in FillMissing: %v", err)
            return nil, err
        }
        return filled, nil
    }
    return caller.ToValues()
}

// ToTimeValue will format the SQL rows into a JSON format for the API.
// [{"timestamp": "2006-01-02T15:04:05Z", "amount": 468293}]
// TODO redo this entire function, use better SQL query to group by time
func (b *GroupQuery) ToTimeValue() (*TimeVar, error) {
	rows, err := b.db.Rows()
	if err != nil {
		return nil, err
	}
	var data []*TimeValue
	for rows.Next() {
		var timeframe string
		var amount int64
		if err := rows.Scan(&timeframe, &amount); err != nil {
			log.Error(err, timeframe)
		}
		trueTime, _ := b.db.ParseTime(timeframe)
		newTs := types.FixedTime(trueTime, b.Group)

		tv := &TimeValue{
			Timeframe: newTs,
			Amount:    amount,
		}
		data = append(data, tv)
	}
	return &TimeVar{b, data}, nil
}

func (b *GroupQuery) ToTimeValueForFailures() (*TimeVar, error) {
    rows, err := b.db.Rows()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var data []*TimeValue
    for rows.Next() {
        var timeframe string
        var amount int64
        var outageType string

        if err := rows.Scan(&timeframe, &amount, &outageType); err != nil {
            log.Errorln(err, timeframe)
            continue // ou return nil, err selon le comportement souhaité
        }
        trueTime, _ := b.db.ParseTime(timeframe)
        newTs := types.FixedTime(trueTime, b.Group)

        tv := &TimeValue{
            Timeframe:  newTs,
            Amount:     amount,
            OutageType: outageType,
        }
        data = append(data, tv)
    }
    return &TimeVar{b, data}, nil
}


// FillMissing fills in missing time slots between the provided current and end times,
// using aggregated data from t.data. For each time slot, the function aggregates any
// available entries by summing their amounts and by selecting the outage type with the
// highest severity. Severity is ranked as follows:
//   "critical" > "major" > "minor" > "" (no outage)
// 
// The display color (for UI purposes) is derived from the aggregated result:
//   - If no data is available for a time slot, the color is grey.
//   - If data is available but the total failure amount is 0, the color is green.
//   - If the highest severity is "minor" or "major", the color is orange.
//   - If the highest severity is "critical", the color is red.
//
// The function returns a slice of pointers to TimeValue structs representing each
// time slot from current until end. If no record exists for a time slot, the slot is
// filled with an amount of 0 and an empty outage type.
//
// Examples (where the three columns represent, respectively, the color of individual
// entries being aggregated, the color of the newly added entry, and the final aggregated
// result color):
//
//   Input(s)           -> Aggregated Result (Display Color)
//   -----------------------------------------------------------
//   green - green      -> green   // (0 failure, and additional 0 → green)
//   green - orange     -> orange  // (0 failure combined with a minor/major outage → orange)
//   green - red        -> red     // (0 failure combined with a critical outage → red)
//   [no data]          -> gray    // (no record for this time slot → gray)
//
// More generally:
//
//   Condition                          			| Outage Type       	| Display Color
//   -----------------------------------------------|-----------------------|--------------
//   No data                            			| N/A               	| Gray
//   Data exists, total failures = 0    			| "" (empty)        	| Green
//   Data exists, highest severity = minor/major 	| "minor" or "major" 	| Orange
//   Data exists, highest severity = critical      	| "critical"       		| Red
//
// Args:
//   current: The starting time of the interval to fill.
//   end: The ending time of the interval to fill.
//
// Returns:
//   A slice of pointers to TimeValue, one for each time slot, containing the
//   aggregated amount and outage type (from which a display color can be derived),
//   and an error if any occurred.
func (t *TimeVar) FillMissing(current, end time.Time) ([]*TimeValue, error) {
	// aggregated holds the sum and the highest severity outage type for a given timeframe.
	type aggregated struct {
		amount     int64
		outageType string
	}

	// severityRank defines the relative severity of outage types.
	// "critical" is the highest, followed by "major" and "minor".
	severityRank := map[string]int{
		"critical": 3,
		"major":    2,
		"minor":    1,
		"":         0,
	}

	// Build a map of timeframe → aggregated data.
	aggMap := make(map[string]aggregated)
	for _, v := range t.data {
		key := v.Timeframe
		if existing, ok := aggMap[key]; ok {
			// Sum the amounts.
			newAmount := existing.amount + v.Amount
			// Choose the outage type with the higher severity.
			if severityRank[v.OutageType] > severityRank[existing.outageType] {
				aggMap[key] = aggregated{amount: newAmount, outageType: v.OutageType}
			} else {
				aggMap[key] = aggregated{amount: newAmount, outageType: existing.outageType}
			}
		} else {
			aggMap[key] = aggregated{amount: v.Amount, outageType: v.OutageType}
		}
	}

	var validSet []*TimeValue
	// Iterate over each time slot from current to end.
	for {
		currentStr := types.FixedTime(current, t.g.Group)
		if agg, ok := aggMap[currentStr]; ok {
			validSet = append(validSet, &TimeValue{
				Timeframe:  currentStr,
				Amount:     agg.amount,
				OutageType: agg.outageType,
			})
		} else {
			// No record for this time slot: fill with zero amount and no outage type.
			validSet = append(validSet, &TimeValue{
				Timeframe:  currentStr,
				Amount:     0,
				OutageType: "",
			})
		}

		current = current.Add(t.g.Group)
		if current.After(end) {
			break
		}
	}

	return validSet, nil
}

type isObject interface {
	Db() Database
}

func ParseRequest(r *http.Request) (*GroupQuery, error) {
	fields := parseGet(r)
	grouping := fields.Get("group")
	startField := utils.ToInt(fields.Get("start"))
	endField := utils.ToInt(fields.Get("end"))
	limit := utils.ToInt(fields.Get("limit"))
	offset := utils.ToInt(fields.Get("offset"))
	fill, _ := strconv.ParseBool(fields.Get("fill"))
	orderBy := fields.Get("order")
	if limit == 0 {
		limit = 10000
	}

	if grouping == "" {
		grouping = "1h"
	}
	groupDur, err := time.ParseDuration(grouping)
	if err != nil {
		log.Errorln(err)
		groupDur = 1 * time.Hour
	}

	query := &GroupQuery{
		Start:     time.Unix(startField, 0).UTC(),
		End:       time.Unix(endField, 0).UTC(),
		Group:     groupDur,
		Order:     orderBy,
		Limit:     int(limit),
		Offset:    int(offset),
		FillEmpty: fill,
	}

	if query.Start.After(query.End) {
		return nil, errors.New("start time is after ending time")
	}

	return query, nil
}

func ParseQueries(r *http.Request, o isObject) (*GroupQuery, error) {
	return ParseQueriesForTable(r, o, "")
}

func ParseQueriesForTable(r *http.Request, o isObject, whereTable string) (*GroupQuery, error) {
	fields := parseGet(r)
	grouping := fields.Get("group")
	startField := utils.ToInt(fields.Get("start"))
	endField := utils.ToInt(fields.Get("end"))
	limit := utils.ToInt(fields.Get("limit"))
	offset := utils.ToInt(fields.Get("offset"))
	fill, _ := strconv.ParseBool(fields.Get("fill"))
	orderBy := fields.Get("order")
	if limit == 0 {
		limit = 10000
	}

	q := o.Db()

	if grouping == "" {
		grouping = "1h"
	}
	groupDur, err := time.ParseDuration(grouping)
	if err != nil {
		log.Errorln(err)
		groupDur = 1 * time.Hour
	}
	if endField == 0 {
		endField = utils.Now().Unix()
	}

	query := &GroupQuery{
		Start:     time.Unix(startField, 0).UTC(),
		End:       time.Unix(endField, 0).UTC(),
		Group:     groupDur,
		Order:     orderBy,
		Limit:     int(limit),
		Offset:    int(offset),
		FillEmpty: fill,
		db:        q,
	}

	if query.Start.After(query.End) {
		return nil, errors.New("start time is after ending time")
	}

	if startField == 0 {
		query.Start = utils.Now().Add(-7 * types.Day)
	}
	if endField == 0 {
		query.End = utils.Now()
	}
	if query.Limit != 0 {
		q = q.Limit(query.Limit)
	}
	if query.Offset > 0 {
		q = q.Offset(query.Offset)
	}

	if len(whereTable) > 0 {
		whereTable = fmt.Sprintf("%s.", whereTable)
	}

	q = q.Where(fmt.Sprintf("%screated_at BETWEEN ? AND ?", whereTable), q.FormatTime(query.Start), q.FormatTime(query.End))

	if query.Order != "" {
		q = q.Order(query.Order)
	}
	query.db = q

	return query, nil
}

func parseForm(r *http.Request) url.Values {
	r.ParseForm()
	return r.PostForm
}

func parseGet(r *http.Request) url.Values {
	r.ParseForm()
	return r.Form
}
