package sequencediagram

import "testing"

func TestCalcOffsets(t *testing.T) {
	sd, _ := ParseFromText("ayy->lmao:fam\nlmao->lmao:message\nayyy->ayy:msg")
	t.Errorf("\n%v", calcOffsets(sd))
	t.Errorf("\n%v", sd.AsTextDiagram())
}
