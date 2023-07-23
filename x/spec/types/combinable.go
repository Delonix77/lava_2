package types

import fmt "fmt"

type CurrentContainer struct {
	index             int
	currentCombinable Combinable
}

type Combinable interface {
	Differeniator() string
	GetEnabled() bool
	Equal(interface{}) bool
	Overwrite(other Combinable) (overwritten Combinable, isOverwritten bool) // combinable.Overwrite(other) is called when combinable is overwriting other
}

func (api *Api) Differeniator() string {
	return api.Name
}

func (h *Header) Differeniator() string {
	return h.Name
}

func (parsing *ParseDirective) Differeniator() string {
	return parsing.FunctionTag.String()
}

func (e *Extension) Differeniator() string {
	return e.Name
}

func (v *Verification) Differeniator() string {
	return v.Name
}

func (api *Api) Overwrite(other Combinable) (Combinable, bool) {
	return api, false
}

func (h *Header) Overwrite(Combinable) (Combinable, bool) {
	return h, false
}

func (parsing *ParseDirective) Overwrite(Combinable) (Combinable, bool) {
	return parsing, false
}

func (e *Extension) Overwrite(Combinable) (Combinable, bool) {
	return e, false
}

func (v *Verification) Overwrite(other Combinable) (Combinable, bool) {
	if v.ParseDirective == nil {
		if otherVerification, ok := other.(*Verification); ok && otherVerification.ParseDirective != nil {
			v.ParseDirective = otherVerification.ParseDirective
			return v, true
		}
	}
	return v, false
}

// Api has GetEnabled

func (h *Header) GetEnabled() bool {
	return true
}

func (h *ParseDirective) GetEnabled() bool {
	return true
}

func (h *Extension) GetEnabled() bool {
	return true
}

func (h *Verification) GetEnabled() bool {
	return true
}

func CombineFields[T Combinable](currentMap map[string]CurrentContainer, mergedMap map[string]interface{}, combineFrom []T, toCombine []T, allowOverwrite bool) ([]T, map[string]interface{}, error) {
	for _, combinable := range combineFrom {
		if !combinable.GetEnabled() {
			continue
		}
		if existing, ok := mergedMap[combinable.Differeniator()]; ok {
			// only matters if they are not the same
			if !combinable.Equal(existing) {
				// was already existing in mergedApis
				if !allowOverwrite {
					// if we don't allow overwriting by the calling collection, then a duplication is an error, but only if it's not equal
					return nil, nil, fmt.Errorf("try overwrite existing collection combination %s", combinable.Differeniator())
				}
				if _, found := currentMap[combinable.Differeniator()]; !found {
					return nil, nil, fmt.Errorf("duplicate imported combinable: %s", combinable.Differeniator())
				}
			}
		}
		mergedMap[combinable.Differeniator()] = combinable
		toCombine = append(toCombine, combinable)
	}
	return toCombine, mergedMap, nil
}

func GetCurrentFromCombinable[T Combinable](current []T) (currentMap map[string]CurrentContainer) {
	currentMap = map[string]CurrentContainer{}
	for idx, combinable := range current {
		currentMap[combinable.Differeniator()] = CurrentContainer{
			index:             idx,
			currentCombinable: combinable,
		}
	}
	return currentMap
}

func CombineUnique[T Combinable](appendFrom []T, appendTo []T, currentMap map[string]CurrentContainer, allowOverwrite bool) ([]T, error) {
	for _, combinable := range appendFrom {
		if current, found := currentMap[combinable.Differeniator()]; !found {
			appendTo = append(appendTo, combinable)
		} else if !allowOverwrite {
			// if they are the same it's allowed
			if !combinable.Equal(current.currentCombinable) {
				return nil, fmt.Errorf("existing combinable in collection combination %s ", combinable.Differeniator())
			}
		} else {
			// overwriting the inherited field might need Overwrite actions
			if overwritten, isOverwritten := current.currentCombinable.Overwrite(combinable); isOverwritten {
				if appendTo[current.index].Differeniator() != combinable.Differeniator() {
					return nil, fmt.Errorf("differentiator mismatch in overwrite %s vs %s", combinable.Differeniator(), appendTo[current.index].Differeniator())
				}
				appendTo[current.index] = overwritten.(T)
			}
		}
	}
	return appendTo, nil
}
