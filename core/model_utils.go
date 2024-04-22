package elemental

import (
	"elemental/utils"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m Model[T]) addToFilters(key string, value any) Model[T] {
	stage := "$match"
	foundMatchStage := false
	m.pipeline = lo.Map(m.pipeline, func(stg bson.D, _ int) bson.D {
		filters := e_utils.Cast[primitive.M](e_utils.CastBSON[bson.M](stg)[stage])
		if filters != nil {
			foundMatchStage = true
			if (m.orConditionActive) {
				if filters["$or"] == nil {
					filters["$or"] = []primitive.M{
						{m.whereField: primitive.M{key: value}},
					}
				} else {
					filters["$or"] = append(filters["$or"].([]primitive.M), primitive.M{m.whereField: primitive.M{key: value}})
				}
				for k, v := range filters {
					if k != "$or" {
						filters["$or"] = append(filters["$or"].([]primitive.M), primitive.M{k: v})
						delete(filters, k)
					}
				}	
				m.orConditionActive = false
			} else {
				filterExistsWithinAndOperator := false
				if filters["$and"] != nil {
					for _, filter := range filters["$and"].([]primitive.M) {
						if filter[m.whereField] != nil {
							filterExistsWithinAndOperator = true
							filters["$and"] = append(filters["$and"].([]primitive.M), primitive.M{m.whereField: primitive.M{key: value}})
						}
					}
				}
				if !filterExistsWithinAndOperator {
					if filters[m.whereField] == nil {
						filters[m.whereField] = primitive.M{key: value}
					} else {
						filters["$and"] = []primitive.M{
							{m.whereField: filters[m.whereField]},
							{m.whereField: primitive.M{key: value}},
						}
						delete(filters, m.whereField)
					}
				}
			}
			return bson.D{{Key: stage, Value: filters}}
		}
		return stg
	})
	if !foundMatchStage {
		m.pipeline = append(m.pipeline, bson.D{{Key: stage, Value: primitive.M{m.whereField: primitive.M{key: value}}}})
		return m
	}
	return m
}

func (m Model[T]) addToPipeline(stage, key string, value any) Model[T] {
	foundStage := false
	m.pipeline = lo.Map(m.pipeline, func(stg bson.D, _ int) bson.D {
		stageObject := e_utils.Cast[bson.D](e_utils.CastBSON[bson.D](stg).Map()[stage])
		if stageObject != nil {
			foundStage = true
			if stageObject.Map()[key] == nil {
				stageObject = append(stageObject, bson.E{Key: key, Value: value})
			}
			return bson.D{{Key: stage, Value: stageObject}}
		}
		return stg
	})
	if !foundStage {
		m.pipeline = append(m.pipeline, bson.D{{Key: stage, Value: bson.D{{Key: key, Value: value}}}})
		return m
	}
	return m
}

func (m Model[T]) checkConditionsAndPanic(results []T) {
	if m.failWith != nil && len(results) == 0 {
		panic(*m.failWith)
	}
}
