//go:build porttest

package legacy

func themeInvokeGo(op int, v [6]uint32) uint32 {
	switch op {
	case 0:
		return mapThemeFile(v[0], v[1])
	case 1:
		return mapThemeRead(v[0], mapThemeByte(v[1], 0))
	case 2:
		return mapThemeRaw(v[0], mapThemeByte(v[1], 0))
	case 3:
		return mapThemeSkipLine(v[0])
	case 4:
		return mapThemeControl(v[0])
	case 5:
		return mapThemeSkip(v[0], false)
	case 6:
		return mapThemeSkip(v[0], true)
	case 7:
		return mapThemeCondition(v[0], populationWord(v[1], 0))
	case 8:
		return mapThemeOperator(v[0], populationWord(v[1], 0))
	case 9:
		return mapThemeAlgorithm(v[0], v[1])
	case 10:
		return mapThemeSpells(v[0], v[1])
	case 11:
		return mapThemeEquipment(v[0], v[1], false)
	case 12:
		mapThemeEquipmentFree(v[0])
		return 0
	case 13:
		return mapThemeAttributes(v[0], v[1])
	case 14:
		return mapThemeEquipment(v[0], v[1], true)
	case 15:
		return mapThemeExit(v[0], v[1])
	case 16:
		return mapThemeDecor(v[0], v[1])
	case 17:
		return mapThemeWallFloor(v[0], v[1])
	case 18:
		return mapThemeEdging(v[0], int32(v[1]), v[2])
	case 19:
		return mapThemeDecorSet(v[0], v[1])
	case 20:
		return mapThemeChoices(v[0])
	case 21:
		return mapThemeForeach(v[0])
	case 22:
		return mapThemeDecorCopy(v[0], v[1], v[2])
	case 23, 24, 25, 26, 27, 28:
		return mapThemeProperty(v[0], v[1], op)
	case 29:
		return mapThemeValidate(v[0])
	case 30:
		return mapThemePrefabs(v[0], v[1])
	case 31:
		return mapThemeNewPrefab(v[0], v[1])
	case 32:
		return mapThemeCleanup(v[0])
	}
	panic("theme operation")
}
