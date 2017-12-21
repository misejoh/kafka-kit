package kafkazk

import (
	"testing"
)

func TestBrokerMapFromTopicMap(t *testing.T) {
	zk := &zkmock{}
	bm, _ := zk.GetAllBrokerMeta()
	pm, _ := PartitionMapFromString(testGetMapString("test_topic"))
	forceRebuild := false

	brokers := BrokerMapFromTopicMap(pm, bm, forceRebuild)

	expected := brokerMap{
		0:    &broker{id: 0, replace: true},
		1001: &broker{id: 1001, locality: "a", used: 3, replace: false},
		1002: &broker{id: 1002, locality: "b", used: 3, replace: false},
		1003: &broker{id: 1003, locality: "c", used: 2, replace: false},
		1004: &broker{id: 1004, locality: "a", used: 2, replace: false},
	}

	for id, b := range brokers {
		switch {
		case b.id != expected[id].id:
			t.Errorf("Expected id %d, got %d for broker %d",
				expected[id].id, b.id, id)
		case b.locality != expected[id].locality:
			t.Errorf("Expected locality %s, got %s for broker %d",
				expected[id].locality, b.locality, id)
		case b.used != expected[id].used:
			t.Errorf("Expected used %d, got %d for broker %d",
				expected[id].used, b.used, id)
		case b.replace != expected[id].replace:
			t.Errorf("Expected replace %b, got %b for broker %d",
				expected[id].replace, b.replace, id)
		}
	}
}

// func TestBestCandidate(t *testing.T) {}
// func TestConstraintsAdd(t *testing.T) {}
// func TestConstraintsPasses(t *testing.T) {}
// func TestMergeConstraints(t *testing.T) {}
// func TestUpdate(t *testing.T) {}
// func TestFilteredList(t *testing.T) {}
// func TestBrokerStringToSlice(t *testing.T) {}