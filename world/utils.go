package world

import (
	"log"
	"math"
)

func split(in byte) (byte, byte) {
	return in >> 4, ((in << 4) >> 4)
}

func unsplit(top byte, bot byte) byte {
	return top*16 + bot
}

func toHalf(inlow byte, inhigh byte) (outtop byte, outbot byte) {
	inlowtop, inlowbot := split(inlow)
	inhightop, inhighbot := split(inhigh)

	outtop = unsplit(inhightop, inlowtop)
	outbot = unsplit(inhighbot, inlowbot)
	return
}

func toDouble(intop byte, inbot byte) (outlow byte, outhigh byte) {
	intoptop, intopbot := split(intop)
	inbottop, inbotbot := split(inbot)

	outlow = unsplit(intopbot, inbotbot)
	outhigh = unsplit(intoptop, inbottop)
	return
}

func Half(arrin []byte, top bool) (arrout []byte) {
	lenin := len(arrin)

	if math.Mod(float64(lenin), 2) != 0 {
		log.Panicf("lenin %d not even", lenin)
	}

	arrout = make([]byte, lenin/2)

	for i := range arrout {
		outtop, outbot := toHalf(arrin[i/2], arrin[i/2+1])
		if top {
			arrout[i] = outtop
		} else {
			arrout[i] = outbot
		}
	}
	return
}

func Halve(arrin []byte) (arrtop []byte, arrbot []byte) {
	lenin := len(arrin)

	if math.Mod(float64(lenin), 2) != 0 {
		log.Panicf("lenin %d not even", lenin)
	}

	arrtop = make([]byte, lenin/2)
	arrbot = make([]byte, lenin/2)

	for i := range arrtop {
		arrtop[i], arrbot[i] = toHalf(arrin[i/2], arrin[i/2+1])
	}
	return
}

func Double(top []byte, bot []byte) (full []byte) {
	lentop := len(top)
	lenbot := len(bot)

	if lentop != lenbot {
		log.Panicf("lentop %d must equal lenbot %d")
	}

	full = make([]byte, lentop+lenbot)

	for i := range top {
		di := i * 2
		full[di], full[di+1] = toDouble(top[i], bot[i])
	}
	return
}
