package input

import (
	"fmt"
	"regexp"

	"github.com/lamasutra/bg-music/wt-client/model"
	"github.com/lamasutra/bg-music/wt-client/ui"
)

const anyKillPattern = `shot\s+down|destroyed`
const playerMadeKillPatternTemplate = `%s.+(shot\s+down|destroyed)`
const playerMadeSeverDamagPatternTemplate = `%s.+severely damaged`
const playerIsShotDownPatternTemplate = `shot down %s`
const playerHasCrashedPatternTemplate = `%s.+has crashed`
const playerIsBurningTemplate = `set afire %s`
const playerIsCritDamagedTemplate = `critically damaged %s`
const playerIsSeverlyDamagedTemplate = `severely damaged %s`

var anyKillRegExp regexp.Regexp
var playerMadeKilledRegExp regexp.Regexp
var playerMadeSeverDamagedRegExp regexp.Regexp
var playerIsShotDownRegExp regexp.Regexp
var playerHasCrashedRegExp regexp.Regexp
var playerIsBurningRegExp regexp.Regexp
var playerIsCritDamagedRegExp regexp.Regexp
var playerIsSeverlyDamagedRegExp regexp.Regexp

// const playerShotDownRegExpTemplate = `%s\s+\([^\)]+\)\s+(shot\s+down|destroyed)\s+` // [^\(]+\([^\)]+\)

// const playerShootDownEnemyRegexp = `[^\(]+\s+\([^\)]+\)\s+shot\s+down\s+[^\(]+\([^\)]+\)`

func initPatterns(conf *model.Config) {
	playerMadeKilledRegExp = *regexp.MustCompile(fmt.Sprintf(playerMadeKillPatternTemplate, conf.Nickname))
	playerMadeSeverDamagedRegExp = *regexp.MustCompile(fmt.Sprintf(playerMadeSeverDamagPatternTemplate, conf.Nickname))
	anyKillRegExp = *regexp.MustCompile(anyKillPattern)
	playerHasCrashedRegExp = *regexp.MustCompile(fmt.Sprintf(playerHasCrashedPatternTemplate, conf.Nickname))
	playerIsShotDownRegExp = *regexp.MustCompile(fmt.Sprintf(playerIsShotDownPatternTemplate, conf.Nickname))
	playerIsBurningRegExp = *regexp.MustCompile(fmt.Sprintf(playerIsBurningTemplate, conf.Nickname))
	playerIsCritDamagedRegExp = *regexp.MustCompile(fmt.Sprintf(playerIsCritDamagedTemplate, conf.Nickname))
	playerIsSeverlyDamagedRegExp = *regexp.MustCompile(fmt.Sprintf(playerIsSeverlyDamagedTemplate, conf.Nickname))

	ui.Debug("pattern initialized\n",
		"made kill", playerMadeKilledRegExp, "\n",
		"made sever damage", playerMadeSeverDamagedRegExp, "\n",
		"any kill", anyKillRegExp, "\n",
		"is shot down", playerIsShotDownRegExp, "\n",
		"is burning", playerIsBurningRegExp, "\n",
		"is crit damaged", playerIsCritDamagedRegExp, "\n",
		"is severely damaged", playerIsSeverlyDamagedRegExp, "\n")
}
