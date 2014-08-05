package world

import "testing"

func Test_sectionWrite(t *testing.T) {
	s := MakeSection()

	// populate it
	for i := range s.Blocks {
		s.Blocks[i] = 0x31
	}
	for i := range s.Add {
		s.Add[i] = 0x12
		s.Data[i] = 0x34
		s.BlockLight[i] = 0x56
		s.SkyLight[i] = 0x78
	}

	sTagPayload := s.write(1)

	requiredTags := map[string]bool{
		"Y":          false,
		"Blocks":     false,
		"Add":        false,
		"Data":       false,
		"BlockLight": false,
		"SkyLight":   false,
	}

	for _, tag := range sTagPayload {
		if _, ok := requiredTags[tag.Name]; ok {
			requiredTags[tag.Name] = true
			switch tag.Name {
			case "Y":
				if tag.Payload.(byte) != 0x01 {
					t.Errorf("Y value does not match")
				}
			case "Blocks":
				if tag.Payload.([]byte)[0] != 0x31 {
					t.Errorf("Block value does not match")
				}
			case "Add":
				if tag.Payload.([]byte)[0] != 0x12 {
					t.Errorf("Add value does not match")
				}
			case "Data":
				if tag.Payload.([]byte)[0] != 0x34 {
					t.Errorf("Data value does not match")
				}
			case "BlockLight":
				if tag.Payload.([]byte)[0] != 0x56 {
					t.Errorf("BlockLight value does not match")
				}
			case "SkyLight":
				if tag.Payload.([]byte)[0] != 0x78 {
					t.Errorf("SkyLight value does not match")
				}
			}
		} else {
			t.Errorf("tag name %s not required for section", tag.Name)
		}
	}

	for rtkey, rtval := range requiredTags {
		if rtval == false {
			t.Errorf("tag name %s required for section but not found", rtkey)
		}
	}
}
