package world

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

func Half(arrin FullByte, top bool) (arrout HalfByte) {
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

func Halve(arrin FullByte) (arrtop HalfByte, arrbot HalfByte) {
	for i := range arrtop {
		arrtop[i], arrbot[i] = toHalf(arrin[i/2], arrin[i/2+1])
	}
	return
}

func Double(top HalfByte, bot HalfByte) (full FullByte) {
	for i := range top {
		di := i * 2
		full[di], full[di+1] = toDouble(top[i], bot[i])
	}
	return
}
