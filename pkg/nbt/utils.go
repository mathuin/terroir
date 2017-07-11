package nbt

type CompoundElem struct {
	Key   string
	Tag   byte
	Value interface{}
}

func MakeCompoundPayload(elems []CompoundElem) []Tag {

	compoundPayload := make([]Tag, 0)

	for _, elem := range elems {
		newTag := MakeTag(elem.Tag, elem.Key)
		newTag.SetPayload(elem.Value)
		compoundPayload = append(compoundPayload, newTag)
	}

	return compoundPayload
}

func MakeCompound(name string, elems []CompoundElem) Tag {

	compoundPayload := MakeCompoundPayload(elems)

	compoundTag := MakeTag(TAG_Compound, name)
	compoundTag.SetPayload(compoundPayload)

	return compoundTag
}
